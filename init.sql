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
    '$2a$10$GH8J3KJOwwclDHOgEs6qZ.KY18HADgN.hHHBao9oTZ13W7ian75Cm' -- qwerty123
);