# Enigma

# Table of content
- [Introduction](#introduction)
- [Configuration](#Configuration)
- [Must read](#must-read)
- [Configuration](#configuration)
- [What's that `GetHandler`?](#whats-that-gethandler)
- [Working directory](#working-directory)
- [cURLs](#curls)
- [TO-DO](#to-do)

## Introduction
Its main purpose is to provide a service for both authentication and authorization.

## Must read
- Before any commit (and pull-request) you should use `gofmt`
- Do not add any dependency without asking first (and of course make use of `go mod tidy`)

## Configuration
You **must** add the environment variables, either this won't work (ex. "development")
The variables that the system need are:

```Bash
    export DB_USERNAME="value" // Database username
    export DB_PASSWORD="value" // Database password
    export DB_HOSTNAME="value" // Database hostname
    export DB_PORT="value" // Database port
    export DB_NAME="value" // Database name
    export PASSWORD_HASHING_KEY="value" // used for salts, SHOULD BE AS PRIVATE AS POSSIBLE
    export ARGON_MEMORY="value"
    export ARGON_ITERATIONS="value"
    export ARGON_PARALLELISM="value"
    export ARGON_SALT_LENGTH="value"
    export ARGON_KEY_LENGTH="value"
```

Any other configuration should be provided in the `config.{SCOPE}.yml` file. Also, please check [working directory](#working-directory).

## What's that `GetHandler`?
Since `gin-gonic` sucks we had to find an easy (and temporal solution) for the `wildcard route conflicts with existing children` error. This happens when you're trying to use
 an RESTful standard (such as /users/{userId}/address), since the httprouter library priorizes speed over the standard.
