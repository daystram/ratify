# :lock: Ratify

[![Gitlab Pipeline Status](https://img.shields.io/gitlab/pipeline/daystram/ratify/master)](https://gitlab.com/daystram/ratify/-/pipelines)
[![Docker Pulls](https://img.shields.io/docker/pulls/daystram/ratify)](https://hub.docker.com/r/daystram/ratify)
[![MIT License](https://img.shields.io/github/license/daystram/ratify)](https://github.com/daystram/ratify/blob/master/LICENSE)

Ratify is a Central Authentication Service (CAS) implementing OAuth 2.0 and OpenID Connect (OIDC) protocols, as defined in [RFC 6749](https://tools.ietf.org/html/rfc6749).

## Features
- Implements various authorization flows
- Implements OpenID Connect protocol layer
- Register new applications to use Ratify
- Manage registered users (with email verification)
- Multi-factor authentication using Time-based One-Time Password (TOTP)
- Universal login
- User authentication and incident log
- _WIP: Active session management_

## Supported Authorizaton Flows
- Authorization Code
- Authorization Code with PKCE
- _WIP: Client Credentials_

## Client Libraries
Use the following libraries to easily integrate your application with Ratify's authentication service.
- JavaScript: [ratify-client-js](https://github.com/daystram/ratify-client-js)

## Develop
Ratify is split into two sub-applications, `ratify-be` (backend) and `ratify-fe` (frontend). `ratify-fe` itself acts as stand-alone application to `ratify-be` and thus utilizes an access token it self-issued via the _Authorization Code with PKCE_ flow to authenticate users.

### ratify-be
`ratify-be` is developed for Go with [Gin](https://github.com/gin-gonic/gin) framework. The project uses [Go Modules](https://blog.golang.org/using-go-modules) module/dependency manager, hence at least Go 1.11 is required. To ease with development, [comstrek/air](https://github.com/cosmtrek/air) is used to live-reload the application. Create and populate the `config.yaml` file based on the provided template.

To begin developing, simply enter the sub-directory and run the development server:
```console
$ cd ratify-be
$ go mod tidy
$ air
```

### ratify-fe
`ratify-fe` is written in TypeScript and uses [Vue](https://github.com/vuejs/vue) framework. The `.env.development` file needs to be populated with the Client ID that `ratify-be` provides. 

To begin developing, simply enter the sub-directory and run the development server:
```console
$ cd ratify-fe
$ yarn
$ yarn serve
```

## Deploy
Both `ratify-be` and `ratify-fe` are containerized and pushed to [Docker Hub](https://hub.docker.com/r/daystram/ratify). They are tagged based on their application name and version, e.g. `daystram/ratify:be` or `daystram/ratify:be-v0.9.1`.

For `ratify-be`, a configuration file has to be bound to the container. Use the provided template [example.config.yaml](./ratify-be/config/example.config.yaml).

To run `ratify-be`, run the following:
```console
$ docker run --name ratify-be -v /path_to_config/config.yaml:/config.yaml:ro -p 8080:8080 -d daystram/ratify:be
```

And `ratify-fe` as follows:
```console
$ docker run --name ratify-fe -p 80:80 -d daystram/ratify:fe
```

### Dependencies
The following are required for `ratify-be` to function properly:
- PostgreSQL
- Redis
- SMTP Server

Their credentials must be provided in the configuration file.

UUID support is also required in PostgreSQL. For modern PostgreSQL versions (9.1 and newer), the contrib module `uuid-ossp` can be enabled as follows:
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### Docker Compose
For ease of deployment, the following `docker-compose.yml` file can be used to orchestrate the stack deployment:
```yaml
version: "3"
services:
  ratify-fe:
    image: daystram/ratify:fe
    ports:
      - 80:80
    restart: unless-stopped
  ratify-be:
    image: daystram/ratify:be
    ports:
      - 8080:8080
    volumes:
      - /path_to_config/config.yaml:/config.yaml
    restart: unless-stopped
  postgres:
    image: postgres:13.1-alpine
    volumes:
      - /path_to_pg_data:/var/lib/postgresql/data
    restart: unless-stopped
  redis:
    image: redis:6.0-alpine
    volumes:
      - /path_to_rd_data:/data
    restart: unless-stopped
```

## License
This project is licensed under the [MIT License](./LICENSE).
 