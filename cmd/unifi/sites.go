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

func sitesCommand() *cli.Command {
	return &cli.Command{
		Name:    "sites",
		Aliases: []string{"s"},
		Usage:   "Manage UniFi sites",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all sites",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Maximum number of sites to return (0-200)",
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

					params := &unifi.ListSitesParams{
						Limit:  c.Int("limit"),
						Offset: c.Int("offset"),
					}

					ctx := context.Background()
					resp, err := client.ListSites(ctx, params)
					if err != nil {
						return fmt.Errorf("failed to list sites: %w", err)
					}

					if c.Bool("json") {
						return json.NewEncoder(os.Stdout).Encode(resp)
					}

					// Table output
					fmt.Printf("%-36s %-24s\n", "ID", "NAME")
					fmt.Println(strings.Repeat("-", 62))
					for _, site := range resp.Data {
						fmt.Printf("%-36s %-24s\n",
							site.ID,
							truncateString(site.Name, 23),
						)
					}

					fmt.Printf("\nShowing %d of %d sites (offset: %d)\n",
						resp.Count, resp.TotalCount, resp.Offset)
					return nil
				},
			},
		},
	}
}
