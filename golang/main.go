package main

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//StoreNewPassword()
	/*awsClient, hash := Login()
	if awsClient == nil {
		return
	}
	ChangeRRSet(awsClient, types.ChangeActionUpsert, hash, "phucmai", "test")*/
	decrypt()
}
