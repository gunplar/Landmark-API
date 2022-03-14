/*
Copyright © 2022 PHUC MAI <phuc.mai@here.com>

*/
package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/spf13/cobra"
	"landmark/internal"
)

// publishUserDataCmd represents the publishUserData command
var deleteEncryptedKeyCmd = &cobra.Command{
	Use:   "delete-shared",
	Short: "Remove the encrypted key on a DNS domain.",
	Long:  `Remove the encrypted key shared with a postal service on a DNS domain.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ChangeRRSet(types.ChangeActionDelete, subDomain, "")
	},
}

func init() {
	userCmd.AddCommand(deleteEncryptedKeyCmd)
	deleteEncryptedKeyCmd.Flags().StringVar(&subDomain, "domain", "",
		"The subdomain containing the encrypted key")
	err := deleteEncryptedKeyCmd.MarkFlagRequired("domain")
	check(err)
}
