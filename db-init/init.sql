CREATE
DATABASE users_db;

\c users_db;

CREATE TABLE users
(
    id            varchar(256) primary key,
    login         varchar(256) unique NOT NULL,
    password      varchar(256)        NOT NULL,
    name          varchar(256)        NOT NULL,
    date_of_birth timestamp           NOT NULL
);

CREATE TABLE sessions
(
    id      varchar(256) primary key,
    user_id varchar(256) NOT NULL
);

CREATE
DATABASE documents_db;

\c documents_db;

CREATE
EXTENSION ltree;

CREATE TABLE sections
(
    id       varchar(256) primary key,
    path     ltree,
    name     varchar(256) NOT NULL,
    owner_id varchar(256) NOT NULL
);

CREATE TABLE documents
(
    id         varchar(256) primary key,
    name       varchar(256) NOT NULL,
    text       varchar(256) NOT NULL,
    owner_id   varchar(256) NOT NULL,
    section_id varchar(256) NOT NULL,
    FOREIGN KEY (section_id) REFERENCES sections (id) ON DELETE CASCADE
);
