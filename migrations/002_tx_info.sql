create table tx_info
(
    id         serial primary key,
    fee        numeric,
    sent_value integer,
    timestamp  time
);

create table block2tx
(
    block_id serial REFERENCES block_info (id),
    tx_id    serial REFERENCES tx_info (id)
);


---- create above / drop below ----

drop table block2tx;
drop table tx_info;
