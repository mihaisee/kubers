.PHONY: install
install:
	@echo "Building..."
	@go build
	@echo "Moving binary to /usr/local/bin"
	@sudo mv kubers /usr/local/bin
