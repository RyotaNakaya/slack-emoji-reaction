package repository

// func TestSave(t *testing.T) {
// 	DB = PrepareDBConnection("root", "", "localhost", "3306", "slack_reaction_development", 100, 100, "1h")
// 	mr := MessageReaction{
// 		ChannelID:     "cid",
// 		MessageID:     "mid3",
// 		ReactionName:  "rname",
// 		ReactionCount: 12,
// 		MessageTS:     1621868732,
// 		YYYYMM:        "202101",
// 		CreatedAt:     1621868732,
// 	}

// 	err := mr.save()
// 	if err != nil {
// 		t.Fatalf("err: %#v", err)
// 	}
// }

// func TestMessageReactionsSave(t *testing.T) {
// 	DB = PrepareDBConnection("root", "", "localhost", "3306", "slack_reaction_development", 100, 100, "1h")

// 	mrs := MessageReactions{}

// 	for i := 0; i <= 5; i++ {
// 		mr := MessageReaction{
// 			ChannelID:     "cid",
// 			MessageID:     fmt.Sprintf("mid_%d", i),
// 			ReactionName:  "rname",
// 			ReactionCount: 12,
// 			MessageTS:     1621868732,
// 			YYYYMM:        "202101",
// 			CreatedAt:     1621868732,
// 		}
// 		mrs.MessageReactions = append(mrs.MessageReactions, &mr)
// 	}

// 	err := mrs.save()
// 	if err != nil {
// 		t.Fatalf("err: %#v", err)
// 	}
// }
