# Bundle and install the app
.PHONY: install
install: bundle
	@echo "Installing application to Applications folder..."
	cp -R $(BINDIR)/$(APPNAME).app /Applications/
	@echo "Installation complete. You can find $(APPNAME) in your Applications folder."
	@echo "Note: You may need to grant accessibility permissions to the app."
	@echo "Go to System Preferences > Security & Privacy > Privacy > Accessibility"
	@echo "and add $(APPNAME).app to the list of allowed applications."# Makefile for macOS ARM Mouse Jiggler

# Variables
APPNAME = MouseJiggler
BINDIR = bin
GOFILES = main.go
GOFLAGS = -ldflags="-s -w" # Strip debugging information to reduce binary size
GO = go
GOGET = go get
GOBUILD = $(GO) build
GOCLEAN = $(GO) clean
GOINSTALL = $(GO) install

# Default target
.PHONY: all
all: deps build

# Create bin directory if it doesn't exist
$(BINDIR):
	mkdir -p $(BINDIR)

# Initialize module and install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@if [ ! -f go.mod ]; then \
		echo "Initializing Go module..."; \
		$(GO) mod init mouse-jiggler; \
	fi
	$(GO) get github.com/getlantern/systray
	$(GO) get github.com/go-vgo/robotgo/mouse
	$(GO) mod tidy

# Build the application
.PHONY: build
build: $(BINDIR) deps
	@echo "Building $(APPNAME) for macOS ARM..."
	GOARCH=arm64 GOOS=darwin $(GOBUILD) $(GOFLAGS) -o $(BINDIR)/$(APPNAME) $(GOFILES)
	@echo "Build complete: $(BINDIR)/$(APPNAME)"

# Build for both Intel and ARM (Universal Binary)
.PHONY: universal
universal: $(BINDIR)
	@echo "Building $(APPNAME) for macOS Intel..."
	GOARCH=amd64 GOOS=darwin $(GOBUILD) $(GOFLAGS) -o $(BINDIR)/$(APPNAME)-intel $(GOFILES)
	@echo "Building $(APPNAME) for macOS ARM..."
	GOARCH=arm64 GOOS=darwin $(GOBUILD) $(GOFLAGS) -o $(BINDIR)/$(APPNAME)-arm $(GOFILES)
	@echo "Creating universal binary..."
	lipo -create -output $(BINDIR)/$(APPNAME) $(BINDIR)/$(APPNAME)-intel $(BINDIR)/$(APPNAME)-arm
	rm $(BINDIR)/$(APPNAME)-intel $(BINDIR)/$(APPNAME)-arm
	@echo "Universal build complete: $(BINDIR)/$(APPNAME)"

# Install system dependencies (macOS)
.PHONY: sysdeps
sysdeps:
	@echo "Installing system dependencies with Homebrew..."
	brew install cmake pkg-config
	brew install automake libtool libpng
	brew install tesseract
	@echo "System dependencies installed successfully."
	@echo "Now run 'make deps' to install Go dependencies."

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BINDIR)

# Create a macOS application bundle
.PHONY: bundle
bundle: build
	@echo "Creating macOS application bundle..."
	mkdir -p $(BINDIR)/$(APPNAME).app/Contents/MacOS
	mkdir -p $(BINDIR)/$(APPNAME).app/Contents/Resources
	cp ./appinfo/*.icns $(BINDIR)/$(APPNAME).app/Contents/Resources/
	cp $(BINDIR)/$(APPNAME) $(BINDIR)/$(APPNAME).app/Contents/MacOS/
	cat ./appinfo/Info.plist | sed 's/mousejiggler/$(APPNAME)/g' > $(BINDIR)/$(APPNAME).app/Contents/Info.plist
	@echo "Bundle created: $(BINDIR)/$(APPNAME).app"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all       - Install dependencies and build the application (default)"
	@echo "  deps      - Install Go dependencies"
	@echo "  sysdeps   - Install system dependencies with Homebrew"
	@echo "  build     - Build the application for macOS ARM"
	@echo "  universal - Create a universal binary (Intel + ARM)"
	@echo "  bundle    - Create a macOS application bundle"
	@echo "  install   - Install the application to /Applications"
	@echo "  clean     - Remove build artifacts"
	@echo "  help      - Display this help message"