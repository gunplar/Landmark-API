package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"net"
	"os"
	"path/filepath"
)

func ChangeRRSet(client *route53.Client, operation types.ChangeAction, subDomain string, rrContent string) {

	path, err := os.Getwd()
	check(err)
	configDir := filepath.Join(path, "resources", "config.json")
	config := LoadConfiguration(configDir)

	if operation == types.ChangeActionDelete {
		res, err := net.LookupTXT(subDomain + "." + config.ZoneName)
		check(err)
		rrContent = res[0]
	}

	input := route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: operation,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(subDomain + "." + config.ZoneName),
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
		HostedZoneId: aws.String(config.ZoneId),
	}

	_, err = client.ChangeResourceRecordSets(context.Background(), &input)
	check(err)
}
