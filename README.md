# Icinga Check Rabbitmq
This is a learning repo for simple icinga checks in golang.

---
## Commands

**Max Command**
Help Output *(./icinga_check_rabbitmq help max)*:
```
Fails with exit code 2 , if a queue is higher then defined treshold

Usage:
  icinga_check_rabbitmq max [flags]

Flags:
  -a, --auth string         basic auth for header authentication. format=username:password
  -e, --exit int            exit code to be used for fail (default 2)
  -h, --help                help for max
  -m, --max int             defines minimum amount of docs that are required when command fails (default 1000)
      --url string          string of url like http://localhost:15672
  -v, --vhosts strings      comma-separated  whitlist for vhosts
  -w, --whitelist strings   comma-separated  whitlist for queue
```

Example Call:
```shell
./icinga_check_rabbitmq max --url http://localhost:15672 -a guest:guest -m 2 -w "hello,world" -v "/blah,/blub,/"
```


---

## Todo 
* doc command overview (max)
* doc param overview for max
* min command

---

## Docs
Good api documentation: [https://cdn.rawgit.com/rabbitmq/rabbitmq-management/v3.7.12/priv/www/api/index.html](https://cdn.rawgit.com/rabbitmq/rabbitmq-management/v3.7.12/priv/www/api/index.html)