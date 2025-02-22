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

func hotspotVouchersCommand() *cli.Command {
	return &cli.Command{
		Name:    "vouchers",
		Aliases: []string{"v"},
		Usage:   "Manage UniFi hotspot vouchers",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all hotspot vouchers",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Maximum number of vouchers to return",
						Value: 25,
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

					params := &unifi.ListHotspotVouchersParams{
						Limit: c.Int("limit"),
					}

					ctx := context.Background()
					resp, err := client.ListHotspotVouchers(ctx, c.String("site"), params)
					if err != nil {
						return fmt.Errorf("failed to list vouchers: %w", err)
					}

					if c.Bool("json") {
						return json.NewEncoder(os.Stdout).Encode(resp.Data)
					}

					// Table output
					fmt.Printf("%-24s %-12s %-15s %-10s %-8s\n", "NOTE", "CODE", "EXPIRES", "LIMIT", "STATUS")
					fmt.Println(strings.Repeat("-", 80))
					for _, voucher := range resp.Data {
						expires := "Never"
						if voucher.ExpiresAt != "" {
							expires = voucher.ExpiresAt
						}
						status := "Active"
						if voucher.Expired {
							status = "Expired"
						}

						fmt.Printf("%-24s %-12s %-15s %-10d %-8s\n",
							truncateString(voucher.Name, 23),
							voucher.Code,
							expires,
							voucher.TimeLimitMinutes,
							status,
						)
					}

					return nil
				},
			},
			{
				Name:  "create",
				Usage: "Create a new hotspot voucher",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
					&cli.StringFlag{
						Name:     "note",
						Usage:    "Voucher note",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "duration",
						Usage:    "Duration in minutes",
						Required: true,
					},
					&cli.IntFlag{
						Name:  "guest-limit",
						Usage: "Maximum number of guests per voucher",
					},
					&cli.IntFlag{
						Name:  "data-limit",
						Usage: "Data usage limit in MB",
					},
					&cli.IntFlag{
						Name:  "down-limit",
						Usage: "Download rate limit in Kbps",
					},
					&cli.IntFlag{
						Name:  "up-limit",
						Usage: "Upload rate limit in Kbps",
					},
					&cli.IntFlag{
						Name:  "count",
						Usage: "Number of vouchers to create",
						Value: 1,
					},
				},
				Action: func(c *cli.Context) error {
					client, err := createClient(c)
					if err != nil {
						return err
					}

					request := &unifi.CreateHotspotVoucherRequest{
						Note:             c.String("note"),
						Duration:         c.Int("duration"),
						TimeLimitMinutes: c.Int("duration"),
						Count:            c.Int("count"),
					}

					if c.IsSet("guest-limit") {
						request.AuthorizeGuestLimit = c.Int("guest-limit")
					}
					if c.IsSet("data-limit") {
						request.DataUsageLimitMB = c.Int("data-limit")
					}
					if c.IsSet("down-limit") {
						request.DownRateLimitKbps = c.Int("down-limit")
					}
					if c.IsSet("up-limit") {
						request.UpRateLimitKbps = c.Int("up-limit")
					}

					ctx := context.Background()
					resp, err := client.CreateHotspotVoucher(ctx, c.String("site"), request)
					if err != nil {
						return fmt.Errorf("failed to create voucher: %w", err)
					}

					return json.NewEncoder(os.Stdout).Encode(resp.Data)
				},
			},
			{
				Name:  "generate",
				Usage: "Generate multiple hotspot vouchers",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
					&cli.StringFlag{
						Name:     "name",
						Usage:    "Voucher note (applied to all generated vouchers)",
						Required: true,
					},
					&cli.IntFlag{
						Name:  "count",
						Usage: "Number of vouchers to generate (1-10000)",
						Value: 1,
					},
					&cli.IntFlag{
						Name:     "time-limit",
						Usage:    "Time limit in minutes (1-1000000)",
						Required: true,
					},
					&cli.IntFlag{
						Name:  "guest-limit",
						Usage: "Maximum number of guests per voucher",
					},
					&cli.IntFlag{
						Name:  "data-limit",
						Usage: "Data usage limit in MB (1-1046576)",
					},
					&cli.IntFlag{
						Name:  "down-limit",
						Usage: "Download rate limit in Kbps (2-100000)",
					},
					&cli.IntFlag{
						Name:  "up-limit",
						Usage: "Upload rate limit in Kbps (2-100000)",
					},
				},
				Action: func(c *cli.Context) error {
					client, err := createClient(c)
					if err != nil {
						return err
					}

					request := &unifi.GenerateHotspotVouchersRequest{
						Count:            c.Int("count"),
						Name:             c.String("name"),
						TimeLimitMinutes: c.Int("time-limit"),
					}

					if c.IsSet("guest-limit") {
						request.AuthorizeGuestLimit = c.Int("guest-limit")
					}
					if c.IsSet("data-limit") {
						request.DataUsageLimitMB = c.Int("data-limit")
					}
					if c.IsSet("down-limit") {
						request.RxRateLimitKbps = c.Int("down-limit")
					}
					if c.IsSet("up-limit") {
						request.TxRateLimitKbps = c.Int("up-limit")
					}

					ctx := context.Background()
					resp, err := client.GenerateHotspotVouchers(ctx, c.String("site"), request)
					if err != nil {
						return fmt.Errorf("failed to generate vouchers: %w", err)
					}

					return json.NewEncoder(os.Stdout).Encode(resp.Data)
				},
			},
			{
				Name:  "get",
				Usage: "Get voucher details",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "Voucher ID",
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
					voucher, err := client.GetVoucherDetails(ctx, c.String("site"), c.String("id"))
					if err != nil {
						return fmt.Errorf("failed to get voucher details: %w", err)
					}

					return json.NewEncoder(os.Stdout).Encode(voucher)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a voucher",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "Voucher ID",
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
					err = client.DeleteHotspotVoucher(ctx, c.String("site"), c.String("id"))
					if err != nil {
						return fmt.Errorf("failed to delete voucher: %w", err)
					}

					fmt.Printf("Successfully deleted voucher %s\n", c.String("id"))
					return nil
				},
			},
		},
	}
}
