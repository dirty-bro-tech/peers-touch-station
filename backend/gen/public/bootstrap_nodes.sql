create table bootstrap_nodes
(
    id                         serial
        primary key,
    peer_id                    varchar(255) not null
        unique,
    multi_addresses            text[]       not null,
    protocol_version           varchar(20),
    region                     varchar(50),
    is_active                  boolean                  default true,
    created_at                 timestamp with time zone default CURRENT_TIMESTAMP,
    last_updated               timestamp with time zone default CURRENT_TIMESTAMP,
    last_successful_connection timestamp with time zone,
    failure_count              integer                  default 0
);

alter table bootstrap_nodes
    owner to peer;

create index idx_nodes_peer_id
    on bootstrap_nodes (peer_id);

