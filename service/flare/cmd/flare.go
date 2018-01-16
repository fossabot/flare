// Copyright 2018 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/diegobernardes/flare/service/flare"
)

func main() {
	var configPath string

	cmdStart := &cobra.Command{
		Use:   "start",
		Short: "Start Flare service",
		Long: `This command is used to start the Flare service. The application gonna
look for a 'flare.toml' file at the same directory as the binary.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := readConfig(configPath)
			if err != nil {
				fmt.Println(errors.Wrap(err, "could not load configuration file"))
				os.Exit(1)
			}

			var options []func(*flare.Client)
			if config != "" {
				options = append(options, flare.ClientConfig(config))
			}

			client, err := flare.NewClient(options...)
			if err != nil {
				fmt.Println(errors.Wrap(err, "error during client initialization"))
				os.Exit(1)
			}

			if err := client.Start(); err != nil {
				fmt.Println(errors.Wrap(err, "error during client start"))
				os.Exit(1)
			}

			chanExit := make(chan os.Signal, 1)
			signal.Notify(chanExit, os.Interrupt)
			<-chanExit

			if err := client.Stop(); err != nil {
				fmt.Println(errors.Wrap(err, "error during client stop"))
				os.Exit(1)
			}
		},
	}
	cmdStart.PersistentFlags().StringVarP(&configPath, "config", "c", "./flare.toml", "")

	cmdSetup := &cobra.Command{
		Use:   "setup",
		Short: "Setup the required resources",
		Long:  "Based at the configuration, it run the setup on all required resources.",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := readConfig(configPath)
			if err != nil && configPath != "./flare.toml" {
				fmt.Println(errors.Wrap(err, "could not load configuration file"))
				os.Exit(1)
			}

			ctx, ctxCancel := context.WithCancel(context.Background())
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				<-c
				ctxCancel()
			}()

			var options []func(*flare.Client)
			if config != "" {
				options = append(options, flare.ClientConfig(config))
			}

			client, err := flare.NewClient(options...)
			if err != nil {
				fmt.Println(errors.Wrap(err, "error during client initialization"))
				os.Exit(1)
			}

			if err := client.Setup(ctx); err != nil {
				fmt.Println(errors.Wrap(err, "error during client setup"))
				os.Exit(1)
			}
		},
	}
	cmdSetup.PersistentFlags().StringVarP(&configPath, "config", "c", "./flare.toml", "")

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Show the Flare Version",
		Long:  "Show information about the Go, Repository and Flare version.",
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)

			if flare.Version != "" {
				fmt.Fprintln(w, fmt.Sprintf("Version:\t%s", flare.Version))
			}

			if flare.Commit != "" {
				fmt.Fprintln(w, fmt.Sprintf("Commit:\t%s", flare.Commit))
			}

			if flare.BuildTime != "" {
				fmt.Fprintln(w, fmt.Sprintf("Build Time:\t%s", flare.BuildTime))
			}

			fmt.Fprintln(w, fmt.Sprintf("Go Version:\t%s", flare.GoVersion))
			w.Flush()
		},
	}

	var rootCmd = &cobra.Command{Use: "flare"}
	rootCmd.AddCommand(cmdStart, cmdSetup, cmdVersion)
	rootCmd.Execute()
}

func readConfig(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
