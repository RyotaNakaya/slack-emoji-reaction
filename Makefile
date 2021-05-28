test:
	go clean -testcache & go test ./...

goose_staus:
	goose -dir ./db/migration/ mysql "root:@/slack_reaction_development?parseTime=true" status

goose_up:
	goose -dir ./db/migration/ mysql "root:@/slack_reaction_development?parseTime=true" up

goose_down:
	goose -dir ./db/migration/ mysql "root:@/slack_reaction_development?parseTime=true" down

dockerbuild:
	docker build -f build/package/Dockerfile.aggregate_reaction  ./ -t aggregate-reaction:latest