package internal

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func ChangeRRSet(
	operation types.ChangeAction,
	subDomain string,
	rrContent string) {

	path, err := os.Getwd()
	check(err)
	configDir := filepath.Join(path, "resources", "appConfig.json")
	appConfig := LoadConfiguration(configDir)

	if operation == types.ChangeActionDelete {
		res, err := net.LookupTXT(subDomain + "." + appConfig.ZoneName)
		if len(res) == 0 {
			fmt.Println(subDomain + "." + appConfig.ZoneName + " does not exist.")
			return
		}
		check(err)
		rrContent = res[0]
	}

	rrContent = SplitLongRoute53String(rrContent)

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
						TTL: aws.Int64(60),
					},
				},
			},
			Comment: aws.String("Update using Go SDK."),
		},
		HostedZoneId: aws.String(appConfig.ZoneId),
	}
	client := NewAWSClient()
	_, err = client.ChangeResourceRecordSets(context.Background(), &input)
	check(err)
}

func GetAesKey() string {
	path, err := os.Getwd()
	check(err)
	passDir := filepath.Join(path, "resources", "pass")
	awsKey, err := os.ReadFile(passDir)
	check(err)
	hash := sha256.Sum256(awsKey)
	return hex.EncodeToString(hash[:])
}

func ModifyUserData(
	operation types.ChangeAction,
	subDomain string,
	rrContent string) {

	//Compute AES key from AWS credential
	aesKey := GetAesKey()

	var nonce []byte
	rrContent, nonce = AESencrypt(rrContent, aesKey)
	nonceString := hex.EncodeToString(nonce)

	ChangeRRSet(operation, subDomain, rrContent)
	ChangeRRSet(operation, "nonce."+subDomain, nonceString)
}

func PublishEncryptedAESkey(
	subDomain string,
	postalDomain string) {

	//Load configs
	path, err := os.Getwd()
	check(err)
	configDir := filepath.Join(path, "resources", "appConfig.json")
	appConfig := LoadConfiguration(configDir)
	//Compute AES key from AWS credential
	aesKey := GetAesKey()
	//DNS query the public key from postal service
	var res []string
	res, err = net.LookupTXT(postalDomain + "." + appConfig.ZoneName)
	check(err)
	//Parsing
	res[0] = strings.Replace(res[0], "-----BEGIN RSA PUBLIC KEY-----", "-----BEGIN RSA PUBLIC KEY-----\n", 2)
	res[0] = strings.Replace(res[0], "-----END RSA PUBLIC KEY-----", "\n-----END RSA PUBLIC KEY-----", 2)
	block, _ := pem.Decode([]byte(res[0]))
	if block == nil {
		panic("Cannot parse public key from postal service.")
	}
	key, err := x509.ParsePKCS1PublicKey(block.Bytes) //TODO support other key types?
	check(err)
	//Encrypt the AES key with the postal service public key
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, []byte(aesKey), nil)
	check(err)
	//Publish the cipher on a DNS RR
	ChangeRRSet(types.ChangeActionUpsert, postalDomain+"."+subDomain, SplitLongRoute53String(hex.EncodeToString(encryptedBytes)))
}
