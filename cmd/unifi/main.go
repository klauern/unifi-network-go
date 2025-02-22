package main

import (
	"fmt"
	"os"

	"github.com/klauern/unifi-network-go"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "unifi",
		Usage: "UniFi Network API CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "UniFi Network Controller URL",
				EnvVars:  []string{"UNIFI_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "api-key",
				Usage:    "UniFi Network API Key",
				EnvVars:  []string{"UNIFI_API_KEY"},
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "insecure",
				Usage:   "Skip TLS certificate verification",
				EnvVars: []string{"UNIFI_INSECURE"},
			},
		},
		Commands: []*cli.Command{
			clientsCommand(),
			devicesCommand(),
			hotspotVouchersCommand(),
			sitesCommand(),
			appInfoCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func createClient(c *cli.Context) (*unifi.Client, error) {
	client, err := unifi.NewClient(
		c.String("url"),
		unifi.WithAPIKey(c.String("api-key")),
		unifi.WithInsecure(c.Bool("insecure")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}
