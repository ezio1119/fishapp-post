protoc-post:
	protoc \
  -I . \
  -I ${GOPATH}/src \
  -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:. \
  --validate_out="lang=go:." \
	./src/post/controllers/post_grpc/post.proto

protoc-entry:
	protoc \
  -I . \
  -I ${GOPATH}/src \
  -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:. \
  --validate_out="lang=go:." \
	./src/entry/controllers/entry_post_grpc/entry_post.proto