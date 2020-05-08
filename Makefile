DC = docker-compose
CURRENT_DIR = $(shell pwd)
API = post

sqldoc:
	docker run --rm --net=fishapp-net -v $(CURRENT_DIR)/db:/work ezio1119/tbls \
	doc -f -t svg mysql://root:password@${API}-db:3306/${API}_DB ./

proto:
	docker run --rm -v $(CURRENT_DIR)/interfaces/controllers/${API}_grpc:$(CURRENT_DIR) \
	-v $(CURRENT_DIR)/schema/${API}:/schema \
	-w $(CURRENT_DIR) ezio1119/protoc \
	-I/schema \
	-I/go/src/github.com/envoyproxy/protoc-gen-validate  \
	--doc_out=markdown,README.md:/schema \
	--go_out=plugins=grpc:. \
	--validate_out="lang=go:." \
	${API}.proto

	docker run --rm -v $(CURRENT_DIR)/schema:/schema -v $(CURRENT_DIR)/interfaces/controllers:/work ezio1119/protoc \
	-I/schema \
	--doc_out=/schema \
	--doc_opt=markdown,README.md \
	--go_out=:/work \
	/schema/event/event.proto

cli:
	docker run --rm --net=fishapp-net znly/grpc_cli \
	call ${API}:50051 ${API}_grpc.PostService.$(m) "$(q)"

migrate:
	docker run --rm --name migrate --net=fishapp-net \
	-v $(CURRENT_DIR)/db/sql:/sql migrate/migrate:latest \
	-path /sql/ -database "mysql://root:password@tcp($(API)-db:3306)/$(API)_DB" down

# seed:
# 	docker run --rm --name seed arey/mysql-client sh


newsql:
	docker run --rm -it --name newsql -v $(CURRENT_DIR)/db/sql:/sql \
	migrate/migrate:latest create -ext sql -dir ./sql ${n}

test:
	$(DC) exec ${API} sh -c "go test -v -coverprofile=cover.out ./... && \
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
	$(DC) exec ${API} sh

logs:
	$(DC) logs -f post-db

dblog:
	$(DC) exec ${API}-db tail -f /var/log/mysql/query.log