ARG = argument
M = grpc_method
Q = query
DC = docker-compose
CURRENT_DIR = $(shell pwd)

sql:
	docker run --rm -v $(CURRENT_DIR)/migrate/sql:/sql \
	migrate/migrate:latest create -ext sql -dir /sql ${ARG}

sql-doc:
	docker run --rm --net=api-gateway_default -v $(CURRENT_DIR)/db:/work ezio1119/tbls \
	doc -f -t svg mysql://root:password@post-db:3306/post_DB ./

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

cli:
		docker run --rm --net=api-gateway_default namely/grpc-cli \
		call post:50051 post_grpc.PostService.$(M) "$(Q)" $(OPT)

test:
	$(DC) exec post sh -c "go test -v -coverprofile=cover.out -coverpkg=$(ARG) $(ARG) && \
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

exec:
	$(DC) exec post sh

logs:
	$(DC) logs -f