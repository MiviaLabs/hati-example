BUILDDIR = ./build
GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run
GOBUILD = env GO111MODULE=on go build

.PHONY: node build dev clean

build:
	make clean && make app
dev:
	@echo "\n> --- run in development mode --"
	HATI_DEBUG=true HATI_DATA_DIR=./build go run ./main.go start $(cmd)
app:
	mkdir -p $(GOBIN)
	go fmt ./... && $(GOBUILD) -ldflags "-w" -o $(GOBIN)/app
	# cd ./cmd/app/ && go fmt ./... && $(GOBUILD) -ldflags "-w" -o ./../../$(GOBIN)/app
	# cp config.example.yml ./build/config.yml
	# cp config.yml ./build/config.yml
	chmod +x $(GOBIN)/app

	@echo "\n> ---"
	@echo "> Build successful. Executable in: \"$(GOBIN)/app\" "
	@echo "> ---\n"
clean:
	rm -rf $(GOBIN)