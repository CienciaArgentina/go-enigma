# Enigma

## Introduction
Its main purpose is to provide a service for both authentication and authorization.

## Configuration
You **must** add the environment variables, either this won't work.
The variables that the system needs are:
```shell
    export DB_PASSWORD value
    export DB_HOSTNAME value
```

## Entities

#### User
***
```json
{
  "userId": int,
  "username": string,
  "normalizedUsername": string,
  "passwordHash": string,
  "lockoutEnabled": bool,
  "lockoutEnd": datetime,
  "failedLoginAttempts": int,
  "dateCreated": datetime,
  "dateModified": datetime,
  "modificationMadeByUserId": int,
  "securityToken": string, // this will be used to reset passwords and will be changed everytime a reset is requested
  "verificationToken": string, // used for email and device verification, changed after a request is made
  "dateDeleted": datetime
}
```

#### User emails
***
```json
{
  "userEmailId:" int,
  "userId": int,
  "email": string,
  "normalizedEmail": string,
  "verifiedEmail": bool,
  "verificationDate": datetime,
  "dateCreated": datetime,
  "modificationDate": datetime,
  "modificationMadeByUserId": int,
  "dateDeleted": datetime
}
```

#### User roles
***
```json
{
  "userId": int,
  "roleId": int
}
```

#### Roles
***
```json
{
  "roleId": int,
  "name": string,
  "normalizedName": string
}
```

#### SSO user login
***
```json
{
  "ssoUserLoginId": int,
  "userId": int,
  "providerKey": string,
  "providerDisplayName": string
}
```

#### Role claims
***
```json
{
  "roleClaimId": int,
  "roleId": int,
  "claimType": int,
  "claimValue":  string
}
```









