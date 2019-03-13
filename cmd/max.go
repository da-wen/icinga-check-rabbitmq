// Copyright © 2019 dawen <dan.wendlandt@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var url, index, field, auth string
var ecode, period, cTimeout int
var treshold int64
var qw, vhosts []string

type ResultJSON []struct {
	Name     string `json:"name"`
	Messages int64  `json:"messages"`
	Vhost    string `json:"vhost"`
}

// maxCmd represents the max command
var maxCmd = &cobra.Command{
	Use:   "max",
	Short: "Checks if max messages treshold is reached",
	Long: `Fails with exit code 2 , if doc count for a period is higher then defined treshold
`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(maxCmd)
	// required flags
	maxCmd.Flags().StringVar(&url, "url", "", "string of url like http://localhost:15672")
	maxCmd.MarkFlagRequired("url")
	//optional
	maxCmd.Flags().Int64VarP(&treshold, "max", "m", 1000, "defines minimum amount of docs that are required when command fails")
	maxCmd.Flags().IntVarP(&ecode, "exit", "e", 2, "exit code to be used for fail")
	maxCmd.Flags().StringVarP(&auth, "auth", "a", "", "basic auth for header authenticatio. format=username:password")
	maxCmd.Flags().StringSliceVarP(&qw, "whitelist", "w", nil, "comma-separated  whitlist for queue")
	maxCmd.Flags().StringSliceVarP(&vhosts, "vhosts", "v", nil, "comma-separated  whitlist for vhosts")
}

func logExit(m string, err error) {
	if err != nil {
		fmt.Println(m, err)
	} else {
		fmt.Println(m)
	}

	os.Exit(5)
}

func run() {
	p := fmt.Sprintf("%s/api/queues", url)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: false,
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(http.MethodGet, p, nil)
	if err != nil {
		logExit("Request can not be created", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if len(auth) != 0 {
		bAuth := []byte(auth)
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(bAuth))
	}

	res, err := client.Do(req)
	if err != nil {
		logExit("Rabbitmq request error", err)
	}
	defer res.Body.Close()

	//fmt.Printf("%+v", res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logExit("Status response not readable", err)
	}
	if res.StatusCode >= 400 {
		var resErr map[string]interface{}
		json.Unmarshal(body, &resErr)
		logExit(fmt.Sprintf("Response code invalid: %s", res.Status), nil)
	}

	var result ResultJSON
	json.Unmarshal(body, &result)
	//fmt.Printf("%+v", string(body))
	fmt.Printf("%+v", result)
	result = filterVhost(&vhosts, result)

	var failed []string
	for _, queue := range result {
		if contains(&qw, queue.Name) && queue.Messages >= treshold {
			failed = append(failed, queue.Name)
		}
	}

	if len(failed) > 0 {
		fmt.Printf("Critical Queues ["+strings.Join(failed, ",")+"] over threshold [%d] | queues=%d", treshold, len(failed))
		os.Exit(ecode)
	}

	fmt.Printf("OK Queues under threshold  [%d] | queues=%d", treshold, len(result))
}

func filterVhost(allowed *[]string, list ResultJSON) ResultJSON {
	if len(*allowed) == 0 {
		return list
	}

	f := list[:0]
	for _, a := range *allowed {
		for _, l := range list {
			if a == l.Vhost {
				f = append(f, l)
			}
		}
	}
	return f
}

func contains(list *[]string, name string) bool {
	if len(*list) == 0 {
		return true
	}

	for _, n := range *list {
		if name == n {
			return true
		}
	}

	return false
}
