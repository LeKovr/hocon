
GOGENS_IMG ?= ghcr.io/apisite/gogens:latest

buf.lock:
	@id=$$(docker create $(BUF_IMG)) ; \
	docker cp $$id:/app/$@ $@ ; \
	docker rm -v $$id

gen:
	docker run --rm  -v `pwd`:/mnt/pwd -w /mnt/pwd $(GOGENS_IMG) --debug generate --template buf.gen.yaml --path proto

buf:
	docker run --rm -it  -v `pwd`:/mnt/pwd -w /mnt/pwd $(GOGENS_IMG) $(CMD)

js: static/js/api.js

static/js/api.js: zgen/ts/proto/service.pb.ts
	docker run --rm  -v `pwd`:/mnt/pwd -w /mnt/pwd --entrypoint /go/bin/esbuild $(GOGENS_IMG)  \
	zgen/ts/proto/service.pb.ts --bundle --outfile=/mnt/pwd/static/js/api.js --global-name=AppAPI

#	--sourcemap --target=chrome58 \
#	--minify --sourcemap --target=chrome58,firefox57,safari11,edge16 \

run:
	go run cmd/hocon/*.go
