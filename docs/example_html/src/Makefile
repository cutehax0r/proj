# Config
CONTAINER_NAME := name_goes_here
PROD_TAG := $(CONTAINER_NAME):prod

.PHONY help

# Targets
default: help

help:
	@echo "make help		- Display this message"
	@echo "make build		- build the container image"

build:
	podman build -t $(PROD_TAG)
