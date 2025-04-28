# cloudru-balancer

This project is a [test task](./task/README.md) for internship to cloud.ru.

## Build and run

### With docker compose

Run

```shell
docker compose up --build -d 
```

to start services with docker compose. In addition to balancer docker compose starts dummy-backend that logs methhod uri and body of incoming requests.

### With Makefile

To build balancer service run (the `balancer` executable will be created):

```shell
make build-balancer
```

To build dummy backend run (the `dummy` executable will be created):

```shell
make build-dummy
```

Run given executables

## Command line arguments

### balancer

`--config`
> Path to config file. Default value is `/etc/cloudru_balancer/balancer.yml`

`--print-config`
> If present prints balancer config.

## Configuration

### balancer

Please see [example](./configs/balancer_config_example.yml)

## Responses

If error occurs while processing request (for example there is no available backends to handle request), balancer responses with `5xx` status code and following body:

```json
{
    "msg": "the error message",
    "code": 503
}
```
