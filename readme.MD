# macOS ARM Mouse Jiggler

A simple, efficient mouse jiggler application written in Go that runs in your macOS menu bar. Keeps your system active by simulating small mouse movements at regular intervals.

## Features

- Menu bar application with intuitive controls
- Customizable jiggling radius (2, 5, 10, or 20 pixels)
- Adjustable jiggling intervals (1, 3, 5, or 10 seconds)
- Automatic return to original cursor position
- Full compatibility with Apple Silicon (ARM) Macs
- Also works on Intel Macs (with universal build option)

## Installation

### Prerequisites

- Go 1.18 or newer
- macOS 10.12 or newer
- Administrative access (for system dependencies)

### Quick Start

1. Clone this repository:
   ```
   git clone https://github.com/ViralOne/go-mouse-jiggler.git
   cd go-mouse-jiggler
   ```

2. Run the initialization script:
   ```
   chmod +x initialize.sh
   ./initialize.sh
   ```

3. Install system dependencies:
   ```
   make sysdeps
   ```

4. Build the application:
   ```
   make build
   ```

4.1 Your application is ready in the `bin/MouseJiggler` directory!

5. Create a proper macOS application bundle:
   ```
   make bundle
   ```

6. Your application is ready in the `bin/MouseJiggler.app` directory!

## Manual Setup

If you prefer to set up manually:

1. Initialize Go module:
   ```
   go mod init mouse-jiggler
   ```

2. Install Go dependencies:
   ```
   go get github.com/getlantern/systray
   go get github.com/go-vgo/robotgo
   go mod tidy
   ```

3. Install system dependencies:
   ```
   brew install cmake pkg-config
   brew install automake libtool libpng
   ```

4. Build using the Makefile:
   ```
   make build
   ```

## Usage

1. Launch the application by running the executable
2. Click on the cursor icon in the menu bar
3. Select "Start Jiggling" to begin
4. Customize the radius and interval through the menu options
5. Select "Stop Jiggling" when you're done

## Build Options

The Makefile provides several useful targets:

- `make` or `make all` - Install dependencies and build the application
- `make deps` - Install required Go dependencies
- `make sysdeps` - Install system dependencies via Homebrew
- `make build` - Build the application for macOS ARM
- `make universal` - Create a universal binary (Intel + ARM)
- `make bundle` - Package the application as a proper macOS .app bundle (WIP)
- `make clean` - Remove build artifacts

## Troubleshooting

### Common Issues

1. **Build fails with package errors**
   
   Make sure you've run the initialization script or manually set up the Go module:
   ```
   ./initialize.sh
   ```

2. **System dependency errors**
   
   Make sure you've installed all required system libraries:
   ```
   make sysdeps
   ```

3. **Permission issues when running the app**
   
   macOS may restrict the app. Go to System Preferences > Security & Privacy and allow the app to run.

4. **Accessibility permissions**
   
   The app needs accessibility permissions to control the mouse. You'll be prompted for this when you first run the app.

## Acknowledgments

- [systray](https://github.com/getlantern/systray) for menu bar integration
- [robotgo](https://github.com/go-vgo/robotgo) for mouse control capabilities