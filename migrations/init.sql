create table employee (
    id bigserial primary key,
    username text not null,
    password text not null,
    balance integer not null,
    create_time timestamp with time zone default now()
);
create unique index employee_username on employee (username);
