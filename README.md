# auth
GIN JWT authentication stored in cookies

## Usage example

#### `main.go`
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oyaro-tech/auth"
)

func main() {
	router := gin.Default()
	auth.RegisterRoutes(router)
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
```
-- Create users database
CREATE DATABASE users;

-- Create user role enum
create type user_role as enum ('user', 'administrator');

-- Creation of users table
create table if not exists users (
    id SERIAL NOT NULL,
    privileges user_role DEFAULT 'user',
    email varchar(1024) NOT NULL,
    username varchar(64) NOT NULL,
    password varchar(64) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Insert test users
insert into users (privileges, email, username, password)
values (
    'administrator',
    'example@gmail.com',
    'admin',
    '$2a$10$.lWUct/xzfsd8OccI/Fn0ue8aiDMmU/HCffzOTcD8KwsNlldHkOE6' -- qwerty123
);
```

Run Postgres in Docker and init database
```
docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 --rm postgres
cat init.sql | sudo docker exec -i postgres psql -U postgres
```

Init package and install requirements
```
go mod init github.com/$USER/auth-example
go mod tidy
go get
```

Running example
```
env $(cat .env) go run ./...
```

Use Postman as client
![Postman output](https://github.com/oyaro-tech/auth/example/postman-output.png)