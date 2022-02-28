package main

import "github.com/aws/aws-sdk-go-v2/service/route53/types"

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
	PublishUserData(awsClient, types.ChangeActionUpsert, hash, "phucmai", "test")
	decrypt()
	PublishNewKeyPostalService(awsClient, "real.dhl")
}
