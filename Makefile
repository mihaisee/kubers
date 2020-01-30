install:
	@go build
	@sudo mv ./kubectl-ext /usr/local/bin