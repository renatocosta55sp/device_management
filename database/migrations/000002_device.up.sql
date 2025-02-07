create sequence domain_event_entry_seq start with 1 increment by 1;

create table domain_event_entry
(
    global_index         bigint       not null,
    aggregate_identifier varchar(255) not null,
    event_identifier     varchar(255) not null unique,
    payload_revision     varchar(255),
    sequence_number      bigint       not null,
    time_stamp           varchar(255) not null,
    type                 varchar(255),
    meta_data            jsonb,
    payload              jsonb          not null,
    primary key (global_index),
    unique (aggregate_identifier, sequence_number)
);

create sequence snapshot_event_entry_seq start with 1 increment by 50;

create table snapshot_event_entry
(
    sequence_number      bigint       not null,
    aggregate_identifier varchar(255) not null,
    event_identifier     varchar(255) not null unique,
    payload_revision     varchar(255),
    payload_type         varchar(255) not null,
    time_stamp           varchar(255) not null,
    type                 varchar(255) not null,
    meta_data            oid,
    payload              oid          not null,
    primary key (sequence_number, aggregate_identifier, type)
);