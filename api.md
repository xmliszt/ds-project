# E-SAFE API Documentation (v1)

[TOC]

## Base URL

```
localhost:8080/
```

## User

### /login POST [User log in]

Request body:

```json
{
    "username": "user",
    "password": "password"
}
```

Success Response:

200

```json
{
    "success": true
}
```

Failure Response:

400

```json
{
    "success": false,
    "error": "Error message"
}
```

### /register POST [Create a user]

Request body:

```json
{
    "username": "user",
    "password": "password",
    "role": 1
}
```

Success Response:

200

```json
{
    "success": true
}
```

Failure Response:

400

```json
{
    "success": false,
    "error": "Error message"
}
```

## Secret

### api/v1/secret?alias={...} GET [Get a specific secret]

Request params:

```
alias: the alias (name) to the secret
```

Success Response:

200

```json
{
    "success": true,
    "data": [
        {
            "value": "Vict0r1aSecret",
    		"role": 1
        }
    ]
}
```

Failure Response:

400

```json
{
    "success": false,
    "error": "Error message"
}
```

### api/v1/secret PUT [put a secret]

Request body:

```json
{
    "alias": "workspace_1_secret",
    "value": "Vict0r1aSecret",
    "role": 1
}
```

Success Response:

200

```json
{
    "success": true,
    "data": [
        {
            "value": "hashed_Vict0r1aSecret",
            "role": 1
        }
    ]
}
```

Failure Response:

400

```json
{
    "success": false,
    "error": "Error message"
}
```

### api/v1/secret?alias={...} DELETE [Delete a secret]

Request params:

```
alias: the alias (name) to the secret
```

Success Response:

200

```json
{
    "success": true,
}
```

Failure Response:

400

```json
{
    "success": false,
    "error": "Error message"
}
```

### api/v1/secrets GET [Get all secrets under this role]

Success Response:

200

```json
{
    "success": true,
    "data": {
        "role": 1,
        "data": [
            {
                "alias": "workspace_1_secret",
                "value": "hashed_Vict0r1aSecret",
                "role": 1
            },
            {
                "alias": "workspace_2_secret",
                "value": "hashed_Vict5r2aSecret",
                "role": 1
            }
    	],
    }
}
```

Failure Response:

400

```json
{
    "success": false,
    "error": "Error message"
}
```

401

```json
{
    "success": false,
    "error": "Unauthorized"
}
```



