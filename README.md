# TestFlight Monitor

A Go application that monitors TestFlight slots for selected apps and notifies you when a slot becomes available using Pushover and Discord.

## Features

- Monitors multiple TestFlight apps.
- Sends notifications via Pushover and Discord when slots are available.
- Configurable check interval.
- Logs status updates to the console.

## Prerequisites

- [Go](https://go.dev/dl/) 1.19 or higher installed.
- Discord webhook URL.
- Pushover API credentials (User Key and API Token).

## Setup

### Clone the Repository

```bash
git clone <repository-url>
cd main
```

### Configuration

1. Edit the `discordWebhookURL` variable in the code to include your Discord webhook URL:

   ```go
   const discordWebhookURL = "your_discord_webhook_url"
   ```

2. Edit the `pushoverConfig` variable in the code to include your Pushover User Key and API Token:

   ```go
   var pushoverConfig = PushoverConfig{
       UserKey:  "your_pushover_user_key",
       APIToken: "your_pushover_api_token",
   }
   ```

3. Edit the `appsToMonitor` array in the code to add or modify the TestFlight apps you want to monitor.

### Build and Run

1. Build the application:

   ```bash
   go build -o main
   ```

2. Run the application:

   ```bash
   ./main
   ```

## How It Works

1. The script periodically checks the TestFlight links for the apps defined in the `appsToMonitor` array.
2. If slots are available, it sends notifications to:
   - Pushover.
   - Discord (via webhook).
3. Logs the current status in the console.

## Adding New Apps

1. Find the TestFlight link for the app you want to monitor.

2. Add a new entry in the `appsToMonitor` array:

   ```go
   {"AppName", "https://testflight.apple.com/join/your_testflight_code", 0},
   ```

3. Restart the application.

## Example Output

```plaintext
TestFlight Monitor started...
No slot available for WhatsApp (2x)
Slot available for Capcut! Sending notifications...
No slot available for Instagram (3x)
```

## Troubleshooting

### Pushover Notifications Not Working

- Ensure your Pushover User Key and API Token are correctly set in the code.
- Check the API limits for Pushover notifications.

### Discord Notifications Not Working

- Ensure the Discord webhook URL is correctly set in the code.
- Verify that the webhook is active in your Discord server settings.

## Contributing

Feel free to fork this repository, make improvements, and submit a pull request.

