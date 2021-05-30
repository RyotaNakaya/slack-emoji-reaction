package repository

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MessageReaction struct {
	ChannelID     string `db:"channel_id"`
	MessageID     string `db:"message_id"`
	ReactionName  string `db:"reaction_name"`
	ReactionCount uint   `db:"reaction_count"`
	MessageUserID string `db:"message_user_id"`
	MessageTSNano string `db:"message_ts_nano"`
	MessageTS     uint   `db:"message_ts"`
	YYYYMM        string `db:"yyyymm"`
	CreatedAt     uint   `db:"created_at"`
}

func (m *MessageReaction) save() error {
	tx := DB.MustBegin()
	// _, err := tx.NamedExec(`
	// 	INSERT INTO message_reactions (channel_id, message_id, reaction_name, reaction_count, message_ts, yyyymm, created_at)
	// 	VALUES (:channel_id, :message_id, :reaction_name, :reaction_count, :message_ts, :yyyymm, :created_at)`, &m)
	_, err := tx.Exec(buildUpsertQuery([]*MessageReaction{m}))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// バルクインサート用
type MessageReactions struct {
	MessageReactions []*MessageReaction
}

func (m *MessageReactions) Save() error {
	if len(m.MessageReactions) == 0 {
		return nil
	}

	tx := DB.MustBegin()
	// _, err := tx.NamedExec(`
	// 	INSERT INTO message_reactions (channel_id, message_id, reaction_name, reaction_count, message_ts, yyyymm, created_at)
	// 	VALUES (:channel_id, :message_id, :reaction_name, :reaction_count, :message_ts, :yyyymm, :created_at)`, m.MessageReactions)
	_, err := tx.Exec(buildUpsertQuery(m.MessageReactions))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// sqlx の NamedExec で upsert 文が作れなさそうだったので、頑張って文字列連結してクエリを作る関数
func buildUpsertQuery(m []*MessageReaction) string {
	q := `
	INSERT INTO message_reactions
		(channel_id, message_id, reaction_name, reaction_count, message_user_id, message_ts_nano, message_ts, yyyymm, created_at)
	VALUES 
	`

	s := []string{}
	for _, v := range m {
		s = append(s, fmt.Sprintf("('%v', '%v', '%v', %v, '%v', %v, '%v', '%v', %v)",
			v.ChannelID, v.MessageID, v.ReactionName, v.ReactionCount, v.MessageUserID, v.MessageTSNano, v.MessageTS, v.YYYYMM, v.CreatedAt))
	}

	// duplicate entry の時は一応 reaction_count だけ更新する
	dup := ` ON DUPLICATE KEY UPDATE reaction_count = VALUES(reaction_count)`

	return q + strings.Join(s, ",") + dup
}
