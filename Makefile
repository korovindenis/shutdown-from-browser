APP_BUILD_NAME = SFB
PATH_MAIN_GO = ./cmd/sfb/main.go
OS = linux

build: clean get build-web build-app

build-web:
	@echo "  >  Building web-components"
	@cd ./web && npm run build
	
build-app:
	@echo "  >  Building go app"
	@go mod download && CGO_ENABLED=0 GOOS=$(OS) go build -ldflags "-w" -a -o $(APP_BUILD_NAME) $(PATH_MAIN_GO)
	@rice append -i ./server/ --exec $(APP_BUILD_NAME)

build-swagger:
	@echo "  >  Building api"
	@swag init  -g .\cmd\sfb\main.go --parseDependency -o api

get:
	@echo "  >  Checking dependencies"
	@go install $(PATH_MAIN_GO)

clean:
	@echo "  >  Сlearing folder"
	@rm -f ./cmd/sfb/$(APP_BUILD_NAME)
	@rm -rf ./web/build

install:
	@echo "  >  Installing app as service"
	@if [ -d /etc/init.d ]; then cp $(APP_BUILD_NAME) /etc/init.d/ ; fi
	@if [ -d /etc/rc.d ]; then cp $(APP_BUILD_NAME) /etc/rc.d/ ; fi