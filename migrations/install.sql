INSERT INTO block_info (id, hash, net_code, prev_id, next_id, size)
VALUES (1, 'aa', 'BTC', 1, 1, 5),
       (2, 'bb', 'LTC', 1, 2, 6),
       (3, 'cc', 'DOGE', 2, 3, 7),
       (4, 'dd', 'BTC', 3, 4, 8);
INSERT INTO tx_info(id, fee, sent_value)
VALUES (1, 1.4, 2),
       (2, 2.4, 3),
       (3, 3.5, 4),
       (4, 2.4, 3),
       (5, 3.5, 4),
       (6, 2.4, 3),
       (7, 3.5, 4),
       (8, 5.1, 5);
INSERT INTO block2tx(block_id, tx_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (1, 4),
       (1, 5),
       (1, 6);