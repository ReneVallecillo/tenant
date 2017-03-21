package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"gopkg.in/urfave/cli.v2"
)

func main() {
	var language string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lang",
				Value:       "english",
				Usage:       "language for the greeting",
				Destination: &language,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "up",
				Usage: "Start docker database",
				Action: func(c *cli.Context) error {
					// messages := make(chan int)
					var wg sync.WaitGroup

					cmdName := "docker-compose"
					cmdArgs := []string{"-f", "./docker/docker-compose.yml", "up", "-d"}
					fmt.Println("Starting Docker: ")
					cmd := exec.Command(cmdName, cmdArgs...)
					cmdReader, err := cmd.StderrPipe()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
						return err
					}
					if err := cmd.Start(); err != nil {
						fmt.Fprintln(os.Stderr, "Error executing command", err)
						return err
					}

					scanner := bufio.NewScanner(cmdReader)
					wg.Add(1)
					go func() {
						defer wg.Done()
						for scanner.Scan() {
							fmt.Println(scanner.Text())
						}
					}()

					// err = cmd.Start()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
						return err
					}

					err = cmd.Wait()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error waiting for Cmd", cmd.Stderr)
						return err
					}
					wg.Wait()
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "Stoping database",
				Action: func(c *cli.Context) error {
					fmt.Println("Stopping Docker")
					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
