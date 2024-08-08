EXECUTABLE=1st
DOCKER_EXPOSED_PORT=8080
HOST_EXPOSED_PORT=8080
BUILD_PATH=/usr/src/1st
BIN_PATH=$(BUILD_PATH)/cmd
HOST=localhost

TEST_JSON_SIZE=10000
TEST_REQUEST_NUMBER=100000
TEST_CONCURRENCY_LEVEL=10


default:
	run

clean:
	[ ! -f ./$(EXECUTABLE) ] || rm -f ./$(EXECUTABLE)

build: clean
	docker run --rm -v .:$(BUILD_PATH) -w $(BUILD_PATH) \
    	$(GO_IMAGE) \
    	sh -c "CGO_ENABLED=0 GOOS=linux go build  -gcflags="-m"  -a -installsuffix cgo -o $(BUILD_PATH)/$(EXECUTABLE) $(BIN_PATH)/$(EXECUTABLE)"

run: build
	./$(EXECUTABLE)

test:
	go test -v -race ./... 

bench:
	go test -bench=. ./internal/utils -benchmem -cpuprofile cpu.out


dbuild: dclean build
	docker build -t $(EXECUTABLE) .

drun: dbuild
	docker stop $(EXECUTABLE) || true
	docker rm $(EXECUTABLE) || true
	docker run -p $(HOST_EXPOSED_PORT):$(DOCKER_EXPOSED_PORT) --name $(EXECUTABLE) $(EXECUTABLE):latest

dclean:
	docker rmi 3services-service1:latest --force

smoke:
	@echo "Checking if ApacheBench is installed..."
	@command -v ab >/dev/null 2>&1 || (echo "Apache Bench (ab) not installed, skipping performance tests." && exit 1)
	@echo "ApacheBench is installed. Performing smoke test."
	@ab -n ${TEST_REQUEST_NUMBER} -c ${TEST_CONCURRENCY_LEVEL} -l http://$(HOST):$(HOST_EXPOSED_PORT)/generate/json/${TEST_JSON_SIZE} || echo "Error during benchmarking. Check if container is running"
