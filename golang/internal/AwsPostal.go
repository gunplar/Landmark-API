package internal

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func RetrieveUserData(
	postalDomain string,
	userDomain string) {
	//Load private key from file
	path, err := os.Getwd()
	check(err)
	keyDir := filepath.Join(path, "resources", "privatekey.pem")
	var privateKeyByte []byte
	privateKeyByte, err = os.ReadFile(keyDir)
	block, _ := pem.Decode(privateKeyByte)
	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	check(err)
	//Load configs
	configDir := filepath.Join(path, "resources", "appConfig.json")
	appConfig := LoadConfiguration(configDir)
	//DNS query the encrypted AES key from the end user
	var res []string
	res, err = net.LookupTXT(postalDomain + "." + userDomain + "." + appConfig.ZoneName)
	check(err)
	encryptedBytes, err := hex.DecodeString(res[0])
	check(err)
	//Decrypt to get the AES key
	aesKey, err := privkey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	check(err)
	//Query user's info and decrypt
	res, err = net.LookupTXT(userDomain + "." + appConfig.ZoneName)
	check(err)
	nonce, err := net.LookupTXT("nonce." + userDomain + "." + appConfig.ZoneName)
	check(err)
	AESdecrypt(string(aesKey), res[0], nonce[0])
}

func PublishNewKeyPostalService(
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
	pubkeyPemString := strings.Replace(string(pubkeyPem), "\n", "", -1)
	ChangeRRSet(types.ChangeActionUpsert, subDomain, pubkeyPemString)
}
