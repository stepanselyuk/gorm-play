VERSION ?= $(shell git describe --tags --always)

IMAGE = pocketmedia/alpine
PKG = github.com/stepanselyuk/gorm-play
DOCKER_CONTAINERS = $(docker ps -a -q)

LDFLAGS = "-w -X main.Version=$(VERSION)"

OS ?= linux
ARCH ?= amd64

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/main -a -tags netgo -ldflags $(LDFLAGS)

test:
	go test ./...

start:
	if [ ! -f "bin/main" ]; then make build; fi
	docker-compose up -d

stop:
	docker-compose stop

logs:
	docker logs -f gorm-playroom

rebuild:
	make build
	# stop and remove old containers
	docker-compose down
	# build and start containers
	docker-compose up -d

clean:
	#stop and remove old containers including volumes (mysql)
	docker-compose down --volumes

#npm install -g multi-file-swagger
#generate:
#	sh scripts/generate.sh