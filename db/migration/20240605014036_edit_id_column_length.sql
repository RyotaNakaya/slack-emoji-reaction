-- +goose Up
-- +goose StatementBegin
alter table message_reactions
  MODIFY channel_id VARCHAR(255),
  MODIFY message_id VARCHAR(255),
  MODIFY reaction_name VARCHAR(255),
  MODIFY message_user_id VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table message_reactions
  MODIFY channel_id VARCHAR(20),
  MODIFY message_id VARCHAR(20),
  MODIFY reaction_name VARCHAR(20),
  MODIFY message_user_id VARCHAR(20);
-- +goose StatementEnd
