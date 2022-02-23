package main

import (
	"fmt"
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

	username := "a1234"
	// register(username)
	//登录
	user, err := model.LoginByUsername(username)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
}

func register(username string) {
	var err error
	// 注册
	if !ureg.Username(username) {
		log.Fatal("username")
	}
	oauth := &pb.TableOauth{
		Pk:     model.GetOauthPk(username),
		Sk:     model.GetOauthPk(username),
		UserID: uid.Gen64Def(),
	}
	err = model.OauthRegister(oauth)
	if err != nil {
		log.Fatal(err)
	}
}
