package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/klauern/unifi-network-go"
	"github.com/urfave/cli/v2"
)

func devicesCommand() *cli.Command {
	return &cli.Command{
		Name:    "devices",
		Aliases: []string{"d"},
		Usage:   "Manage UniFi network devices",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all network devices",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Maximum number of devices to return",
						Value: 25,
					},
					&cli.StringFlag{
						Name:  "type",
						Usage: "Filter by device type",
					},
					&cli.BoolFlag{
						Name:  "json",
						Usage: "Output in JSON format",
						Value: false,
					},
				},
				Action: func(c *cli.Context) error {
					client, err := createClient(c)
					if err != nil {
						return err
					}

					params := &unifi.ListDevicesParams{
						Limit: c.Int("limit"),
						Type:  c.String("type"),
					}

					ctx := context.Background()
					resp, err := client.ListDevices(ctx, c.String("site"), params)
					if err != nil {
						return fmt.Errorf("failed to list devices: %w", err)
					}

					if c.Bool("json") {
						return json.NewEncoder(os.Stdout).Encode(resp.Data)
					}

					// Table output
					fmt.Printf("%-24s %-18s %-15s %-12s %-8s\n", "NAME", "MAC", "IP", "MODEL", "STATUS")
					fmt.Println(strings.Repeat("-", 80))
					for _, device := range resp.Data {
						status := "Offline"
						if device.State == 1 {
							status = "Online"
						}
						if device.Disabled {
							status = "Disabled"
						}

						fmt.Printf("%-24s %-18s %-15s %-12s %-8s\n",
							truncateString(device.Name, 23),
							device.MAC,
							device.IP,
							device.Model,
							status,
						)
					}

					return nil
				},
			},
			{
				Name:  "get",
				Usage: "Get device details",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "Device ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
				},
				Action: func(c *cli.Context) error {
					client, err := createClient(c)
					if err != nil {
						return err
					}

					ctx := context.Background()
					device, err := client.GetDevice(ctx, c.String("site"), c.String("id"))
					if err != nil {
						return fmt.Errorf("failed to get device: %w", err)
					}

					return json.NewEncoder(os.Stdout).Encode(device)
				},
			},
			{
				Name:  "stats",
				Usage: "Get device statistics",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "Device ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
				},
				Action: func(c *cli.Context) error {
					client, err := createClient(c)
					if err != nil {
						return err
					}

					ctx := context.Background()
					stats, err := client.GetDeviceStatistics(ctx, c.String("site"), c.String("id"))
					if err != nil {
						return fmt.Errorf("failed to get device statistics: %w", err)
					}

					return json.NewEncoder(os.Stdout).Encode(stats)
				},
			},
			{
				Name:  "action",
				Usage: "Execute device action (restart, adopt, forget)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "Device ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
					&cli.StringFlag{
						Name:     "action",
						Usage:    "Action to perform (restart, adopt, forget)",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					client, err := createClient(c)
					if err != nil {
						return err
					}

					action := &unifi.DeviceAction{
						Action: c.String("action"),
					}

					ctx := context.Background()
					err = client.ExecuteDeviceAction(ctx, c.String("site"), c.String("id"), action)
					if err != nil {
						return fmt.Errorf("failed to execute device action: %w", err)
					}

					fmt.Printf("Successfully executed %s action on device %s\n", action.Action, c.String("id"))
					return nil
				},
			},
			{
				Name:  "port",
				Usage: "Execute port action (reset, enable, disable)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "Device ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
					&cli.StringFlag{
						Name:     "action",
						Usage:    "Action to perform (reset, enable, disable)",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "port-idx",
						Usage:    "Port index number",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "port-id",
						Usage: "Port identifier",
					},
				},
				Action: func(c *cli.Context) error {
					client, err := createClient(c)
					if err != nil {
						return err
					}

					action := &unifi.DevicePortAction{
						PortIDX: c.Int("port-idx"),
						PortID:  c.String("port-id"),
						Action:  c.String("action"),
					}

					ctx := context.Background()
					err = client.ExecutePortAction(ctx, c.String("site"), c.String("id"), action)
					if err != nil {
						return fmt.Errorf("failed to execute port action: %w", err)
					}

					fmt.Printf("Successfully executed %s action on port %d of device %s\n",
						action.Action, action.PortIDX, c.String("id"))
					return nil
				},
			},
		},
	}
}
