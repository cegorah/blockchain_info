create table block_info
(
    id        serial primary key,
    hash      varchar NOT NULL,
    net_code  varchar NOT NULL,
    next_id   serial REFERENCES block_info (id),
    prev_id   serial REFERENCES block_info (id),
    size      int default 0,
    timestamp date
);

---- create above / drop below ----

drop table block_info;
