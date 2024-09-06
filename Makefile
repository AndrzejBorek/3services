include .env
export $(shell sed 's/=.*//' .env)

default: run

build: 
	$(MAKE) -C $(SERVICE1_PATH) build
	$(MAKE) -C $(SERVICE2_PATH) build 
	
run: build
	docker compose up --build 

smoke:
	$(MAKE) -C $(SERVICE1_PATH) smoke
	$(MAKE) -C $(SERVICE2_PATH) smoke