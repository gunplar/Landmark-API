package main

import (
	"context"
	"encoding/hex"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"net"
	"os"
	"path/filepath"
)

func ChangeRRSet(client *route53.Client,
	operation types.ChangeAction,
	subDomain string,
	rrContent string) {

	path, err := os.Getwd()
	check(err)
	configDir := filepath.Join(path, "resources", "appConfig.json")
	appConfig := LoadConfiguration(configDir)

	if operation == types.ChangeActionDelete {
		res, err := net.LookupTXT(subDomain + "." + appConfig.ZoneName)
		check(err)
		rrContent = res[0]
	}

	input := route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: operation,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(subDomain + "." + appConfig.ZoneName),
						Type: types.RRTypeTxt,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String("\"" + rrContent + "\""),
							},
						},
						TTL: aws.Int64(600),
					},
				},
			},
			Comment: aws.String("Update using Go SDK."),
		},
		HostedZoneId: aws.String(appConfig.ZoneId),
	}

	_, err = client.ChangeResourceRecordSets(context.Background(), &input)
	check(err)
}

func PublishUserData(client *route53.Client,
	operation types.ChangeAction,
	aesKey string,
	subDomain string,
	rrContent string) {

	var nonce []byte
	rrContent, nonce = encrypt(rrContent, aesKey)
	nonceString := hex.EncodeToString(nonce)

	ChangeRRSet(client, operation, subDomain, rrContent)
	ChangeRRSet(client, operation, "nonce."+subDomain, nonceString)
}
