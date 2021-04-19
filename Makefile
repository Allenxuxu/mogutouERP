
.PHONY: build_web
build_web: 
	cd web && npm install && npm run build && cp -rf dist/* ..

build_macos: build_web
	set GOOS=darwin
	set GOARCH=amd64
	mkdir -p _output/macos
	go build -o _output/macos/mogutou main.go router.go
	cp -rf conf _output/macos/

build_linux: build_web
	set GOOS=linux
	set GOARCH=386
	mkdir -p _output/linux
	go build -o _output/linux/mogutou main.go router.go
	cp -rf conf _output/linux/

build_windows: build_web
	set GOOS=windows
	set GOARCH=386
	mkdir -p _output/windows
	go build -o _output/windows/mogutou main.go router.go
	cp -rf conf _output/windows/

.PHONY: run
run: 
	go run main.go router.go -c conf/

.PHONY: docker
docker: build_web
	CGO_ENABLED=0 GOOS=linux go build -o mogutou main.go router.go
	docker build . -t xuxu123/mogutou:v0.1.0
	rm mogutou
