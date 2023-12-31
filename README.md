# Go Service

An example of Go API and Kubernetes config focusing on performance, high availability, and scalability.

## Table of Contents

- [Go Service](#go-service)
  - [Table of Contents](#table-of-contents)
  - [Project Layout](#project-layout)
  - [API Reference](#api-reference)
    - [Swagger](#swagger)
  - [Web Framework](#web-framework)
    - [Why Fiber and FastHTTP?](#why-fiber-and-fasthttp)
    - [JSON Serializing](#json-serializing)
  - [Database](#database)
    - [Why PostgreSQL and Gorm?](#why-postgresql-and-gorm)
    - [Caches Prepared Statement](#caches-prepared-statement)
    - [Index Hints](#index-hints)
    - [Trigram (pg\_trgm)](#trigram-pg_trgm)
    - [Connection pool](#connection-pool)
  - [Scalibility](#scalibility)
    - [Kubernetes Horizontal Pod Autoscaler (HPA)](#kubernetes-horizontal-pod-autoscaler-hpa)
    - [Kubernetes StatefulSets](#kubernetes-statefulsets)
  - [Environment Variables](#environment-variables)
  - [Development](#development)
    - [Install Go Packages](#install-go-packages)
    - [Build Docker images](#build-docker-images)
    - [Run Postgres](#run-postgres)
    - [Run Migration](#run-migration)
    - [Run Server](#run-server)
  - [Unit Testing](#unit-testing)
    - [GoMock](#gomock)
    - [Run Tests](#run-tests)
  - [Deployment](#deployment)
    - [Load Testing](#load-testing)
  - [References](#references)

## Project Layout

The project layout is following <https://github.com/golang-standards/project-layout> a non official stadard.

```
├── cmd
│   ├── migrate                 - main file for migration
│   │   └── migrations          - migration files
│   └── service                 - main file for service
├── deployments                 - Docker, keys, and Kubernetes.
├── docs                        - documentation and swagger
└── internal                    - internal packages will not export.
    ├── config                  - configuration and environment
    ├── crypto                  - cryptography package
    │   └── mock_jwt
    ├── handlers                - server handlers and endpoints
    ├── models                  - database models
    └── storage                 - database storage and repositories
        └── mock_storage
```

## API Reference

[Link to API Reference](/docs/api_reference.md)

### Swagger

The swagger generates by <https://github.com/swaggo/swag> which allows generating the swagger schema from comments.

on the default port, the swagger URL will be `http://localhost:8080/swagger/`.

To generate swagger

```sh
# Install swag cli
go get -u github.com/swaggo/swag/cmd/swag
# Generate swagger
make gen-swagger
```

* This API only supports [OAuth 2 Client Credentials](https://www.oauth.com/oauth2-servers/access-tokens/client-credentials/) grant types.

## Web Framework

### Why Fiber and FastHTTP?

[Fiber](https://github.com/gofiber/fiber) is a web framework built on top of [FastHTTP](https://github.com/valyala/fasthttp), the fastest HTTP engine for Go.

1. More performance, because they develop their own HTTP and serializer but it's a trade off with non built-in `net/http`. This makes the framework not compatible with many HTTP libraries such as Gin and Gorilla Mux.

2. Zero/less memory allocation, `*fiber.Ctx` are not immutable so it can be reused anywhere during request.

3. Fast HTTP has its own [Workerpool](https://github.com/valyala/fasthttp/blob/master/workerpool.go) instead of creating a goroutine for every request like other frameworks do. Because they say the Go routine is too expensive.

4. Prefork Support, Preforking makes use of single Go processes but will load balance connections on the OS level by [SO_REUSEPORT](https://lwn.net/Articles/542629/) basically for running multiple servers using the same port. But when using Kubernetes they prefer to run a separate server instance per CPU core with `GOMAXPROCS=1` and `Prefork` to False to spawn pods by auto-scaling instead.

### JSON Serializing

[Sonic](https://github.com/bytedance/sonic) is a fast JSON Serializer, much faster than the built-in from my texting. The library does not enable HTML escaping and SortKeys features by default if we enable it will lose some performance. The downside is they say it may not support well on all environments.

## Database

### Why PostgreSQL and Gorm?

PostgreSQL is an open-source database with rich features, distribution, scalability, and high availability.

Gorm is an open-source ORM library for Go. Gorm provides a simple and efficient way to interact with databases.

### Caches Prepared Statement

Gorm can prepared statement when executing any SQL and caches them to speed up future calls.

### Index Hints

GORM provides optimizer/index/comment hints support.

example

```go
db.Clauses(hints.UseIndex("users_name_trgm_idx")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`users_name_trgm_idx`)
```

### Trigram (pg_trgm)

For Partial/Full Text Seach, PostgreSQL has built-in [pg_trgm](https://www.postgresql.org/docs/current/pgtrgm.html) module for determining the similarity of alphanumeric text based on trigram matching. The pg_trgm module support GiST and GIN index for text columns.

By using trigrams, we can compare similar trigrams using `SIMILARITY` or `%` operators.

example `SELEC * FROM users WHERE SIMILARITY(name,'John') > 0.4 ;`

### Connection pool

For a bettter Database connection may considor the connection pool to keep the connection opening and ready for the query.

## Scalibility

### Kubernetes Horizontal Pod Autoscaler (HPA)

[HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/) allows automatically adjusting the number of pods based on resource utilization.

Parameters

`--cpu-percent=50` the HPA controller will increase and decrease the number of replicas to maintain an average CPU utilization across all Pods of 50%.

`kubectl autoscale deployment go-service-hpa --cpu-percent=50 --min=1 --max=10`

check HPA status

```bash
kubectl get hpa
# or force mode
kubectl get hpa go-service-hpa --watch
```

Or using the yaml file

`kubectl apply -f ./deployments/hpa.yaml`

### Kubernetes StatefulSets

Deploying a Database on K8S might be too complex for auto-scaling. [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) allows creating clusters/replicas sets of PostgreSQL by a specific number of instances.

## Environment Variables

| Key               | Description                                                                   | Example                                       |
| ----------------- | ----------------------------------------------------------------------------- | --------------------------------------------- |
| APP_PORT          | port number, default: `8080`                                                  | 8080                                          |
| POSTGRES_URL      | **Required** PostgreSQL URL                                                   | postgresql://user:pass@localhost:5432/db_name |
| POSTGRES_READ_URL | Postgres URL for read queries, default will be `POSTGRES_URL` if not defined. | postgresql://user:pass@db-slave:5432/db_name  |
| DEV_MODE          | is developement mode, default: `false`                                        | true                                          |
| PREFORK           | is Prefork mode, default: `true`                                              | false                                         |
| ES256_PRIVATE_KEY | **Required** path to ES256 private key                                        | ./deployments/ec_private.pem                  |
| ES256_PUBLIC_KEY  | **Required** path to ES256 public key                                         | ./deployments/ec_public.pem                   |
| HASH_SALT         | hash salt, default: `change_this_salt`                                        | salt                                          |

To make it easier, it has an example Dotenv file to make a simple run for testing (not recommend on prod).

Create `.env` file from `.env.example`

`cp .env.example .env`

## Development

### Install Go Packages

If you want to run without docker you need to install [Go](https://go.dev/doc/install) (version > 1.8) on your local and install packages by:

```sh
go mod download
```

### Build Docker images

```sh
docker-compose build --no-cache
# or
make docker-build
```

### Run Postgres

```sh
docker-compose up service-postgres
# or
make docker-postgres
```

### Run Migration

Required to run migrate database before running the server.

```sh
go run ./cmd/migrate/migrate.go -mode=up
# or
docker-compose up service-migrate
```

For migrate down

`go run ./cmd/migrate/migrate.go -mode=down`

### Run Server

```sh
go run ./cmd/service/main.go
# or
docker-compose up service-app
# or
make docker-up
```

## Unit Testing

### GoMock

[gomock](https://github.com/golang/mock) is a mocking framework for Go.

Install gomock

`GO111MODULE=on go get github.com/golang/mock/mockgen@v1.6.0`

Mockgen

`make mockgen`

### Run Tests

```sh
make test
# or
go test ./...
```

## Deployment

To deploy we using Kubernetes and Kubectl.

```sh
# apply master postgres
kubectl apply -f ./deployments/postgres/master.yaml
# wait until the master start, apply slave postgres
kubectl apply -f ./deployments/postgres/slave.yaml

# apply application
kubectl apply -f ./deployments
```

Pod Forwarding for testing, \*Note this is not load balancing because the port-forwarding only pick the first pod.

`kubectl port-forward service/go-service-service 8080:8080`

### Load Testing

If you run the service in Local you can just use [k6](https://github.com/grafana/k6)

Update the host inside the `./deployments/loadtest/loadtest.js` file then `k6 run ./deployments/loadtest/loadtest.js`.

If you run the service in Kubernetes, I prefer to use [k6-operator](https://github.com/grafana/k6-operator) to run inside Kubernetes instead. Read <https://k6.io/blog/running-distributed-tests-on-k8s/> for installation.

1. Create configmap from file

`kubectl create configmap loadtest --from-file ./deployments/loadtest/loadtest.js --dry-run=client -o yaml > ./deployments/loadtest/configmap.yaml`

2. Run K6 Operator

`kubectl apply -f ./deployments/loadtest/`

To check resources used, run `kubectl top pods`

To check HPA status, run `kubectl get hpa`

## References

<https://github.com/stefanprodan/podinfo>

<https://earthly.dev/blog/optimize-golang-for-kubernetes/>
