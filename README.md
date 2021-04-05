# :lock: Ratify

[![GitHub Actions Status](https://github.com/daystram/ratify/actions/workflows/build.yml/badge.svg)](https://github.com/daystram/ratify/actions/workflows/build.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/daystram/ratify)](https://hub.docker.com/r/daystram/ratify)
[![MIT License](https://img.shields.io/github/license/daystram/ratify)](https://github.com/daystram/ratify/blob/master/LICENSE)

**Ratify** is a Central Authentication Service (CAS) implementing OAuth 2.0 and OpenID Connect (OID) protocols, as defined in [RFC 6749](https://tools.ietf.org/html/rfc6749).

## Features

- Implements various authorization flows
- Implements OpenID Connect protocol layer
- Register new applications to use **Ratify**
- Manage registered users (with email verification)
- Multi-factor authentication using Time-based One-Time Password (TOTP)
- Universal login
- User authentication and incident log
- Active session management

## Supported Authorizaton Flows

- Authorization Code
- Authorization Code with PKCE
- _WIP: Client Credentials_

## Client Libraries

Use the following libraries to easily integrate your application with **Ratify**'s authentication service.

- JavaScript/TypeScript: [ratify-client-js](https://github.com/daystram/ratify-client-js)

## Services

The application comes in two parts:

| Name      |  Code Name  | Stack                                                                                                                                                                               |
| --------- | :---------: | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Back-end  | `ratify-be` | [Go](https://golang.org/), [Gin](https://github.com/gin-gonic/gin) + [Gorm](https://github.com/go-gorm/gorm), [PostgreSQL](https://www.postgresql.org/), [Redis](https://redis.io/) |
| Front-end | `ratify-fe` | [TypeScript](https://www.typescriptlang.org/), [Vue.js](https://vuejs.org/)                                                                                                         |

## Develop

`ratify-fe` itself acts as stand-alone application to `ratify-be`, thus it utilizes an access token it self-issued via the _Authorization Code with PKCE_ flow to authenticate users.

### ratify-be

`ratify-be` uses [Go Modules](https://blog.golang.org/using-go-modules) module/dependency manager, hence at least Go 1.11 is required. To ease development, [comstrek/air](https://github.com/cosmtrek/air) is used to live-reload the application. Swagger is used for API documentation, [swaggo/swag](https://github.com/swaggo/swag) is used to generate the docs. Install the tools as documented.

To begin developing, simply enter the sub-directory and run the development server:

```shell
$ cd ratify-be
$ swag init
$ go mod tidy
$ air
```

### ratify-fe

Populate `.env.development` with the required credentials. Use the Client ID that `ratify-be` provides.

To begin developing, simply enter the sub-directory and run the development server:

```shell
$ cd ratify-fe
$ yarn
$ yarn serve
```

## Deploy

Both `ratify-be` and `ratify-fe` are containerized and pushed to [Docker Hub](https://hub.docker.com/r/daystram/ratify). They are tagged based on their application name and version, e.g. `daystram/ratify:be` or `daystram/ratify:be-v1.1.0`.

To run `ratify-be`, run the following:

```shell
$ docker run --name ratify-be --env-file ./.env -p 8080:8080 -d daystram/ratify:be
```

And `ratify-fe` as follows:

```shell
$ docker run --name ratify-fe -p 80:80 -d daystram/ratify:fe
```

### Dependencies

The following are required for `ratify-be` to function properly:

- PostgreSQL
- Redis
- SMTP Server

Their credentials must be provided in the configuration file.

### Helm Chart

To deploy to a Kubernetes cluster, Helm charts could be used. Add the [repository](https://charts.daystram.com):

```shell
$ helm repo add daystram https://charts.daystram.com
$ helm repo update
```

Ensure you have the secrets created for `ratify-be` by providing the secret name in `values.yaml`, or creating the secret from a populated `.env` file (make sure it is on the same namespace as `ratify` installation):

```shell
$ kubectl create secret generic secret-ratify-be --from-env-file=.env
```

And install `ratify`:

```shell
$ helm install ratify daystram/ratify
```

You can override the chart values by providing a `values.yaml` file via the `--values` flag.

Pre-release and development charts are accessible using the `--devel` flag. To isntall the development chart, provide the `--set image.tag=dev` flag, as development images are deployed with the suffix `dev`.

### Docker Compose

For ease of deployment, the following `docker-compose.yml` file can be used to orchestrate the stack deployment:

```yaml
version: "3"
services:
  ratify-be:
    image: daystram/ratify:be
    ports:
      - "8080:8080"
    env_file:
      - /path_to_env_file/.env
    restart: unless-stopped
  ratify-fe:
    image: daystram/ratify:fe
    ports:
      - "80:80"
    restart: unless-stopped
  postgres:
    image: postgres:13.1-alpine
    volumes:
      - /path_to_postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
  redis:
    image: redis:6.0-alpine
    expose:
      - 6379
    volumes:
      - /path_to_redis_data:/data
    restart: unless-stopped
```

### PostgreSQL UUID Extension

UUID support is also required in PostgreSQL. For modern PostgreSQL versions (9.1 and newer), the contrib module `uuid-ossp` can be enabled as follows:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

## License

This project is licensed under the [MIT License](./LICENSE).
