# auth
GIN JWT authentication with token stored in cookies

## Code Example

### `main.go`
```go
package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oyaro-tech/auth"
)

func main() {
	router := gin.New()

	_ = router.SetTrustedProxies(nil)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "UPDATE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/favicon.ico"},
	}))
	router.Use(gin.Recovery())

	auth.RegisterRoutes(router)
	router.GET("/welcome", auth.TokenAuthMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "Welcome admin!")
	})
	router.Run()
}
```

### `.env`
```
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

### `init.sql`
```sql
-- Create database
CREATE DATABASE users;

-- Create user role enum
create type user_role as enum ('user', 'administrator');

-- Create users table
create table if not exists users (
    id SERIAL NOT NULL,
    privileges user_role DEFAULT 'user',
    email varchar(1024) NOT NULL,
    username varchar(64) NOT NULL,
    password varchar(64) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Insert test user
insert into users (privileges, email, username, password)
values (
    'administrator',
    'example@gmail.com',
    'admin',
    '$2a$10$.lWUct/xzfsd8OccI/Fn0ue8aiDMmU/HCffzOTcD8KwsNlldHkOE6' -- qwerty123
);
```

### Run Postgres in Docker and init database
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 --rm postgres
cat init.sql | docker exec -i postgres psql -U postgres
```

### Init package and install requirements
```bash
go mod init github.com/$USER/auth-example
go mod tidy
go get
```

#### Running
```bash
env $(cat .env) go run ./...
```

## Usage

### Try accessing the `/welcome` endpoint
```
curl -v localhost:8080/welcome -X GET
```

```
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /welcome HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.74.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 401 Unauthorized
< Content-Type: application/json; charset=utf-8
< Set-Cookie: access_token=; Path=/; Max-Age=0; HttpOnly; Secure
< Date: Tue, 01 Mar 2022 09:23:09 GMT
< Content-Length: 33
< 
* Connection #0 to host localhost left intact
"no access_token found in cookie"
```

### Login with invalid credentials
```
curl -v localhost:8080/login -X POST -H "Content-Type: application/json" -d '{"username": "admin", "password": "admin"}' 
```

```
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /login HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.74.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 42
> 
* upload completely sent off: 42 out of 42 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 401 Unauthorized
< Content-Type: application/json; charset=utf-8
< Date: Tue, 01 Mar 2022 09:35:36 GMT
< Content-Length: 36
< 
* Connection #0 to host localhost left intact
"Please provide valid login details"
```

### Login with valid credentials
```
curl -v localhost:8080/login -X POST -H "Content-Type: application/json" -d '{"username": "admin", "password": "qwerty123"}'
```

```
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /login HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.74.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 46
> 
* upload completely sent off: 46 out of 46 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Set-Cookie: access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.tJMQeyTTqaFkHbzImAyYcGRzlJYwA04tKZ61OZ3lKqg; Path=/; Max-Age=43200; HttpOnly; Secure
< Date: Tue, 01 Mar 2022 09:31:40 GMT
< Content-Length: 0
<
```

### Try accessing the `/welcome` endpoint with invalid jwt token in cookies
```
curl -v localhost:8080/welcome -b "access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.Fzb932Uj1qCIvi4ggTFMG634mJ-T63lan_G-1tRi9Ek; Path=/; Max-Age=43200; HttpOnly; Secure"
```

```
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /welcome HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.74.0
> Accept: */*
> Cookie: access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.Fzb932Uj1qCIvi4ggTFMG634mJ-T63lan_G-1tRi9Ek; Path=/; Max-Age=43200; HttpOnly; Secure
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 401 Unauthorized
< Content-Type: application/json; charset=utf-8
< Set-Cookie: access_token=; Path=/; Max-Age=0; HttpOnly; Secure
< Date: Tue, 01 Mar 2022 09:45:04 GMT
< Content-Length: 22
< 
* Connection #0 to host localhost left intact
"signature is invalid"
```

### Try accessing the `/welcome` endpoint with valid jwt token in cookies
```
curl -v localhost:8080/welcome -X GET -b "access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.tJMQeyTTqaFkHbzImAyYcGRzlJYwA04tKZ61OZ3lKqg; Path=/; Max-Age=43200; HttpOnly; Secure"
```

```
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /welcome HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.74.0
> Accept: */*
> Cookie: access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.tJMQeyTTqaFkHbzImAyYcGRzlJYwA04tKZ61OZ3lKqg; Path=/; Max-Age=43200; HttpOnly; Secure
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 202 Accepted
< Content-Type: application/json; charset=utf-8
< Date: Tue, 01 Mar 2022 09:37:18 GMT
< Content-Length: 16
< 
* Connection #0 to host localhost left intact
"Welcome admin!"
```

## TODO
- [ ] Create middleware for user with administrator privileges
- [x] Generate ACCESS_SECRET on init
