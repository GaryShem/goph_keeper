package postgresql

const createUserTableSQL = `create table if not exists users
(
    username text not null
        constraint users_pk
            primary key,
    password text not null
);`

const createDataTableQSL = `create table if not exists data
(
    username     text not null
        constraint data_users_name_fk
            references users,
    dataname     text not null,
    data_type    text not null,
    text_data    text,
    card_data    text,
    blobsky_data bytea,
    constraint data_pk
        primary key (username, dataname)
);`
