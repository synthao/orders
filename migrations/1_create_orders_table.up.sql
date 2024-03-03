create table orders
(
    id         bigint generated always as identity primary key,
    status     smallint default 0,
    sum        numeric(10,2) default 0,
    created_at timestamp default current_timestamp NOT NULL,
    updated_at timestamp default current_timestamp NOT NULL
);