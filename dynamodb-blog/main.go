package main

import (
	"log"

	"aws-client-example/dynamodb/initial"
	"aws-client-example/dynamodb/utils/uid"
)

func main() {
	uid.Init()

	var err error
	err = initial.InitAwsService()
	if err != nil {
		log.Fatal(err)
	}

	// 用户
	username := "a12341"
	// register(username)
	login(username)
	select {}

}
