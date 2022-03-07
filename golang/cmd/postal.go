/*
Copyright © 2022 PHUC MAI <phuc.mai@here.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// postalCmd represents the postal command
var postalCmd = &cobra.Command{
	Use:   "postal",
	Short: "Including the commands performed on the postal service side.",
	Long:  `Including the commands performed on the postal service side.`,
}

func init() {
	rootCmd.AddCommand(postalCmd)
}
