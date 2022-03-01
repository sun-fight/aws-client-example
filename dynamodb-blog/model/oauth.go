package model

import (
	"aws-client-example/dynamodb/define/derr"
	"aws-client-example/dynamodb/pb"
	"aws-client-example/dynamodb/utils/uid"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sun-fight/aws-client/mdynamodb"
)

const (
	PkPreOauth = "oauth#"
)

type Oauth struct {
	*pb.TableOauth
}

func NewOauthDao() *Oauth {
	return &Oauth{}
}

func GetOauthPk(t pb.EnumOauthT, username string) string {
	return PkPreOauth + t.String() + "#" + username
}

func (item *Oauth) GetUserOauths(userPk string) (res []Oauth, err error) {
	cond := expression.Key(Gsi1Pk).Equal(expression.Value(userPk)).
		And(expression.KeyBeginsWith(expression.Key(Gsi1Sk), PkPreOauth))
	exp, err := expression.NewBuilder().WithKeyCondition(cond).Build()
	if err != nil {
		return
	}
	out, err := mdynamodb.NewItemDao(TableName).Query(mdynamodb.ReqQueryInput{
		IndexName:                 aws.String(GsiIdx1),
		KeyConditionExpression:    exp.KeyCondition(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
	})
	if err != nil {
		return
	}
	for _, v := range out.Items {
		var m Oauth
		err = attributevalue.UnmarshalMap(v, &m)
		if err != nil {
			return
		}
		res = append(res, m)
	}
	return
}

func OauthRegister(oauth *pb.TableOauth) (err error) {
	nowTime := int32(time.Now().Unix())
	userID := uid.Gen64Def()
	userPk := GetUserPk(userID)

	oauth.Sk = oauth.Pk
	oauth.Gsi1Pk = userPk
	oauth.Gsi1Sk = oauth.Pk
	oauth.CreatedAt = nowTime
	oauth.Version = 1

	exp := GetPkExp()

	oauthItemMap, err := attributevalue.MarshalMap(&oauth)
	if err != nil {
		return
	}
	userInfo := &pb.TableUser{
		Pk:          userPk,
		Sk:          userPk,
		Gsi1Pk:      userPk,
		Gsi1Sk:      userPk,
		UserID:      userID,
		CreatedAt:   nowTime,
		LastLoginAt: nowTime,
		Version:     1,
	}
	userItemMap, err := attributevalue.MarshalMap(&userInfo)
	if err != nil {
		return
	}
	dao := mdynamodb.NewTransactDao()
	_, err = dao.TransactWriteItems(mdynamodb.ReqTransactWriteItems{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:                 GetTableName(),
					Item:                      userItemMap,
					ConditionExpression:       exp.Condition(),
					ExpressionAttributeNames:  exp.Names(),
					ExpressionAttributeValues: exp.Values(),
				},
			},
			{
				Put: &types.Put{
					TableName:                 GetTableName(),
					Item:                      oauthItemMap,
					ConditionExpression:       exp.Condition(),
					ExpressionAttributeNames:  exp.Names(),
					ExpressionAttributeValues: exp.Values(),
				},
			},
		},
	})
	return
}

func LoginByUsername(username string) (user User, err error) {
	return login(&pb.TableOauth{
		Pk: GetOauthPk(pb.EnumOauthT_OauthTUsername, username),
	})
}

func login(oauth *pb.TableOauth) (user User, err error) {
	dao := mdynamodb.NewItemDao(TableName)
	_, err = dao.GetItem(mdynamodb.ReqGetItem{
		Key:            GetPkSkMap(oauth.Pk, oauth.Pk),
		ConsistentRead: aws.Bool(true),
	}, &oauth)
	if err != nil {
		return
	}

	tableUser := &pb.TableUser{
		Pk: oauth.Gsi1Pk,
		Sk: oauth.Gsi1Pk,
	}
	_, err = dao.GetItem(mdynamodb.ReqGetItem{
		Key:            GetPkSkMap(tableUser.Pk, tableUser.Pk),
		ConsistentRead: aws.Bool(true),
	}, &tableUser)
	if err != nil {
		return
	}
	user.TableUser = tableUser
	return
}

func (item *Oauth) ToUpdateBuilder(updateItems []*pb.ExpUpdateItem) (updateBuilder expression.UpdateBuilder, err error) {
	if len(updateItems) == 0 {
		err = derr.ErrUpdateItemNoSet
		return
	}
	for _, v := range updateItems {
		switch v.OperationMode {
		case pb.EnumExpUpdateOperationMode_OperationModeSet:
			for _, set := range v.ExpUpdateSets {
				value, err2 := item.NameToVal(set.Name)
				if err2 != nil {
					err = err2
					return
				}
				updateBuilder = SetToUpdateBuilder(set, value, updateBuilder)
			}
		default:
			err = derr.ErrUpdateItemOperationMode
			return
		}
	}
	return
}

func (item *Oauth) NameToVal(name string) (vv expression.ValueBuilder, err error) {
	var val interface{}
	switch name {
	case "Version":
		val = item.Version
	default:
		err = derr.NewErrNamtToVal(name)
		return
	}
	return expression.Value(val), nil
}
