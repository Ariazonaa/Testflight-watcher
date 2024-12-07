package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// AppConfig represents a TestFlight app to monitor
type AppConfig struct {
	Name        string // App name
	URL         string // TestFlight URL
	NoSlotCount int    // Counter for no-slot checks
}

// PushoverConfig represents the configuration for Pushover notifications
type PushoverConfig struct {
	UserKey  string
	APIToken string
}

// Apps to monitor for available slots
var appsToMonitor = []AppConfig{
	{"WhatsApp", "https://testflight.apple.com/join/krUFQpyJ", 0},
	{"Capcut", "https://testflight.apple.com/join/Gu9kI6ky", 0},
	{"Instagram", "https://testflight.apple.com/join/72eyUWVE", 0},
	{"WhatsApp Business", "https://testflight.apple.com/join/oscYikr0", 0},
	{"Snapchat", "https://testflight.apple.com/join/p7hGbZUR", 0},
}

// Pushover configuration (uses environment variables for security)
var pushoverConfig = PushoverConfig{
	UserKey:  os.Getenv("PUSHOVER_USER_KEY"),
	APIToken: os.Getenv("PUSHOVER_API_TOKEN"),
}

const checkInterval = 5 * time.Second
var discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL") // Discord webhook URL

// checkAppStatus checks if a TestFlight slot is available
func checkAppStatus(appURL string) (bool, error) {
	resp, err := http.Get(appURL)
	if err != nil {
		return false, fmt.Errorf("Error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("Unexpected HTTP status code: %d", resp.StatusCode)
	}

	// Check if the page indicates the beta is full
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("Error reading page content: %v", err)
	}
	body := string(bodyBytes)

	return !strings.Contains(body, "This beta is full"), nil
}

// sendPushoverNotification sends a Pushover notification
func sendPushoverNotification(appName, appURL string) error {
	message := url.Values{}
	message.Set("token", pushoverConfig.APIToken)
	message.Set("user", pushoverConfig.UserKey)
	message.Set("title", fmt.Sprintf("Slot available: %s", appName))
	message.Set("message", fmt.Sprintf("A slot for the app '%s' is available. Here is the link: %s", appName, appURL))
	message.Set("url", appURL)
	message.Set("url_title", "TestFlight Link")

	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", message)
	if err != nil {
		return fmt.Errorf("Error sending Pushover notification: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		return fmt.Errorf("Error sending Pushover notification: %v", response)
	}
	return nil
}

// sendDiscordNotification sends a notification to Discord
func sendDiscordNotification(appName, appURL string) error {
	payload := map[string]string{
		"content": fmt.Sprintf("Slot available for **%s**! Check it out here: %s", appName, appURL),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error creating Discord payload: %v", err)
	}

	req, err := http.NewRequest("POST", discordWebhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("Error creating Discord request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending Discord notification: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Unexpected Discord response: %s", body)
	}
	return nil
}

// clearConsole clears the terminal output
func clearConsole() {
	cmd := "clear"
	if strings.Contains(os.Getenv("OS"), "Windows") {
		cmd = "cls"
	}
	fmt.Print("\033[H\033[2J", cmd)
}

func main() {
	fmt.Println("TestFlight Monitor started...")
	for {
		clearConsole()
		for i, app := range appsToMonitor {
			status, err := checkAppStatus(app.URL)
			if err != nil {
				fmt.Printf("Error checking %s: %v\n", app.Name, err)
				continue
			}

			if status {
				fmt.Printf("Slot available for %s! Sending notifications...\n", app.Name)
				if err := sendPushoverNotification(app.Name, app.URL); err != nil {
					fmt.Printf("Error sending Pushover notification for %s: %v\n", app.Name, err)
				}
				if err := sendDiscordNotification(app.Name, app.URL); err != nil {
					fmt.Printf("Error sending Discord notification for %s: %v\n", app.Name, err)
				}
				// Reset NoSlotCount when slots are available
				appsToMonitor[i].NoSlotCount = 0
			} else {
				appsToMonitor[i].NoSlotCount++
				fmt.Printf("No slot available for %s (%dx)\n", app.Name, appsToMonitor[i].NoSlotCount)
			}
		}
		time.Sleep(checkInterval)
	}
}
