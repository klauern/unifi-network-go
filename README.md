# unifi-network-go

A Go client library for interacting with the UniFi Network Controller API. This library provides a simple and intuitive way to manage UniFi network devices, clients, sites, and hotspot vouchers.

## Features

- Site Management
- Device Management
- Network Client Management
- Hotspot Voucher Management
- Pagination Support
- Error Handling
- Customizable HTTP Client

## Installation

```bash
go get github.com/klauern/unifi-network-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    unifi "github.com/klauern/unifi-network-go"
)

func main() {
    // Create a new client
    client, err := unifi.NewClient("https://192.168.1.1:8443")
    if err != nil {
        log.Fatal(err)
    }

    // Get application info
    info, err := client.GetApplicationInfo(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("UniFi Network Version: %s\n", info.ApplicationVersion)
}
```

## Usage Examples

### Site Management

```go
// List all sites
sites, err := client.ListSites(context.Background(), &unifi.ListSitesParams{
    Limit: 100,
})

// Get a specific site
site, err := client.GetSite(context.Background(), "site-id")
```

### Device Management

```go
// List devices in a site
devices, err := client.ListDevices(context.Background(), "site-id", &unifi.ListDevicesParams{
    Limit: 100,
})

// Get a specific device
device, err := client.GetDevice(context.Background(), "site-id", "device-id")

// Get device statistics
stats, err := client.GetDeviceStatistics(context.Background(), "site-id", "device-id")

// Execute device action (restart, adopt, forget)
err := client.ExecuteDeviceAction(context.Background(), "site-id", "device-id", &unifi.DeviceAction{
    Action: "restart",
})
```

### Network Client Management

```go
// List network clients
clients, err := client.ListNetworkClients(context.Background(), "site-id", &unifi.ListNetworkClientsParams{
    Type: "all",
    WithinHours: 24,
})

// Get a specific client
client, err := client.GetNetworkClient(context.Background(), "site-id", "client-id")

// Block/Unblock a client
err := client.BlockNetworkClient(context.Background(), "site-id", "client-id")
err := client.UnblockNetworkClient(context.Background(), "site-id", "client-id")
```

### Hotspot Voucher Management

```go
// Generate vouchers
vouchers, err := client.GenerateHotspotVouchers(context.Background(), "site-id", &unifi.GenerateHotspotVouchersRequest{
    Count: 5,
    Name: "1-Day Pass",
    TimeLimitMinutes: 1440, // 24 hours
})

// List vouchers
voucherList, err := client.ListHotspotVouchers(context.Background(), "site-id", &unifi.ListHotspotVouchersParams{
    Limit: 100,
})

// Get voucher details
voucher, err := client.GetVoucherDetails(context.Background(), "site-id", "voucher-id")

// Delete a voucher
err := client.DeleteHotspotVoucher(context.Background(), "site-id", "voucher-id")
```

## Customization

### Custom HTTP Client

You can provide your own HTTP client with custom settings:

```go
httpClient := &http.Client{
    Timeout: time.Second * 30,
    // Add other customizations...
}

client, err := unifi.NewClient(
    "https://192.168.1.1:8443",
    unifi.WithHTTPClient(httpClient),
)
```

## Error Handling

The library provides detailed error information through the `unifi.Error` type:

```go
if err != nil {
    if apiErr, ok := err.(*unifi.Error); ok {
        fmt.Printf("API Error: %s (Status: %d)\n", apiErr.Message, apiErr.Status)
    }
}
```

## Pagination

Most list operations support pagination through the `Offset` and `Limit` parameters:

```go
response, err := client.ListDevices(context.Background(), "site-id", &unifi.ListDevicesParams{
    Offset: 0,
    Limit: 100,
})

fmt.Printf("Total Count: %d\n", response.TotalCount)
fmt.Printf("Current Page Count: %d\n", response.Count)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This library is distributed under the MIT license. See the LICENSE file for more information.
