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
```
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
```
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
```
{
  "userId": int,
  "roleId": int
}
```

#### Roles
***
```
{
  "roleId": int,
  "name": string,
  "normalizedName": string
}
```

#### SSO user login
***
```
{
  "ssoUserLoginId": int,
  "userId": int,
  "providerKey": string,
  "providerDisplayName": string
}
```

#### Role claims
***
```
{
  "roleClaimId": int,
  "roleId": int,
  "claimType": int,
  "claimValue":  string
}
```









