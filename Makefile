SERVICE1_PATH = ./1st
SERVICE2_PATH = ./2nd
SERVICE3_PATH = ./3rd

export GO_IMAGE=golang:1.22.2

default: run

clean:
	$(MAKE) -C $(SERVICE1_PATH) clean
	$(MAKE) -C $(SERVICE1_PATH) dclean
	$(MAKE) -C $(SERVICE2_PATH) clean
	$(MAKE) -C $(SERVICE2_PATH) dclean

build: clean
	$(MAKE) -C $(SERVICE1_PATH) build
	$(MAKE) -C $(SERVICE2_PATH) build
	

run: build
	docker compose up --build --force-recreate
