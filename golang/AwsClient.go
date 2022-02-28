package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func ChangeUserRRSet(client *route53.Client,
	operation types.ChangeAction,
	aesKey string,
	subDomain string,
	rrContent string) {

	path, err := os.Getwd()
	check(err)
	configDir := filepath.Join(path, "resources", "appConfig.json")
	appConfig := LoadConfiguration(configDir)

	var nonce []byte
	rrContent, nonce = encrypt(rrContent, aesKey)
	nonceString := hex.EncodeToString(nonce)

	if operation == types.ChangeActionDelete {
		res, err := net.LookupTXT(subDomain + "." + appConfig.ZoneName)
		check(err)
		rrContent = res[0]
		res, err = net.LookupTXT("nonce." + subDomain + "." + appConfig.ZoneName)
		check(err)
		nonceString = res[0]
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
				{
					Action: operation,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String("nonce." + subDomain + "." + appConfig.ZoneName),
						Type: types.RRTypeTxt,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String("\"" + nonceString + "\""),
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

func PublishNewKeyPostalService() {

	//Generate the RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	check(err)

	//Write the private key to PEM file
	var privkeyBytes []byte
	privkeyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey)
	check(err)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	path, err := os.Getwd()
	check(err)
	keyDir := filepath.Join(path, "resources", "privatekey.pem")
	err = os.WriteFile(keyDir, privkeyPem, 0644)

	//Generate PEM format of public key
	publicKey := privateKey.PublicKey
	var pubkeyBytes []byte
	pubkeyBytes, err = x509.MarshalPKIXPublicKey(&publicKey)
	check(err)
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	//Create new AWS client and publish the key on a DNS RR
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("aws-global"),
	)
	check(err)
	awsClient := route53.NewFromConfig(cfg)

	configDir := filepath.Join(path, "resources", "config.json")
	appConfig := LoadConfiguration(configDir)

	input := route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String("real.dhl." + appConfig.ZoneName),
						Type: types.RRTypeTxt,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String("\"" + strings.Replace(string(pubkeyPem), "\n", "\"\"", -1) + "\""),
								//TODO the string has "" inserted
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

	_, err = awsClient.ChangeResourceRecordSets(context.Background(), &input)
	check(err)

	/*encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publicKey,
		[]byte("super secret message"),
		nil)
	check(err)

	fmt.Println("encrypted bytes: ", hex.EncodeToString(encryptedBytes))
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	check(err)

	fmt.Println("decrypted message: ", string(decryptedBytes))
	*/
}
