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
var deleteUserDataCmd = &cobra.Command{
	Use:   "delete-stored",
	Short: "Remove the encrypted data on a DNS domain.",
	Long: `Remove the encrypted data on a DNS domain. 
The nonce, used in the AES encryption, in a subdomain 
called "nonce." extending the input domain is also deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ModifyUserData(types.ChangeActionDelete, subDomain, "")
	},
}

func init() {
	userCmd.AddCommand(deleteUserDataCmd)
	deleteUserDataCmd.Flags().StringVar(&subDomain, "domain", "",
		"The subdomain containing the encrypted user data")
	err := deleteUserDataCmd.MarkFlagRequired("domain")
	check(err)
}
