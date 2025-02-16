create table employee (
    id serial primary key,
    username text not null,
    password text not null,
    balance integer not null,
    create_time timestamp with time zone default now()
);
create unique index employee_username on employee (username);

create table merch (
    id serial primary key,
    name text not null,
    price integer not null,
    create_time timestamp with time zone default now()
);
create unique index merch_name on merch (name);

insert into merch (name, price) values
    ('t-shirt',	80),
    ('cup',	20),
    ('book',	50),
    ('pen',	10),
    ('powerbank',	200),
    ('hoody',	300),
    ('umbrella',	200),
    ('socks',	10),
    ('wallet',	50),
    ('pink-hoody',	500)
;

create table inventory (
    employee_id integer not null,
    merch_id integer not null,
    quantity integer not null,
    create_time timestamp with time zone default now(),
    primary key (employee_id, merch_id)
);

create table transaction (
    id serial primary key,
    sender_id integer not null,
    receiver_id integer not null,
    amount integer not null,
    transaction_time timestamp with time zone default now()
);
create index transactions_sender_id on transaction (sender_id);
create index transactions_receiver_id on transaction (receiver_id);
