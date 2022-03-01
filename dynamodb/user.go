package main

import (
	"aws-client-example/dynamodb/model"
	"aws-client-example/dynamodb/pb"
	"aws-client-example/dynamodb/utils/ureg"
	"fmt"
	"log"
)

func register(username string) {
	var err error
	// 注册
	if !ureg.Username(username) {
		log.Fatal("username")
	}
	oauth := &pb.TableOauth{
		Pk: model.GetOauthPk(pb.EnumOauthT_OauthTUsername, username),
	}
	err = model.OauthRegister(oauth)
	if err != nil {
		log.Fatal(err)
	}
}

func login(username string) {
	//登录
	user, err := model.LoginByUsername(username)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	// 获取用户登录方式,查询多个
	res, err := model.NewOauthDao().GetUserOauths(user.Pk)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res {
		fmt.Println(v)
	}
	// 修改用户信息
	err = model.NewUserDao().UpdateUserInfo(user.TableUser, &pb.UpdateCondition{
		ExpUpdateItems: []*pb.ExpUpdateItem{
			{OperationMode: pb.EnumExpUpdateOperationMode_OperationModeSet,
				ExpUpdateSets: []*pb.ExpUpdateSet{
					{Name: "LastLoginAt"},
				}},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
