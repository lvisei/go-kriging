package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	execute()

}

var cmd = &cobra.Command{
	Use:   "go-kriging",
	Short: "geospatial prediction and mapping via ordinary kriging",
	Long: `Golang library for geospatial prediction and mapping via ordinary kriging.
			Complete documentation is available at https://github.com/liuvigongzuoshi/go-kriging`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go Kriging Version: v0.1.0")
		// TODO:
	},
}

func execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
