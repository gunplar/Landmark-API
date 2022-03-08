/*
Copyright © 2022 PHUC MAI <phuc.mai@here.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"landmark/internal"
)

var userKeyDomain string
var postalDomain string

// publishEncryptedKeyCmd represents the publishEncryptedKey command
var publishEncryptedKeyCmd = &cobra.Command{
	Use:   "publish-shared",
	Short: "Publish the AES key encypted with the postal company's public key",
	Long: `Query for the public key on the postal service's DNS record, 
then use that key to encrypt the user's own AES key, 
and publish on a new DNS entry.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.PublishEncryptedAESkey(userKeyDomain, postalDomain)
	},
}

func init() {
	userCmd.AddCommand(publishEncryptedKeyCmd)
	publishEncryptedKeyCmd.Flags().StringVar(&userKeyDomain, "domain", "",
		"The subdomain to publish the encrypted AES key")
	err := publishEncryptedKeyCmd.MarkFlagRequired("domain")
	check(err)
	publishEncryptedKeyCmd.Flags().StringVar(&postalDomain, "postal-domain", "",
		"The DNS domain where the postal service publish their public key")
	err = publishEncryptedKeyCmd.MarkFlagRequired("postal-domain")
	check(err)
}
