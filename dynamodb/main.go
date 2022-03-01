package main

import (
	"log"

	"aws-client-example/dynamodb/initial"
	"aws-client-example/dynamodb/model"
	"aws-client-example/dynamodb/pb"
	"aws-client-example/dynamodb/utils/uid"
	"aws-client-example/dynamodb/utils/utime"

	"github.com/jinzhu/now"
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
	user, err := model.LoginByUsername(username)
	if err != nil {
		log.Fatal(err)
	}
	// login(username)
	//排行榜
	err = model.NewRankDao().Create(model.GetRankPk(pb.EnumRankT_RankTWeek, utime.StrDayByTime(now.BeginningOfWeek())),
		user.UserID, 1)
	if err != nil {
		log.Fatal(err)
	}
	select {}

}
