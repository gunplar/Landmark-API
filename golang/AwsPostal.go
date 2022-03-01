package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"os"
	"path/filepath"
	"strings"
)

func PublishNewKeyPostalService(
	client *route53.Client,
	subDomain string) {

	//Generate the RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	check(err)

	//Write the private key to PEM file
	var privkeyBytes []byte
	privkeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
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
	pubkeyBytes = x509.MarshalPKCS1PublicKey(&publicKey)
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	//Create new AWS client and publish the key on a DNS RR
	ChangeRRSet(client, types.ChangeActionUpsert, subDomain, strings.Replace(string(pubkeyPem), "\n", "\"\"", -1))

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
