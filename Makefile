run-publisher: 
		cd order-publish/cmd && go build -v ./
		cd order-publish/cmd && ./cmd

# start-docker:
# 		docker-compose up

# stop-docker:
# 		docker-compose stop

# test:
# 		cd server-subscriber && go test ./cache && go test ./database

run-subscriber:
		cd order-subscribe/cmd && go build ./
		cd order-subscribe/cmd && ./cmd


.PHONY: start-docker stop-docker test run-publisher run-subscriber