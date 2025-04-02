create table bootstrap_nodes_history
(
    id               integer                  not null,
    version          integer                  not null,
    peer_id          varchar(255)             not null,
    multi_addresses  text[]                   not null,
    protocol_version varchar(20),
    region           varchar(50),
    is_active        boolean,
    valid_from       timestamp with time zone not null,
    valid_to         timestamp with time zone,
    primary key (id, version)
);

alter table bootstrap_nodes_history
    owner to peer;

create index idx_history_peer_ver
    on bootstrap_nodes_history (peer_id asc, version desc);

