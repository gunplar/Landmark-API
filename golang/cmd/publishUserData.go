/*
Copyright © 2022 PHUC MAI <phuc.mai@here.com>

*/
package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/spf13/cobra"
	"landmark/internal"
)

var subDomain string
var rrContent string

// publishUserDataCmd represents the publishUserData command
var publishUserDataCmd = &cobra.Command{
	Use:   "publish-stored",
	Short: "Publish the encrypted data on a DNS domain.",
	Long: `Publish the encrypted data on a DNS domain. 
The nonce used in the AES encryption is also published in a subdomain 
called "nonce." extending the input domain.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ModifyUserData(types.ChangeActionUpsert, subDomain, rrContent)
	},
}

func init() {
	userCmd.AddCommand(publishUserDataCmd)
	publishUserDataCmd.Flags().StringVar(&subDomain, "domain", "",
		"The subdomain to publish the encrypted user data")
	err := publishUserDataCmd.MarkFlagRequired("domain")
	check(err)
	publishUserDataCmd.Flags().StringVar(&rrContent, "content", "",
		"The content to be encrypted and published")
	err = publishUserDataCmd.MarkFlagRequired("content")
	check(err)

}
