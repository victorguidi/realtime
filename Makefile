# Define the compiler and compiler flags
GO := go
GOBUILD := $(GO) build
GOFLAGS := -v
BACKEND_DIR := ./
FRONTEND_DIR := ./frontend

# Define the build targets
all:check create-db create-cert frontend-start backend-build
# all:check create-db create-cert backend-build

# Test if golang and nodejs are installed
check:
	@which go || (echo "Go is not installed" && exit 1)

#Create the db file
create-db:
	@test -f $(BACKEND_DIR)/databases/chat.db || touch $(BACKEND_DIR)/databases/chat.db

#Create the db file
create-cert:
	@test -d $(BACKEND_DIR)/selfCertificate || ./certgen.sh

frontend-start:
	@cd $(FRONTEND_DIR) && pnpm install
	@cd $(FRONTEND_DIR) && pnpm run build
	# @cd $(FRONTEND_DIR) && npm run build
	@test -d $(BACKEND_DIR)/static || mkdir $(BACKEND_DIR)/static
	@cp -R $(FRONTEND_DIR)/static/* $(BACKEND_DIR)/static

backend-build:
	@cd $(BACKEND_DIR) && $(GOBUILD) $(GOFLAGS) -o realtime .

run: all
	@cd $(BACKEND_DIR) && ./realtime

clean:
	@rm -rf $(BACKEND_DIR)/bin/*

.PHONY: all create-db create-cert run clean check frontend-start backend-build
