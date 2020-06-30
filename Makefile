OUTPUT?=bin/multiplexer
APP?=multiplexer
GO111MODULE=on
IMAGE_TAG?=
CA_DIR?=certs

LDFLAGS="-s -extldflags '-static'"

all: vendor lint test build

.PHONY: vendor
vendor:
	GO111MODULE=${GO111MODULE} go get ./... && \
	GO111MODULE=${GO111MODULE} go mod tidy && \
	GO111MODULE=${GO111MODULE} go mod vendor

.PHONY: lint
lint:
ifeq (, $(shell which golangci-lint))
	$(error "No golangci-lint in $(PATH). Install it from https://github.com/golangci/golangci-lint")
endif
	golangci-lint run

.PHONY: test
test:
	GO111MODULE=${GO111MODULE} go test -mod vendor ./...

.PHONY: clean
clean:
	rm -rf ${OUTPUT}

.PHONY: build
build: clean
	@echo "+ $@"
	CGO_ENABLED=0 GO111MODULE=${GO111MODULE} go build \
		-mod vendor \
		-tags "netgo std static_all" \
		-ldflags $(LDFLAGS) \
		-o ${OUTPUT} cmd/main.go

.PHONY: certs
certs:
ifeq ("$(wildcard $(CA_DIR)/ca-certificates.crt)","")
	@echo "+ $@"
	@docker run --name ${APP}-certs -d alpine:latest sh -c "apk --update upgrade && apk add ca-certificates && update-ca-certificates"
	@docker wait ${APP}-certs
	@mkdir -p ${CA_DIR}
	@pwd
	@docker cp ${APP}-certs:/etc/ssl/certs/ca-certificates.crt ${CA_DIR}
	@docker rm -f ${APP}-certs
endif

