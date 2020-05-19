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
	--go_opt=paths=source_relative \
	--go_out=plugins=grpc:/pb \
	--validate_out="lang=go,paths=source_relative:/pb" \
	chat/chat.proto event/event.proto

cli:
	docker run --rm --net=fishapp-net znly/grpc_cli \
	call $(API):50051 $(API)_grpc.PostService.$(m) "$(q)"

migrate:
	docker run --rm --name migrate --net=fishapp-net \
	-v $(CURRENT_DIR)/db/sql:/sql migrate/migrate:latest \
	-path /sql/ -database "mysql://root:password@tcp($(API)-db:3306)/$(API)_DB" ${a}

# seed:
# 	docker run --rm --name seed arey/mysql-client sh


newsql:
	docker run --rm -it --name newsql -v $(CURRENT_DIR)/db/sql:/sql \
	migrate/migrate:latest create -ext sql -dir ./sql ${n}

test:
	$(DC) exec $(API) sh -c "go test -v -coverprofile=cover.out ./... && \
	go tool cover -html=cover.out -o ./cover.html" && \
	open ./src/cover.html

up:
	$(DC) up -d post-db 
ps:
	$(DC) ps

build:
	$(DC) build

down:
	$(DC) stop post-db
	$(DC) rm post-db

exec:
	$(DC) exec $(API) sh

logs:
	$(DC) logs -f post-db

dblog:
	$(DC) exec $(API)-db tail -f -n 100 /var/log/mysql/query.log