package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func CreateAWSClient() {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("aws-global"),
	)
	check(err)
	client := route53.NewFromConfig(cfg)

	input := route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionDelete,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String("phucmai1.cmtrd.aws.in.here.com"),
						Type: types.RRTypeTxt,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String("\"LOC 123\""),
							},
						},
						TTL: aws.Int64(172900),
					},
				},
			},
			Comment: aws.String("Update using Go SDK."),
		},
		HostedZoneId: aws.String("Z058992925XBNXACB8HLT"), //todo
	}

	_, err = client.ChangeResourceRecordSets(context.Background(), &input)
	check(err)
}
