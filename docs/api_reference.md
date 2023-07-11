# Go Service

This project is made for assignments for a company. To make an example of an OAuth 2.0 service that focuses on performance, maintainability, and scalability.

## Version: 1.0

**License:** [MIT](https://github.com/ntsd/go-service/blob/main/LICENSE)

### Security

**BasicAuth**

| basic | _Basic_ |
| ----- | ------- |

**OAuth2Application**

| oauth2    | _OAuth 2.0_     |
| --------- | --------------- |
| Flow      | application     |
| Token URL | /v1/oauth/token |

### /oauth/clients

#### POST

##### Description:

create a new client

##### Parameters

| Name   | Located in | Description             | Required | Schema                                                  |
| ------ | ---------- | ----------------------- | -------- | ------------------------------------------------------- |
| client | body       | JSON body of the client | Yes      | [handlers.CreateClientBody](#handlers.CreateClientBody) |

##### Responses

| Code | Description           | Schema                                                  |
| ---- | --------------------- | ------------------------------------------------------- |
| 200  | OK                    | [handlers.CreateClientBody](#handlers.CreateClientBody) |
| 400  | Bad Request           | [handlers.errorResponse](#handlers.errorResponse)       |
| 422  | Unprocessable Entity  | [handlers.errorResponse](#handlers.errorResponse)       |
| 500  | Internal Server Error | [handlers.errorResponse](#handlers.errorResponse)       |

### /oauth/token

#### POST

##### Description:

OAuth2 authentication only support Client Credentials grant type. required `client_id` and `client_secret` on Basic Authentication to.

##### Parameters

| Name       | Located in | Description                                                | Required | Schema |
| ---------- | ---------- | ---------------------------------------------------------- | -------- | ------ |
| grant_type | formData   | The grant_type parameter must must be `client_credentials` | Yes      | string |

##### Responses

| Code | Description           | Schema                                                        |
| ---- | --------------------- | ------------------------------------------------------------- |
| 200  | OK                    | [handlers.AccessTokenResponse](#handlers.AccessTokenResponse) |
| 400  | Bad Request           | [handlers.errorResponse](#handlers.errorResponse)             |
| 401  | Unauthorized          | [handlers.errorResponse](#handlers.errorResponse)             |
| 500  | Internal Server Error | [handlers.errorResponse](#handlers.errorResponse)             |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BasicAuth       |        |

### /users

#### GET

##### Description:

List users.

##### Parameters

| Name   | Located in | Description                                                                  | Required | Schema  |
| ------ | ---------- | ---------------------------------------------------------------------------- | -------- | ------- |
| offset | query      | pagination offset, default is `0`                                            | No       | integer |
| limit  | query      | pagination limit, default is `100`. If more than `100` will be set as `100`. | No       | integer |
| name   | query      | filter name by partial text search                                           | No       | string  |

##### Responses

| Code | Description           | Schema                                            |
| ---- | --------------------- | ------------------------------------------------- |
| 200  | OK                    | [ [models.User](#models.User) ]                   |
| 400  | Bad Request           | [handlers.errorResponse](#handlers.errorResponse) |
| 401  | Unauthorized          | [handlers.errorResponse](#handlers.errorResponse) |
| 500  | Internal Server Error | [handlers.errorResponse](#handlers.errorResponse) |

##### Security

| Security Schema   | Scopes |
| ----------------- | ------ |
| OAuth2Application |        |

#### POST

##### Description:

create a new user

##### Parameters

| Name | Located in | Description           | Required | Schema                                              |
| ---- | ---------- | --------------------- | -------- | --------------------------------------------------- |
| user | body       | JSON body of the user | Yes      | [handlers.CreateUserBody](#handlers.CreateUserBody) |

##### Responses

| Code | Description           | Schema                                            |
| ---- | --------------------- | ------------------------------------------------- |
| 200  | OK                    | [models.User](#models.User)                       |
| 400  | Bad Request           | [handlers.errorResponse](#handlers.errorResponse) |
| 401  | Unauthorized          | [handlers.errorResponse](#handlers.errorResponse) |
| 422  | Unprocessable Entity  | [handlers.errorResponse](#handlers.errorResponse) |
| 500  | Internal Server Error | [handlers.errorResponse](#handlers.errorResponse) |

##### Security

| Security Schema   | Scopes |
| ----------------- | ------ |
| OAuth2Application |        |

### /users/{id}

#### GET

##### Description:

Get user by id

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id   | path       | User ID     | Yes      | string |

##### Responses

| Code | Description           | Schema                                            |
| ---- | --------------------- | ------------------------------------------------- |
| 200  | OK                    | [models.User](#models.User)                       |
| 400  | Bad Request           | [handlers.errorResponse](#handlers.errorResponse) |
| 401  | Unauthorized          | [handlers.errorResponse](#handlers.errorResponse) |
| 404  | Not Found             | [handlers.errorResponse](#handlers.errorResponse) |
| 500  | Internal Server Error | [handlers.errorResponse](#handlers.errorResponse) |

##### Security

| Security Schema   | Scopes |
| ----------------- | ------ |
| OAuth2Application |        |

#### PUT

##### Description:

update a user, it will not insert if not existing.

##### Parameters

| Name | Located in | Description           | Required | Schema                                              |
| ---- | ---------- | --------------------- | -------- | --------------------------------------------------- |
| id   | path       | User ID               | Yes      | string                                              |
| user | body       | JSON body of the user | Yes      | [handlers.UpdateUserBody](#handlers.UpdateUserBody) |

##### Responses

| Code | Description           | Schema                                            |
| ---- | --------------------- | ------------------------------------------------- |
| 200  | OK                    | [models.User](#models.User)                       |
| 400  | Bad Request           | [handlers.errorResponse](#handlers.errorResponse) |
| 401  | Unauthorized          | [handlers.errorResponse](#handlers.errorResponse) |
| 422  | Unprocessable Entity  | [handlers.errorResponse](#handlers.errorResponse) |
| 500  | Internal Server Error | [handlers.errorResponse](#handlers.errorResponse) |

##### Security

| Security Schema   | Scopes |
| ----------------- | ------ |
| OAuth2Application |        |

### Models

#### handlers.AccessTokenResponse

Access Token response body

| Name         | Type    | Description | Required |
| ------------ | ------- | ----------- | -------- |
| access_token | string  |             | No       |
| expires_in   | integer |             | No       |
| token_type   | string  |             | No       |

#### handlers.CreateClientBody

| Name          | Type   | Description | Required |
| ------------- | ------ | ----------- | -------- |
| client_id     | string |             | Yes      |
| client_secret | string |             | Yes      |

#### handlers.CreateUserBody

| Name  | Type   | Description | Required |
| ----- | ------ | ----------- | -------- |
| email | string |             | Yes      |
| name  | string |             | Yes      |

#### handlers.UpdateUserBody

| Name  | Type   | Description | Required |
| ----- | ------ | ----------- | -------- |
| email | string |             | Yes      |
| name  | string |             | Yes      |

#### handlers.errorResponse

Common error response.

| Name    | Type   | Description | Required |
| ------- | ------ | ----------- | -------- |
| message | string |             | No       |

#### models.User

User information includes id, email, and name.

| Name      | Type   | Description | Required |
| --------- | ------ | ----------- | -------- |
| createdAt | string |             | No       |
| email     | string |             | No       |
| id        | string |             | No       |
| name      | string |             | No       |
| updatedAt | string |             | No       |
