# BUILD AND RUN CLI
BINARY_NAME=agent-code

build: clean
	@go build -o $(BINARY_NAME) .

run: build
	@ ./$(BINARY_NAME)

clean:
	@rm -f $(BINARY_NAME)