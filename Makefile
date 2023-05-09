# LOCAL DEV TARGETS #

build:
	@echo "make: Building..."
	@cd src && go build -o ../dist/event-horizon
	@echo "make: Build complete"

run:
	@make build
	@echo "make: Running the microservice"
	./dist/event-horizon

# DOCKER CONTAINER TARGETS #

run-container:
	@echo "make: Building..."
	@go build -o ./event-horizon
	@echo "make: Build complete"
	@echo "make: Running the microservice"
	@WS_RPC="https://some.io" CONTRACTS="a b" ./event-horizon

