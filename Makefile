
BUF_IMG ?= buf

buf-gen:
	docker run --rm  -v `pwd`:/mnt/pwd -w /mnt/pwd $(BUF_IMG) --debug generate --template buf.gen.yaml --path proto

buf:
	docker run --rm -it  -v `pwd`:/mnt/pwd -w /mnt/pwd $(BUF_IMG) $(CMD)

buf.lock:
	@id=$$(docker create $(BUF_IMG)) ; \
	docker cp $$id:/app/$@ $@ ; \
	docker rm -v $$id


swagger-ui:
	 docker-compose -f docker-compose.yaml up

js: static/js/api.js

static/js/api.js: zgen/ts/proto/service.pb.ts
	p=$$PWD ; \
	cd zgen/ts/proto ; $$p/esbuild service.pb.ts --bundle \
	--outfile=$$p/static/js/api.js --global-name=AppAPI

#	--sourcemap --target=chrome58 \
#	--minify --sourcemap --target=chrome58,firefox57,safari11,edge16 \

run:
	go run cmd/hocon/*.go
