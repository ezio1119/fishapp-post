ARG = ARG
DC = docker-compose
CURRENT_DIR = $(shell pwd)

sql:
	docker run --rm -v $(CURRENT_DIR)/migrate/sql:/sql \
	migrate/migrate:latest create -ext sql -dir /sql ${ARG}

sql-doc:
	docker run --rm --net=api-gateway_default -v $(CURRENT_DIR)/migrate:/work ezio1119/tbls \
	doc -f mysql://root:password@post-db:3306/post_DB ./

proto-post:
	docker run --rm -v $(CURRENT_DIR)/src/post/controllers/post_grpc:$(CURRENT_DIR) -w $(CURRENT_DIR) thethingsindustries/protoc \
	-I. \
	-I/usr/include/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:. \
	--validate_out="lang=go:." \
	--doc_out=markdown,README.md:./ \
	post.proto

proto-entry:
	docker run --rm -v $(CURRENT_DIR)/src/entry/controllers/entry_post_grpc:$(CURRENT_DIR) -w $(CURRENT_DIR) thethingsindustries/protoc \
	-I. \
	-I/usr/include/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:. \
	--validate_out="lang=go:." \
	--doc_out=markdown,README.md:./ \
	entry_post.proto

up:
	$(DC) up

ps:
	$(DC) ps

build:
	$(DC) build

down:
	$(DC) down