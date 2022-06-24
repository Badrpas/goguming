



build: build-front build-go

build-go:
	go build .

build-front:
	cd front && npm i && npm run build


run:
	go run .
