# Define the compiler and compiler flags
GO := go
GOBUILD := $(GO) build
GOFLAGS := -v
BACKEND_DIR := ./
FRONTEND_DIR := ./frontend

# Define the build targets
all:check create-db create-cert backend-build frontend-start

# Test if golang and nodejs are installed
check:
	@which go || (echo "Go is not installed" && exit 1)

#Create the db file
create-db:
	@test -f $(BACKEND_DIR)/databases/chat.db || touch $(BACKEND_DIR)/databases/chat.db

#Create the db file
create-cert:
	@test -d $(BACKEND_DIR)/selfCertificate || ./certgen.sh

backend-build:
	@cd $(BACKEND_DIR) && $(GOBUILD) $(GOFLAGS) -o realtime .

frontend-start:
	@cd $(FRONTEND_DIR) && pnpm run dev

run: all
	@cd $(BACKEND_DIR) && ./realtime

clean:
	@rm -rf $(BACKEND_DIR)/bin/*

.PHONY: all create-db create-cert backend-build run clean check frontend-start
