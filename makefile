SERVICE1_PATH = ./1st
SERVICE2_PATH = ./2nd
SERVICE3_PATH = ./3rd


default: runAll

cleanAll:
	$(MAKE) -C $(SERVICE1_PATH) clean
	$(MAKE) -C $(SERVICE2_PATH) clean
	docker compose down --rmi "all" --remove-orphans

buildAll: cleanAll
	$(MAKE) -C $(SERVICE1_PATH) build
	$(MAKE) -C $(SERVICE2_PATH) build
	

runAll: buildAll
	docker compose up --build --force-recreate
