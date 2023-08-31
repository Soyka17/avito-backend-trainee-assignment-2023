drop table user_segments;
drop table segments;
drop table users;

create table if not exists users
(
    id serial primary key not null unique
);

create table if not exists segments
(
    id   SERIAL primary key not null unique,
    slug varchar(100)       not null unique
    );

create table if not exists user_segments
(
    uid        int       not null references users (id) on delete cascade,
    sid        int       not null references segments (id) on delete cascade,
    begin_date timestamp not null,
    end_date   timestamp
    );