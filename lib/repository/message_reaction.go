package repository

import (
	_ "github.com/go-sql-driver/mysql"
)

type MessageReaction struct {
	ChannelID     string `db:"channel_id"`
	MessageID     string `db:"message_id"`
	ReactionName  string `db:"reaction_name"`
	ReactionCount uint   `db:"reaction_count"`
	MessageTS     uint   `db:"message_ts"`
	YYYYMM        string `db:"yyyymm"`
	CreatedAt     uint   `db:"created_at"`
}

func (m *MessageReaction) save() error {
	tx := DB.MustBegin()
	_, err := tx.NamedExec(`
		INSERT INTO message_reactions (channel_id, message_id, reaction_name, reaction_count, message_ts, yyyymm, created_at)
		VALUES (:channel_id, :message_id, :reaction_name, :reaction_count, :message_ts, :yyyymm, :created_at)`, &m)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
