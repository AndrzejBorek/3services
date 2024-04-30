SERVICE1_PATH = ./1st
SERVICE2_PATH = ./2nd
SERVICE3_PATH = ./3rd


default: runAll

buildAll: cleanAll
	$(MAKE) -C $(SERVICE1_PATH) build
	
cleanAll:
	$(MAKE) -C $(SERVICE1_PATH) clean
	docker compose down 

runAll: buildAll
	docker compose up --build
