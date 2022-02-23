package main

import (
	"log"

	"aws-client-example/dynamodb/initial"
	"aws-client-example/dynamodb/model"
	"aws-client-example/dynamodb/pb"
	"aws-client-example/dynamodb/utils/uid"
	"aws-client-example/dynamodb/utils/ureg"
)

func main() {
	uid.Init()

	var err error
	err = initial.InitAwsService()
	if err != nil {
		log.Fatal(err)
	}

	// 注册
	username := "a1234"
	if !ureg.Username(username) {
		log.Fatal("username")
	}
	oauth := &pb.TableOauth{
		Pk:     model.GetOauthPk(username),
		Sk:     model.GetOauthPk(username),
		UserID: uid.Gen64Def(),
	}
	err = model.CreateOauth(oauth)
	if err != nil {
		log.Fatal(err)
	}
	//
}
