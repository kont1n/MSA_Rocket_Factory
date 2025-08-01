-- +goose Up
create table if not exists orders (
    order_uuid uuid primary key default gen_random_uuid(),
    user_uuid uuid not null,
    part_uuid uuid[] not null,
    total_price float not null,
    transaction_uuid uuid,
    payment_method varchar(255),
    status varchar(255) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create unique index if not exists idx_orders_order_uuid_user_uuid on orders (order_uuid, user_uuid);

-- +goose Down
drop index if exists idx_orders_order_uuid_user_uuid;
drop table if exists orders;    