/*
Copyright © 2022 PHUC MAI <phuc.mai@here.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"landmark/internal"
)

var endUserDomain string
var ownPostalDomain string

// retrieveUserDataCmd represents the retrieveUserData command
var retrieveUserDataCmd = &cobra.Command{
	Use:   "show-user-data",
	Short: "Display decrypted user data",
	Long: `Display decrypted user data. This step will query the user DNS entry
for the encrypted symmetric key, then decrypt that key, and subsequently 
use the key to decrypt the user data on their DNS domain.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.RetrieveUserData(ownPostalDomain, endUserDomain)
	},
}

func init() {
	postalCmd.AddCommand(retrieveUserDataCmd)
	retrieveUserDataCmd.Flags().StringVar(&endUserDomain, "user-domain", "",
		"The subdomain where the user data resides")
	err := retrieveUserDataCmd.MarkFlagRequired("user-domain")
	check(err)
	retrieveUserDataCmd.Flags().StringVar(&ownPostalDomain, "postal-domain", "",
		"The DNS domain where the postal service publish their public key")
	err = retrieveUserDataCmd.MarkFlagRequired("postal-domain")
	check(err)
}
