SERVICE=$(FIRST_IMAGE)
EXECUTABLE=$(SERVICE).bin

# Resolve those variables also
HOST=localhost
HOST_EXPOSED_PORT=$(FIRST_IMAGE_HOST_PORT)

TEST_JSON_SIZE=100
TEST_REQUEST_NUMBER=10000
TEST_CONCURRENCY_LEVEL=10

clean:
	[ ! -f ./$(EXECUTABLE) ] || rm -f ./$(EXECUTABLE)

dclean: clean
	docker rmi $(SERVICE):latest --force

build: clean dclean
	# Building service 1
	docker run --name $(SERVICE)-compile \
		-v .:$(BUILD_PATH):ro \
		-w $(BUILD_PATH) \
    	$(GO_IMAGE) \
    	sh -c "$(GO_BUILD_COMMAND) $(EXECUTABLE_PATH)/$(EXECUTABLE) $(SOURCE_CODE_PATH)/$(SERVICE)"
	docker cp $(SERVICE)-compile:$(EXECUTABLE_PATH)/$(EXECUTABLE) .
	docker rm $(SERVICE)-compile --force


# Testing targets
test:
	go test -v -race -coverprofile=coverage.out ./... 

bench:
	go test -bench=. ./internal/utils -benchmem -cpuprofile cpu.out

smoke:
	@echo "Checking if ApacheBench is installed..."
	@command -v ab >/dev/null 2>&1 || (echo "Apache Bench (ab) not installed, skipping performance tests." && exit 1)
	@echo "ApacheBench is installed. Performing smoke test."
	@echo "Trying to smoke $(SERVICE) on $(HOST):$(HOST_EXPOSED_PORT)"
	@ab -n ${TEST_REQUEST_NUMBER} -c ${TEST_CONCURRENCY_LEVEL} -l http://$(HOST):$(HOST_EXPOSED_PORT)/generate/${TEST_JSON_SIZE} || echo "Error during benchmarking. Check if container is running"
