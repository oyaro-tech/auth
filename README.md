# auth
GIN JWT authentication with token stored in cookies

## Usage example

#### `main.go`
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

#### `.env`
```
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

#### `init.sql`
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

Run Postgres in Docker and init database
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 --rm postgres
cat init.sql | sudo docker exec -i postgres psql -U postgres
```

Init package and install requirements
```bash
go mod init github.com/$USER/auth-example
go mod tidy
go get
```

Running example
```bash
env $(cat .env) go run ./...
```

### Use Postman as client
![Postman not_sing_in](https://github.com/oyaro-tech/auth/blob/main/example/not_sing_in.png)
![Postman login](https://github.com/oyaro-tech/auth/blob/main/example/login.png)
![Postman access_granted](https://github.com/oyaro-tech/auth/blob/main/example/access_granted.png)

## TODO
- [ ] Create middleware for user with administrator privileges
