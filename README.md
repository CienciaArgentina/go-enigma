# Enigma

# Table of content
- [Introduction](#introduction)
- [Configuration](#Configuration)
- [Entities](#entities)

## Introduction
Its main purpose is to provide a service for both authentication and authorization.

## Must read
- Before any commit (and pull-request) you should use `gofmt`
- Do not add any dependency without asking first (and of course make use of `go mod tidy`)

## Configuration
You **must** add the environment variables, either this won't work.
The variables that the system need are:

```Bash
    export db_password="value"
    export db_hostname="value"
```
