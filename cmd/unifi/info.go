package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func appInfoCommand() *cli.Command {
	return &cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "Get UniFi Network application information",
		Flags: []cli.Flag{
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

			ctx := context.Background()
			info, err := client.GetApplicationInfo(ctx)
			if err != nil {
				return fmt.Errorf("failed to get application info: %w", err)
			}

			if c.Bool("json") {
				return json.NewEncoder(os.Stdout).Encode(info)
			}

			fmt.Printf("UniFi Network Version: %s\n", info.ApplicationVersion)
			return nil
		},
	}
}
