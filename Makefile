RELEASE_VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
TIME := $(shell date +"%Y-%m-%d.%H:%M:%S")
DAY := $(shell date +"%Y-%m-%d")
GIT_COUNT 	= $(shell git rev-list --all --count)
VERSION     = $(RELEASE_VERSION)-$(GIT_COUNT)-$(COMMIT)

GOOS ?= $(shell uname -s | tr [:upper:] [:lower:])
GOARCH ?= $(shell go env GOARCH)

FLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Revision=$(COMMIT) -X main.BuildDate=$(TIME)"
IMAGE = registry.cn-shanghai.aliyuncs.com/easyai-io/easyai-platform:$(VERSION)

cli:
	@echo [Build] version=$(VERSION) commit=$(COMMIT) date=$(TIME)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(FLAGS) -o client ./cmd/client/main.go

cli-all:
	@echo TODO

web:
	@echo TODO

ent:
	@echo [Generate] ent
	@cd ./internal/server/dao/ent/ && go run -mod=mod entgo.io/ent/cmd/ent generate --feature intercept,schema/snapshot ./schema

wire:
	@echo [Generate] wire
	@wire gen ./internal/server

lint:
	golangci-lint run ./...

app: lint ent wire swagger
	@echo [Build] version=$(VERSION) commit=$(COMMIT) date=$(TIME)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(FLAGS) -o easyai-platform ./cmd/server/main.go

run: app
	@echo [Run] version=$(VERSION) commit=$(COMMIT) date=$(TIME)
	@ ./easyai-platform web --conf ./configs/app.local.toml



swagger:
	@swag init --parseDependency=false --parseInternal --parseDepth=10  --parseGoList=false --exclude=./internal/server/dao \
 		--generalInfo=./cmd/server/main.go   --output=./internal/server/swagger

image: ent wire web cli-all
	@rm -rf ./easyai-platform-linux-amd64
	@GOOS=linux GOARCH=amd64 go build $(FLAGS) -o easyai-platform-linux-amd64 ./cmd/server/main.go
	@docker build -t $(IMAGE) .

image-push:
	@docker push $(IMAGE)

upgrade:
	@echo "kubectl -n easyai patch deploy  easyai-platform-deploy --type='json' -p='[{\"op\":\"replace\",\"path\":\"/spec/template/spec/containers/0/image\",\"value\":\"$(IMAGE)\"}]'"
	@echo kubectl -n easyai-system get pod -l app=easyai-platform-deploy -w

cloc:
	@docker run --rm -v $(shell pwd):/source-code aldanial/cloc /source-code