CWD = $(shell pwd)
SVC = post
DB_SVC = post-db
DB_NAME = post_DB
DB_USER = root
DB_PWD = password
DB_VOL_NAME = post-data
NATS_URL = nats-streaming:4223
NET = fishapp-net
PJT_NAME = $(notdir $(PWD))
TEST = $(shell docker inspect $(NET) > /dev/null 2>&1; echo " $$?")

createnet:
	docker network create $(NET)

sqldoc: migrate
	docker run --rm --name tbls --net $(NET) -v $(CWD)/db:/work ezio1119/tbls \
	doc -f -t svg mysql://$(DB_USER):$(DB_PWD)@$(DB_SVC):3306/$(DB_NAME) ./

proto:
	docker run --rm --name protoc -v $(CWD)/pb:/pb -v $(CWD)/schema:/proto ezio1119/protoc \
	-I/proto \
	-I/go/src/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:/pb \
	--validate_out="lang=go:/pb" \
	post.proto event.proto image.proto

cli:
	docker run --rm --name grpc_cli --net $(NET) znly/grpc_cli \
	call $(HOST):50051 $(HOST).PostService.$(m) "$(q)"

waitdb: updb
	docker run --rm --name dockerize --net $(NET) jwilder/dockerize \
	-wait tcp://$(DB_SVC):3306

waitnats:
	docker run --rm --name dockerize --net $(NET) jwilder/dockerize \
	-timeout 20s \
	-wait tcp://$(NATS_URL)

migrate: waitdb
	docker run --rm --name migrate --net $(NET) \
	-v $(CWD)/db/sql:/sql migrate/migrate:latest \
	-path /sql/ -database "mysql://$(DB_USER):$(DB_PWD)@tcp($(DB_SVC):3306)/$(DB_NAME)" ${a}

newsql:
	docker run --rm --name newsql -v $(CWD)/db/sql:/sql \
	migrate/migrate:latest create -ext sql -dir ./sql ${a}

test:
	docker-compose exec $(DB_SVC) sh -c "go test -v -coverprofile=cover.out ./... && \
	go tool cover -html=cover.out -o ./cover.html" && \
	open ./src/cover.html

up: migrate waitnats
	docker-compose up -d $(SVC)

updb:
	docker-compose up -d $(DB_SVC)

build:
	docker-compose build

down:
	docker-compose down

exec:
	docker-compose exec $(SVC) sh

logs:
	docker logs -f --tail 100 $(PJT_NAME)_$(SVC)_1

dblogs:
	docker logs -f --tail 100 $(PJT_NAME)_$(DB_SVC)_1

rmvol: down
	docker volume rm $(PJT_NAME)_$(DB_VOL_NAME)