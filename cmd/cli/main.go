package main

import (
	client "goredis/cli"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "goredis client",
		Usage: "simple cli to access goredis server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   6379,
				Usage:   "port of redis",
			},
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Value:   "localhost",
				Usage:   "address of redis",
			},
		},
		Action: func(ctx *cli.Context) error {
			client := client.NewClient(
				client.WithAddress(ctx.String("address")),
				client.WithPort(ctx.Int("port")),
			)

			client.Start()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
