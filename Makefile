
.PHONY: build
build: 
	go build -o mogutou main.go router.go

.PHONY: run
run: build
	./mogutou -c conf/

.PHONY: clean
clean:
	rm mogutou

.PHONY: docker
docker:
	CGO_ENABLED=0 GOOS=linux go build -o mogutou main.go router.go
	docker build . -t xuxu123/mogutou:v0.1.0
	rm mogutou
