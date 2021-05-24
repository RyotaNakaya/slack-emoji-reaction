-- +goose Up
-- +goose StatementBegin
create table message_reactions(
    channel_id varchar(20) not null,
    message_id varchar(20) not null,
    reaction_name varchar(20) not null,
    reaction_count int(10) not null,
    message_ts int(11) unsigned not null, -- unix timestamp
    yyyymm char(6) not null, -- YYYYMM
    created_at int(11) unsigned not null -- unix timestamp

    , primary key (message_id, reaction_name)
    , index (channel_id)
    , index (yyyymm)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists
    message_reactions;
-- +goose StatementEnd