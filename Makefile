.PHONY: all

APP_BUILD_NAME = SFB
PATH_MAIN_GO = ./cmd/sfb/main.go
OS = linux
SYSD_FILE = /etc/systemd/system/sfb.service

all: clean get gotest build-web build-app

build-web:
	@echo "  >  Building web-components"
	@cd ./web && npm install && npm run build
	
build-app:
	@echo "  >  Building go app"
	@CGO_ENABLED=0 GOOS=$(OS) go build -ldflags "-w" -a -o $(APP_BUILD_NAME) $(PATH_MAIN_GO)
	@rice append -i ./internal/transport/ --exec $(APP_BUILD_NAME)

build-swagger:
	@echo "  >  Building api"
	@swag init  -g .\cmd\sfb\main.go --parseDependency -o api

gotest:
	go test ./...
	
gotestcover:
	go test ./... -cover
	
get:
	@echo "  >  Checking dependencies"
	@go mod download
	@go install $(PATH_MAIN_GO)
	@sudo apt install golang-rice

clean:
	@echo "  >  Сlearing folder"
	@rm -f ./$(APP_BUILD_NAME)
	@rm -rf ./web/build

install:
	@echo "  >  Installing app as service"
	@cp ./$(APP_BUILD_NAME) /usr/bin
	@mkdir -p /usr/bin/sfb_configs/
	@cp ./configs/config.yml /usr/bin/sfb_configs/
	@echo '[Unit]\nDescription=Linux service for shutdown PC from the browser' > $(SYSD_FILE) 
	@echo '[Service]\nType=simple\nUser=root\nExecStart=/usr/bin/$(APP_BUILD_NAME)\nRestart=on-failure' >> $(SYSD_FILE)
	@echo '[Install]\nWantedBy=multi-user.target' >> $(SYSD_FILE)
	@chmod 644 $(SYSD_FILE)
	@systemctl daemon-reload
	@systemctl enable sfb.service
	@systemctl start sfb.service

uninstall:
	@rm -f $(SYSD_FILE)
	@rm -f /usr/bin/SFB
	@rm -rf /usr/bin/sfb_configs
	@systemctl daemon-reload
