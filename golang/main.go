package main

import (
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//StoreNewPassword()
	//Login()
	ChangeRRSet(Login(), types.ChangeActionDelete, "phucmai1", "test")
}
