# slack-emoji-reaction

slack の emoji reaction を集計して結果を slack にポストするツール

# 環境

- Go 1.16
- MySQL 5.7

# 使い方

## slack app のセッティング

- 対象のワークスペースに slack app を作成する

  - `User Token Scopes` に以下を付与
    - `channels:history`
    - `channels:read`
  - `Bot Token Scopes` に以下を付与
    - `chat:write`

- 払い出された slack app トークンを環境変数にセットする
  - ルートディレクトリに `.env` ファイルを作成し、以下の値をセット
    - `SLACK_BOT_TOKEN`、`SLACK_USER_TOKEN`

## DB マイグレーション

- [goose](https://github.com/pressly/goose) を利用
- Makefile の `goose_` コマンドの内容を自分のローカルの設定に書き換えて実行
  - `make goose_up`

## 集計

- 以下の main 関数を実行すると集計が開始される
  - `cmd/aggregate_reaction/main.go`
- 必要な設定はフラグおよび環境変数で受け取れるようになっているので、適宜セットする
  - `startTime`、`endTime` で集計の対象期間を指定できる（指定しない場合は前月 1 ヶ月間となる）
  - `targetChannelID` を指定すると指定したチャンネルのみを集計する（指定しない場合は全てのパブリックチャンネル）
  - `db~` はデータベース接続情報
- 同じ期間で複数回集計を実行しても多重にデータができることなない
  - `message_ts` と `reaction_name` で複合ユニークになっているため
    - `message_id` はなぜか空の時があるので、仕方なく `message_ts` をキーに利用している

## 結果通知

- 以下の main 関数を実行すると slack に通知される
  - `cmd/post_reaction_result/main.go`
- 必要な設定はフラグおよび環境変数で受け取れるようになっているので、適宜セットする
  - `targetChannelID` には通知先の slack チャンネルの ID を指定
  - 他は集計と同様
