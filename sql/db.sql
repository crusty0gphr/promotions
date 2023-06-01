begin;

create table if not exists t_promotions (
    id serial primary key,
    key uuid,
    price numeric not null default 0,
    expiration_date text default ''
);

commit;
