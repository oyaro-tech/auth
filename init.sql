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
