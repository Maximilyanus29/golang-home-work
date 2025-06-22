create table events (
    id serial primary key,
    title text,
    date_start timestamp not null,
    date_end timestamp not null,
    descr text,
    owner_id bigint not null,
    time_before_notify int
);
create index owner_idx on events (owner_id);