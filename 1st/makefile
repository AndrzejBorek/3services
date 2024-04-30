EXECUTABLE = 1st
DOCKER_EXPOSED_PORT = 8080
HOST_EXPOSED_PORT = 8080

default:
	run
	
build: clean
	export GODEBUG=gctrace=1
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(EXECUTABLE) ./cmd/$(EXECUTABLE)

run: build
	./$(EXECUTABLE)

test:
	go test -v -race ./...

bench:
	go test -bench=. ./internal/utils -benchmem
	
clean:
	[ ! -f ./$(EXECUTABLE) ] || rm ./$(EXECUTABLE)

dbuild: dclean build
	docker build -t $(EXECUTABLE) .

drun: dbuild 
	docker stop $(EXECUTABLE) || true
	docker rm $(EXECUTABLE) || true
	docker run -p $(HOST_EXPOSED_PORT):$(DOCKER_EXPOSED_PORT) --name $(EXECUTABLE) $(EXECUTABLE):latest

dclean:
	docker rmi $(EXECUTABLE):latest --force

smoke: 
	@echo "Checking if ApacheBench is installed..."
	@command -v ab >/dev/null 2>&1 || (echo "Apache Bench (ab) not installed, skipping performance tests." && exit 1)
	@echo "ApacheBench is installed. Performing smoke test..."
	@ab -n 10000 -c 10 -l http://localhost:$(HOST_EXPOSED_PORT)/generate/json/100000 || echo "Error during benchmarking. Check if container is running"