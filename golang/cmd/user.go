/*
Copyright © 2022 PHUC MAI <phuc.mai@here.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Including the commands performed on the user side.",
	Long:  `Including the commands performed on the user side.`,
}

func init() {
	rootCmd.AddCommand(userCmd)
}
