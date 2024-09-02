CREATE TABLE Users
(
    id            int          not null primary key,
    name          varchar(255) not null,
    age           int CHECK (age >= 18) not null,
    email         varchar(255) not null unique,
    password_hash varchar      not null
);