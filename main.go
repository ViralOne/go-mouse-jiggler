package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	"mouse-jiggler/assets/icon"
)

// Config holds the application configuration
type Config struct {
	isJiggling      bool
	jigglingRadius  int
	jigglingInterval time.Duration
	originalX       int
	originalY       int
}

var config = Config{
	isJiggling:      false,
	jigglingRadius:  5,
	jigglingInterval: 3 * time.Second,
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Set up systray - this will immediately create the tray icon
	// without showing any application window
	systray.Run(onReady, onExit)
}

func onReady() {
	// Set the icon and tooltip
	systray.SetIcon(getIcon())
	systray.SetTooltip("Mouse Jiggler - Keep your system active")

	// Menu items - Status first
	mStatus := systray.AddMenuItem("Status: Inactive", "Current status")
	mStatus.Disable() // This is just an indicator, not clickable
	
	// Toggle jiggling
	mToggle := systray.AddMenuItem("Start Jiggling", "Toggle mouse jiggling")
	
	systray.AddSeparator()

	// Radius submenu
	mRadius := systray.AddMenuItem("Jiggling Radius", "Set the jiggling radius")
	mRadius2px := mRadius.AddSubMenuItem("2 pixels", "Set radius to 2 pixels")
	mRadius5px := mRadius.AddSubMenuItem("5 pixels", "Set radius to 5 pixels")
	mRadius10px := mRadius.AddSubMenuItem("10 pixels", "Set radius to 10 pixels")
	mRadius20px := mRadius.AddSubMenuItem("20 pixels", "Set radius to 20 pixels")

	// Add check mark to current selection
	mRadius5px.Check()

	// Interval submenu
	mInterval := systray.AddMenuItem("Jiggling Interval", "Set the jiggling interval")
	mInterval1s := mInterval.AddSubMenuItem("1 second", "Set interval to 1 second")
	mInterval3s := mInterval.AddSubMenuItem("3 seconds", "Set interval to 3 seconds")
	mInterval5s := mInterval.AddSubMenuItem("5 seconds", "Set interval to 5 seconds")
	mInterval10s := mInterval.AddSubMenuItem("10 seconds", "Set interval to 10 seconds")

	// Add check mark to current selection
	mInterval3s.Check()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Create a channel for jiggle ticker
	jigglerChan := make(chan bool)
	var jigglerTicker *time.Ticker

	// Handle signals for clean termination
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Handle menu events
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				config.isJiggling = !config.isJiggling
				if config.isJiggling {
					// Store original mouse position
					config.originalX, config.originalY = robotgo.GetMousePos()
					
					// Start jiggling
					mToggle.SetTitle("Stop Jiggling")
					mStatus.SetTitle("Status: Active")
					jigglerTicker = time.NewTicker(config.jigglingInterval)
					
					// Start the jiggler in a separate goroutine
					go func() {
						for {
							select {
							case <-jigglerTicker.C:
								jiggleMouse()
							case <-jigglerChan:
								return
							}
						}
					}()
				} else {
					// Stop jiggling
					mToggle.SetTitle("Start Jiggling")
					mStatus.SetTitle("Status: Inactive")
					jigglerTicker.Stop()
					jigglerChan <- true
					
					// Return to original position
					robotgo.MoveMouse(config.originalX, config.originalY)
				}

			// Radius options
			case <-mRadius2px.ClickedCh:
				config.jigglingRadius = 2
				uncheckAll(mRadius2px, mRadius5px, mRadius10px, mRadius20px)
				mRadius2px.Check()
			case <-mRadius5px.ClickedCh:
				config.jigglingRadius = 5
				uncheckAll(mRadius2px, mRadius5px, mRadius10px, mRadius20px)
				mRadius5px.Check()
			case <-mRadius10px.ClickedCh:
				config.jigglingRadius = 10
				uncheckAll(mRadius2px, mRadius5px, mRadius10px, mRadius20px)
				mRadius10px.Check()
			case <-mRadius20px.ClickedCh:
				config.jigglingRadius = 20
				uncheckAll(mRadius2px, mRadius5px, mRadius10px, mRadius20px)
				mRadius20px.Check()

			// Interval options
			case <-mInterval1s.ClickedCh:
				setInterval(1*time.Second, jigglerTicker, jigglerChan)
				uncheckAll(mInterval1s, mInterval3s, mInterval5s, mInterval10s)
				mInterval1s.Check()
			case <-mInterval3s.ClickedCh:
				setInterval(3*time.Second, jigglerTicker, jigglerChan)
				uncheckAll(mInterval1s, mInterval3s, mInterval5s, mInterval10s)
				mInterval3s.Check()
			case <-mInterval5s.ClickedCh:
				setInterval(5*time.Second, jigglerTicker, jigglerChan)
				uncheckAll(mInterval1s, mInterval3s, mInterval5s, mInterval10s)
				mInterval5s.Check()
			case <-mInterval10s.ClickedCh:
				setInterval(10*time.Second, jigglerTicker, jigglerChan)
				uncheckAll(mInterval1s, mInterval3s, mInterval5s, mInterval10s)
				mInterval10s.Check()

			// Quit
			case <-mQuit.ClickedCh:
				if config.isJiggling {
					jigglerTicker.Stop()
					jigglerChan <- true
				}
				systray.Quit()
				return

			// Handle system signals
			case <-signalChan:
				if config.isJiggling {
					jigglerTicker.Stop()
					jigglerChan <- true
				}
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	fmt.Println("Mouse Jiggler exiting...")
}

func jiggleMouse() {
	// Store current position
	x, y := robotgo.GetMousePos()

	// Generate random offsets within radius
	offsetX := rand.Intn(config.jigglingRadius*2+1) - config.jigglingRadius
	offsetY := rand.Intn(config.jigglingRadius*2+1) - config.jigglingRadius

	// Move mouse to new position
	robotgo.MoveMouse(x+offsetX, y+offsetY)

	// Short delay for visibility
	time.Sleep(100 * time.Millisecond)

	// Move back to original position
	robotgo.MoveMouse(x, y)
}

func setInterval(newInterval time.Duration, currentTicker *time.Ticker, stopChan chan bool) {
	config.jigglingInterval = newInterval
	
	// Update ticker if jiggling is active
	if config.isJiggling && currentTicker != nil {
		currentTicker.Stop()
		stopChan <- true
		
		// Create new ticker with updated interval
		newTicker := time.NewTicker(newInterval)
		
		// Start the jiggler again
		go func() {
			for {
				select {
				case <-newTicker.C:
					jiggleMouse()
				case <-stopChan:
					return
				}
			}
		}()
	}
}

func uncheckAll(items ...*systray.MenuItem) {
	for _, item := range items {
		item.Uncheck()
	}
}

// getIcon returns the icon for the menu bar
func getIcon() []byte {
	log.Println("Using tray icon from assets/icon/trayicon.go")
	return icon.Data
}