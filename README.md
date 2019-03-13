# Icinga Check Rabbitmq
This is a learning repo for simple icinga checks in golang.

## Usage
```shell
    ./icinga_check_rabbitmq max --url http://localhost:15672 -a guest:guest
```

```shell
    ./icinga_check_rabbitmq max --url http://localhost:15672 -a guest:guest -m 2 -w "hello,world" -v "/blah,/blub,/"
```

## Todo 
* doc command overview (max)
* doc param overview for max
* min command
* whitelist queues
* implement vhost filter or direct asking for a vhost

## Docs
Good api documentation: [https://cdn.rawgit.com/rabbitmq/rabbitmq-management/v3.7.12/priv/www/api/index.html](https://cdn.rawgit.com/rabbitmq/rabbitmq-management/v3.7.12/priv/www/api/index.html)