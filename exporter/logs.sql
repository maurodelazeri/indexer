drop table logs;
create table logs
(
    log_index bigint,
    transaction_hash varchar(66),
    transaction_index bigint,
    address varchar(42),
    data text,
    topic0 varchar(66),
    topic1 varchar(66),
    topic2 varchar(66),
    topic3 varchar(66),
    block_number bigint,
    block_hash varchar(66)
);

alter table logs add constraint logs_pk primary key (transaction_hash, log_index);

CREATE INDEX block_hash_idx ON logs (block_hash);
CREATE INDEX address_idx ON logs (address);
CREATE INDEX block_number_idx ON logs (block_number);
CREATE INDEX topic0_idx ON logs (topic0);
CREATE INDEX topic1_idx ON logs (topic1);
CREATE INDEX topic2_idx ON logs (topic2);
CREATE INDEX topic3_idx ON logs (topic3);