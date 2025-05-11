#!/bin/bash

# Initialize script for Mouse Jiggler
# This script sets up a proper Go module environment

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Setting up Mouse Jiggler project...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    echo "Visit https://golang.org/dl/ for installation instructions."
    exit 1
fi

# Create go.mod file if it doesn't exist
if [ ! -f go.mod ]; then
    echo -e "${GREEN}Initializing Go module...${NC}"
    go mod init mouse-jiggler
    
    echo -e "${GREEN}Adding required dependencies to go.mod...${NC}"
    go get github.com/getlantern/systray@latest
    go get github.com/vcaesar/gops@latest  # This is a dependency without OCR
    go get github.com/robotn/gohook@latest  # This is a dependency without OCR
    go get github.com/robotn/xgb@latest    # This is a dependency without OCR
    go get github.com/robotn/xgbutil@latest  # This is a dependency without OCR
    go get golang.org/x/image/bmp@latest   # This is a dependency without OCR
    
    # Get a version of robotgo that works without OCR
    go get github.com/go-vgo/robotgo@v0.100.10
    
    go mod tidy
    
    echo -e "${GREEN}Go module initialized successfully.${NC}"
else
    echo -e "${GREEN}Go module already exists. Updating dependencies...${NC}"
    go get github.com/getlantern/systray@latest
    go get github.com/vcaesar/gops@latest  # This is a dependency without OCR
    go get github.com/robotn/gohook@latest  # This is a dependency without OCR
    go get github.com/robotn/xgb@latest    # This is a dependency without OCR
    go get github.com/robotn/xgbutil@latest  # This is a dependency without OCR
    go get golang.org/x/image/bmp@latest   # This is a dependency without OCR
    
    # Get a version of robotgo that works without OCR
    go get github.com/go-vgo/robotgo@v0.100.10
    
    go mod tidy
    echo -e "${GREEN}Dependencies updated.${NC}"
fi

# Create bin directory
mkdir -p bin

echo -e "${GREEN}Setup complete! You can now run 'make build' to build the application.${NC}"
echo -e "Run the following commands to install system dependencies and build:"
echo -e "  ${YELLOW}make sysdeps${NC} - Install required system dependencies"
echo -e "  ${YELLOW}make build${NC}   - Build the application"
echo -e "  ${YELLOW}make bundle${NC}  - Create a macOS application bundle"