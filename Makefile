DC = docker-compose
CURRENT_DIR = $(shell pwd)
API = post

sqldoc:
	docker run --rm --net=api-gateway_default -v $(CURRENT_DIR)/db:/work ezio1119/tbls \
	doc -f -t svg mysql://root:password@${API}-db:3306/${API}_DB ./

proto:
	docker run --rm -v $(CURRENT_DIR)/${API}/controllers/${API}_grpc:$(CURRENT_DIR) \
	-v $(CURRENT_DIR)/schema/${API}:/schema \
	-w $(CURRENT_DIR) thethingsindustries/protoc \
	-I/schema \
	-I/usr/include/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:. \
	--validate_out="lang=go:." \
	--doc_out=markdown,README.md:/schema \
	${API}.proto

cli:
	docker run --rm --net=api-gateway_default namely/grpc-cli \
	call ${API}:50051 ${API}_grpc.PostService.$(m) "$(q)" $(o)

migrate:
	docker run --rm -it --name migrate --net=api-gateway_default \
	-v $(CURRENT_DIR)/db/sql:/sql migrate/migrate:latest \
	-path /sql/ -database "mysql://root:password@tcp($(API)-db:3306)/$(API)_DB" up

test:
	$(DC) exec ${API} sh -c "go test -v -coverprofile=cover.out ./... && \
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
	$(DC) exec ${API} sh

logs:
	$(DC) logs -f --tail 100