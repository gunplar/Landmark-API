/*
Copyright © 2022 PHUC MAI <phuc.mai@here.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"landmark/internal"
)

var postalKeyDomain string

// publishPostalKeyCmd represents the publishPostalKey command
var publishPostalKeyCmd = &cobra.Command{
	Use:   "publish-key",
	Short: "Publish the postal company's public key",
	Long: `Publish the postal company's public key.
The key pair will be automatically generated, the private key will be saved in a file.
End users will use this public key to encrypt and share the symmetric key on a DNS entry. 
The symmetric key is needed to decrypt the user's landmark location data.'`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.PublishNewKeyPostalService(postalKeyDomain)
	},
}

func init() {
	postalCmd.AddCommand(publishPostalKeyCmd)
	publishPostalKeyCmd.Flags().StringVar(&postalKeyDomain, "domain", "",
		"The subdomain to publish the public key")
	err := publishPostalKeyCmd.MarkFlagRequired("domain")
	check(err)

}
