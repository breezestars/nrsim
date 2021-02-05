.PHONY: all build api doc clean

all: build

build:
	echo "Not defined yet"
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o 5gc-api-gw main.go

api:
	@echo "\033[32m----- Compiling openapi files -----\033[0m"
	protoc -I=./api/ --go_out=. --go-grpc_out=. api/uegnbsim.proto

cli:
	@echo "\033[32m----- Running CLI -----\033[0m"
	go run cmd/cli/cli.go

master:
	@echo "\033[32m----- Running master -----\033[0m"
	go run cmd/master/master.go cmd/master/client.go cmd/master/server.go

worker:
	@echo "\033[32m----- Running worker -----\033[0m"
	go run cmd/worker/worker.go cmd/worker/client.go cmd/worker/server.go -masterSrvIp=192.168.1.1:50051

test:
	@echo $(filter-out $@,$(MAKECMDGOALS))

doc:
	@echo "\033[32m----- You can view document in -----\033[0m"
	@echo "\033[32m----- http://localhost:6060/pkg/github.com/ng5gc/ -----\033[0m"
	@echo "\033[32m----- Press ctrl+c if you want to exit -----\033[0m"
	godoc -http=:6060

clean:
	@echo "\033[32m----- Clear all environment -----\033[0m"
#	rm -r -f proto/api
#	rm -f *-api-gw