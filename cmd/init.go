/*
Copyright © 2020 Joel Holmes <holmes89@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/Holmes89/got/core"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new repository",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"directory"},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := core.CreateRepository(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, ".got repo initialized\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
