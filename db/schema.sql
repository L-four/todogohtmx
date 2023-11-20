create table kv
(
    key   TEXT not null
        constraint kv_pk
            primary key,
    value TEXT not null
);

create table lists
(
    id             INTEGER
        primary key,
    title          TEXT,
    show_completed INTEGER,
    sort_by        TEXT
);

create table sessions
(
    id          TEXT              not null,
    data        integer TEXT,
    last_access integer default 0 not null
);

create table todos
(
    id        INTEGER
        primary key,
    title     TEXT,
    completed INTEGER,
    list_id   INTEGER not null
        constraint todos_list_id_fk
            references lists
            on delete cascade
);

