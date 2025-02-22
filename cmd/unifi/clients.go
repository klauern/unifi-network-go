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

func clientsCommand() *cli.Command {
	return &cli.Command{
		Name:    "clients",
		Aliases: []string{"c"},
		Usage:   "Manage UniFi network clients",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all network clients",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "site",
						Aliases: []string{"s"},
						Usage:   "Site ID",
						Value:   "default",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Maximum number of clients to return (0-200)",
						Value: 25,
					},
					&cli.IntFlag{
						Name:  "offset",
						Usage: "Starting offset for pagination",
						Value: 0,
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

					params := &unifi.ListNetworkClientsParams{
						Limit:  c.Int("limit"),
						Offset: c.Int("offset"),
					}

					ctx := context.Background()
					resp, err := client.ListNetworkClients(ctx, c.String("site"), params)
					if err != nil {
						return fmt.Errorf("failed to list network clients: %w", err)
					}

					if c.Bool("json") {
						return json.NewEncoder(os.Stdout).Encode(resp)
					}

					// Table output
					fmt.Printf("%-24s %-18s %-15s %-10s\n", "NAME", "MAC", "IP", "TYPE")
					fmt.Println(strings.Repeat("-", 70))
					for _, client := range resp.Data {
						fmt.Printf("%-24s %-18s %-15s %-10s\n",
							truncateString(client.Name, 23),
							client.MACAddress,
							client.IPAddress,
							client.Type,
						)
					}

					fmt.Printf("\nShowing %d of %d clients (offset: %d)\n",
						resp.Count, resp.TotalCount, resp.Offset)
					return nil
				},
			},
		},
	}
}

func truncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length-3] + "..."
}
