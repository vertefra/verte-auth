# verte-auth

## authentication service for multiple projects

## notes about api

### new user signup

From now on the user will be able to create new authentications using the random key generated

`/api/users/signup`

- body request:

```
{
	"email":"verte.fra@gmail.com",
	"password":"verte"
}
```

- body response

```
{
    "createdUser": {
        "ID": 1,
        "CreatedAt": "2020-11-02T16:58:50.679153796-05:00",
        "UpdatedAt": "2020-11-02T16:58:50.679153796-05:00",
        "DeletedAt": "0001-01-01T00:00:00Z",
        "email": "verte.fra@gmail.com",
        "password": "",
        "key": "aQ8{71D7"
    },
    "success": true
}
```

- Notes: key is the random key necessary to forward the authentiacation request for all
  the authentication projects created by the user

### new user login

Once the user is authenticated he can access all the protected routes to create authentication
projects. These first two routes will be used by the client interface.

`/api/users/login`

- body request:

```
{
	"email":"verte.fra@gmail.com",
	"password":"verte"
}
```

- body response

```
{
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6InZlcnRlLmZyYUBnbWFpbC5jb20iLCJleHAiOjE2MDQ0NDE2NDF9.CmDTYMpw890awAhxqaoDXGJeDrwywXfrtyULL2_RNIM",
    "userId": 1
}
```

- Notes: the userID is encrypted in the token's claims and will be removed from the body response

### Create new authentication account

an authentication account is related to a UserID and stores the information of an account, username, hashed password
and api key
`/api/users/:userID/accounts/signup`

- body request:

```
{
	"email":"verte.fra@gmail.com",
	"password":"verte"
}
```

- headers

`key: aQ8{71D7`

- body response

```
{
    "createdAccount": {
        "ID": 1,
        "CreatedAt": "2020-11-02T17:20:44.404432979-05:00",
        "UpdatedAt": "2020-11-02T17:20:44.404432979-05:00",
        "DeletedAt": "0001-01-01T00:00:00Z",
        "username": "new2@gmail.com",
        "password": "",
        "ownerID": 1,
        "key": "aQ8{71D7",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6ImFROHs3MUQ3IiwiZXhwIjoxNjA0NDQyMDQ0fQ.6SLlRCSx8nGfjZv4a_3BixDNv1pYKBr_5-pPUORV_mE"
    },
    "sucess": true
}
```

- Notes: Keys goes in the header with key of "key" and it's the key provided during the first subscription and will
  be used to access all your projects

### Access the authentication account

Will compare the ashed password and will return a token with project username encoded

`/api/users/:userID/accounts/login`

- body request:

```
{
	"email":"verte.fra@gmail.com",
	"password":"verte"
}
```

- headers

`key: aQ8{71D7`

- body response

```
{
    "createdAccount": {
        "ID": 1,
        "CreatedAt": "2020-11-02T17:20:44.404432979-05:00",
        "UpdatedAt": "2020-11-02T17:20:44.404432979-05:00",
        "DeletedAt": "0001-01-01T00:00:00Z",
        "username": "new2@gmail.com",
        "password": "",
        "ownerID": 1,
        "key": "aQ8{71D7",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6ImFROHs3MUQ3IiwiZXhwIjoxNjA0NDQyMDQ0fQ.6SLlRCSx8nGfjZv4a_3BixDNv1pYKBr_5-pPUORV_mE"
    },
    "sucess": true
}
```

- Notes: Keys goes in the header with key of "key" and it's the key provided during the first subscription and will
  be used to access all your projects
