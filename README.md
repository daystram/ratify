# :lock: Ratify

[![Gitlab Pipeline Status](https://img.shields.io/gitlab/pipeline/daystram/ratify/master)](https://gitlab.com/daystram/ratify/-/pipelines)
[![Docker Pulls](https://img.shields.io/docker/pulls/daystram/ratify)](https://hub.docker.com/r/daystram/ratify)
[![MIT License](https://img.shields.io/github/license/daystram/ratify)](https://github.com/daystram/ratify/blob/master/LICENSE)

__Ratify__ is a Central Authentication Service (CAS) implementing OAuth 2.0 and OpenID Connect (OID) protocols, as defined in [RFC 6749](https://tools.ietf.org/html/rfc6749).

## Features
- Implements various authorization flows
- Implements OpenID Connect protocol layer
- Register new applications to use __Ratify__
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
Use the following libraries to easily integrate your application with __Ratify__'s authentication service.
- JavaScript/TypeScript: [ratify-client-js](https://github.com/daystram/ratify-client-js)

## Services
The application comes in two parts:

|Name|Code Name|Stack|
|----|:-------:|-----|
|Back-end|`ratify-be`|[Go](https://golang.org/), [Gin](https://github.com/gin-gonic/gin) + [Gorm](https://github.com/go-gorm/gorm), [PostgreSQL](https://www.postgresql.org/), [Redis](https://redis.io/)|
|Front-end|`ratify-fe`|[TypeScript](https://www.typescriptlang.org/), [Vue.js](https://vuejs.org/)|

## Develop
`ratify-fe` itself acts as stand-alone application to `ratify-be`, thus it utilizes an access token it self-issued via the _Authorization Code with PKCE_ flow to authenticate users.

### ratify-be
`ratify-be` uses [Go Modules](https://blog.golang.org/using-go-modules) module/dependency manager, hence at least Go 1.11 is required. To ease development, [comstrek/air](https://github.com/cosmtrek/air) is used to live-reload the application. Install the tool as documented.

To begin developing, simply enter the sub-directory and run the development server:
```shell
$ cd ratify-be
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
