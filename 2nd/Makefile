SERVICE=$(SECOND_IMAGE)
EXECUTABLE=$(SERVICE).bin

HOST=localhost
HOST_EXPOSED_PORT=$(SECOND_IMAGE_HOST_PORT)

TEST_CSV_SIZE=1000
TEST_REQUEST_NUMBER=10000
TEST_CONCURRENCY_LEVEL=10


# Compile targets 

clean:
	[ ! -f ./$(EXECUTABLE) ] || rm -f ./$(EXECUTABLE)

build: clean dclean
	# Building service 2
	docker run --name $(SERVICE)-compile \
		-v .:$(BUILD_PATH):ro \
		-w $(BUILD_PATH) \
    	$(GO_IMAGE) \
    	sh -c "$(GO_BUILD_COMMAND) $(EXECUTABLE_PATH)/$(EXECUTABLE) $(SOURCE_CODE_PATH)/$(SERVICE)"
	docker cp $(SERVICE)-compile:$(EXECUTABLE_PATH)/$(EXECUTABLE) .
	docker rm $(SERVICE)-compile


dclean:
	docker rmi $(SERVICE):latest --force

# Testing targets

smoke: 
	@echo "Checking if ApacheBench is installed..."
	@command -v ab >/dev/null 2>&1 || (echo "Apache Bench (ab) not installed, skipping performance tests." && exit 1)
	@echo "ApacheBench is installed. Performing smoke test..."
	@echo "Trying to smoke $(SERVICE)/1st on $(HOST):$(HOST_EXPOSED_PORT)"
	@ab -n $(TEST_REQUEST_NUMBER) -c $(TEST_CONCURRENCY_LEVEL) -l http://$(HOST):$(HOST_EXPOSED_PORT)/1st/$(TEST_CSV_SIZE) || echo "Error during benchmarking. Check if container is running"
	@echo "Trying to smoke $(SERVICE)/2nd on $(HOST):$(HOST_EXPOSED_PORT)"
	@ab -n $(TEST_REQUEST_NUMBER) -c $(TEST_CONCURRENCY_LEVEL) -l http://$(HOST):$(HOST_EXPOSED_PORT)/2nd/?id || echo "Error during benchmarking. Check if container is running"


test:
	go test -v -race -coverprofile=coverage.out ./...

bench:
	go test -bench=. ./internal/utils -benchmem