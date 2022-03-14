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
var deleteEPostalKeyCmd = &cobra.Command{
	Use:   "delete-key",
	Short: "Remove the public key on a DNS domain.",
	Long:  `Remove the public key of postal service on a DNS domain.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Perform the same actions as user's deleteEncryptedKeyCmd
		internal.ChangeRRSet(types.ChangeActionDelete, subDomain, "")
	},
}

func init() {
	postalCmd.AddCommand(deleteEPostalKeyCmd)
	deleteEPostalKeyCmd.Flags().StringVar(&subDomain, "domain", "",
		"The subdomain containing the public key")
	err := deleteEPostalKeyCmd.MarkFlagRequired("domain")
	check(err)
}
