CREATE TABLE users
(
    id       bigserial    not null primary key,
    login    varchar(255) not null unique,
    password varchar      not null -- no encrypt --
);

CREATE TABLE articles (
    id bigserial not null primary key,
    title varchar not null unique,
    author bigint not null,
    content varchar not null,
    CONSTRAINT fk_author FOREIGN KEY (author) REFERENCES users(id)
)
