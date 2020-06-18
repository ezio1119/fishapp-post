DC = docker-compose
CURRENT_DIR = $(shell pwd)
API = post

sqldoc:
	docker run --rm --net=fishapp-net -v $(CURRENT_DIR)/db:/work ezio1119/tbls \
	doc -f -t svg mysql://root:password@$(API)-db:3306/$(API)_DB ./

proto:
	docker run --rm -v $(CURRENT_DIR)/pb:/pb -v $(CURRENT_DIR)/schema:/proto ezio1119/protoc \
	-I/proto \
	-I/go/src/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:/pb \
	--validate_out="lang=go:/pb" \
	post.proto event.proto chat.proto image.proto

cli:
	docker run --rm --net=fishapp-net znly/grpc_cli \
	call $(API):50051 $(API).PostService.$(m) "$(q)"

migrate:
	docker run --rm --name migrate --net=fishapp-net \
	-v $(CURRENT_DIR)/db/sql:/sql migrate/migrate:latest \
	-path /sql/ -database "mysql://root:password@tcp($(API)-db:3306)/$(API)_DB" ${a}

newsql:
	docker run --rm -it --name newsql -v $(CURRENT_DIR)/db/sql:/sql \
	migrate/migrate:latest create -ext sql -dir ./sql ${a}

test:
	$(DC) exec $(API) sh -c "go test -v -coverprofile=cover.out ./... && \
	go tool cover -html=cover.out -o ./cover.html" && \
	open ./src/cover.html

up:
	$(DC) up -d

ps:
	$(DC) ps

build:
	$(DC) build

down:
	$(DC) down

stop:
	$(DC) stop

exec:
	$(DC) exec $(API) sh

logs:
	docker logs -f --tail 100 fishapp-post_post_1

dblogs:
	$(DC) logs -f  $(API)-db

push:
	cd schema
	git pull
	git add .
	git commit -m "	$(m)"
	git push