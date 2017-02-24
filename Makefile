# Makefile
# Runs other commands for convenience, because running 8 commands to build is
# way too much to ask

install: build
	@echo "####################"
	@echo "# Installing files #"
	@echo "####################"
	cp ./build/gits $(GOPATH)/bin/gits

uninstall:
	@echo "######################"
	@echo "# Uninstalling files #"
	@echo "######################"
	rm $(GOPATH)/bin/gits

build: clean dependencies
	@echo "##################"
	@echo "# Building files #"
	@echo "##################"
	go build -o ./build/gits ./src

test: lint build
	@echo "#################"
	@echo "# Testing files #"
	@echo "#################"
	cd ./src && go vet -v ./...
	cd ./src && go test ./...

lint: clean dependencies build
	@echo "#################"
	@echo "# Linting files #"
	@echo "#################"
	cd ./src && golint ./...

benchmark: clean dependencies build
	@echo "######################"
	@echo "# Benchmarking files #"
	@echo "######################"
	cd ./src && go test -bench=. -benchmem ./...

dependencies:
	@echo "###########################"
	@echo "# Installing dependencies #"
	@echo "###########################"
	go get github.com/tools/godep
	cd ./src && godep restore

clean:
	@echo "##################"
	@echo "# Cleaning files #"
	@echo "##################"
	rm -rf ./build
	rm -rf ./src/vendor
