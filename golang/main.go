package main

import "fmt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//StoreNewPassword()
	awsClient, hash := Login()
	if awsClient == nil {
		return
	}
	fmt.Println(hash)
	//ModifyUserData(awsClient, types.ChangeActionUpsert, hash, "phucmai", "test")
	//decrypt()
	//PublishNewKeyPostalService(awsClient, "real.dhl")
	//res, err := net.LookupTXT("real.dhl.cmtrd.aws.in.here.com")
	PublishEncryptedAESkey(awsClient, "phucmai", "real.dhl", "lololo")
}
