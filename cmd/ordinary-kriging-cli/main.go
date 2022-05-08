package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// these are set in build step
	version = "unversioned"
	// lint:ignore U1000 embedded by goreleaser
	commit = "?"
	// lint:ignore U1000 embedded by goreleaser
	date = "?"
)

func main() {
	execute()

}

var cmd = &cobra.Command{
	Use:   "go-kriging",
	Short: "geospatial prediction and mapping via ordinary kriging",
	Long: `Golang library for geospatial prediction and mapping via ordinary kriging.
			Complete documentation is available at https://github.com/lvisei/go-kriging`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Go Kriging Version: v%s \n", version)
		// TODO:
	},
}

func execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
