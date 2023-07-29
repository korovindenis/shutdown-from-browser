.PHONY: all

APP_BUILD_NAME = sfb
PATH_MAIN_GO = ./cmd/shutdown-from-browser/main.go
OS_BUILD = linux
SYSD_FILE = /etc/systemd/system/sfb.service

ifeq ($(OS), Windows_NT)
	OS_BUILD = windows
	APP_BUILD_NAME = sfb.exe
endif

all: clean get build-web build-app

build-web:
	@echo "  >  Building web-components"
	@cd ./web && npm install --legacy-peer-deps && npm run build
	
build-app:
	@echo "  >  Building go app"
	@CGO_ENABLED=0 GOOS=$(OS_BUILD) go build -ldflags "-w" -a -o $(APP_BUILD_NAME) $(PATH_MAIN_GO)

gotest:
	go test ./...
	
gotestcover:
	go test ./... -cover
	
get:
	@echo "  >  Checking dependencies"
	@go mod download
	@go install $(PATH_MAIN_GO)

clean:
	@echo "  >  Clearing folder"
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
