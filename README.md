# Cloudflare DDNS

A tiny Go application that updates Cloudflare DNS records with your current public IPv4 address. This tool is meant for users with dynamic IP addresses who want to keep their domain names pointing to their current location.

It should work right now, but this is still a work in progress

## TODO:
- Add releases
- Make this more easily configurable and simple to use

## Configuration

On first run, the application will create a `config.toml` file with default values in the directory where the program was executed. You need to update this file with your Cloudflare credentials and domain information.

### Required Configuration

Edit the `config.toml` file and update the following values:

```toml
zone-id = "your-zone-id"
dns-record-id = "your-dns-record-id"
cloudflare-email = "your-email@example.com"
cloudflare-api-key = "your-api-key"
domain-name = "your-domain.com"
ttl = 1
type = "A"
```

### Configuration Parameters

| Parameter            | Description                                                     |
| -------------------- | --------------------------------------------------------------- |
| `zone-id`            | Your Cloudflare Zone ID                                         |
| `dns-record-id`      | The DNS record ID to update                                     |
| `cloudflare-email`   | Your Cloudflare account email                                   |
| `cloudflare-api-key` | Your Cloudflare API key                                         |
| `domain-name`        | The domain name to update (e.g., "home.example.com")            |
| `ttl`                | Time to live for the DNS record (in seconds). 1 means automatic |
| `type`               | DNS record type (currently supports "A" records)                |

## Getting Your Cloudflare Credentials

You can get the required credentials on the Cloudflare dashboard. A domain registered with cloudflare is required.

## Usage

### Running the Application

```bash
./cloudflare-ddns
```

The application will:

1. Load configuration from `config.toml` (or create it with default values and exit)
2. Check your current public IP address
3. Compare it with the last known IP (stored in `last-ip.txt`)
4. Update the Cloudflare DNS record if the IP has changed
5. Continue monitoring every 30 minutes

## How It Works

1. **IP Detection**: The application fetches your current public IP from `https://api.ipify.org`. This will be configurable in via the config file in the future
2. **State Management**: It stores the last known IP in `last-ip.txt` to avoid unnecessary API calls
3. **DNS Update**: When the IP changes, it updates the specified Cloudflare DNS record via the Cloudflare API
4. **Scheduling**: The application runs continuously, checking for IP changes every 30 minutes. This should be configurable in the future as well.

## Dependencies

- Uses [ipify.org](https://ipify.org) for public IP detection
- Built with [Viper](https://github.com/spf13/viper) for configuration management
- Integrates with [Cloudflare API](https://developers.cloudflare.com/api/)

