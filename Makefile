# Define the compiler and compiler flags
GO := go
GOBUILD := $(GO) build
GOFLAGS := -v
BACKEND_DIR := ./

# Define the build targets
all:check create-db backend-build

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

run: all
	@cd $(BACKEND_DIR) && ./realtime

clean:
	@rm -rf $(BACKEND_DIR)/bin/*

.PHONY: all create-db backend-build run clean
