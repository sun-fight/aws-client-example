package main

import (
	"log"

	"github.com/sun-fight/aws-client/mdynamodb/example/initial"
	"github.com/sun-fight/aws-client/mdynamodb/example/model"
	"github.com/sun-fight/aws-client/mdynamodb/pb"
)

func main() {

	var err error
	err = initial.InitAwsService()
	if err != nil {
		log.Fatal(err)
	}

	userInfo := &pb.UserInfo{}
	err = model.NewUserDao().CreateUserInfo(userInfo)
	if err != nil {
		log.Fatal(err)
	}
}
