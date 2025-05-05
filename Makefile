.PHONY: build install update uninstall test clean run

# Binary name and install path
BIN := gomon
BUILD_DIR := bin
INSTALL_DIR := /usr/local/bin

# Build the application
build:
	go build -o $(BUILD_DIR)/$(BIN) ./cmd/gomon

# Install the binary system-wide
install: build
	sudo cp $(BUILD_DIR)/$(BIN) $(INSTALL_DIR)/$(BIN)
	echo "Installed $(BIN) to $(INSTALL_DIR)"

# Update = build + install
update: build install

# Uninstall the binary
uninstall:
	sudo rm -f $(INSTALL_DIR)/$(BIN)
	echo "Uninstalled $(BIN) from $(INSTALL_DIR)"

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)/
	rm -rf tmp/

# Run locally without install
run:
	go run ./cmd/gomon
