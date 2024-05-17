-- +goose Up
create table route_accesses (
    id serial primary key,
    route text not null,
    role user_role not null
);

-- +goose Down
drop table route_accesses;
