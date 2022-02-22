package model

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sun-fight/aws-client/mdynamodb"
	"github.com/sun-fight/aws-client/mdynamodb/example/define/derr"
	"github.com/sun-fight/aws-client/mdynamodb/pb"
)

const (
	PkUser          = "user#"
	_gsiOneUsername = "user#username#"
)

type User struct {
	*pb.UserInfo
}

func NewUser(userInfo *pb.UserInfo) *User {
	return &User{
		UserInfo: userInfo,
	}
}
func NewUserDao() *User {
	return &User{}
}

func GetUserPk(id int64) string {
	return GetPk(PkUser, id)
}

func GetUserPkMap(id int64) map[string]types.AttributeValue {
	return GetPkSkMap(PkUser, id)
}

func (item *User) CreateUserInfo(userInfo *pb.UserInfo) (err error) {
	var keyB expression.ConditionBuilder
	keyB = keyB.And(expression.Name(GsiOnePk).AttributeNotExists()).
		And(expression.Name(Pk).AttributeNotExists())
	exp, err := expression.NewBuilder().WithCondition(keyB).Build()
	if err != nil {
		return
	}
	itemMap, err := attributevalue.MarshalMap(&item.UserInfo)
	if err != nil {
		return
	}
	dao := mdynamodb.NewItemDao(TableName)
	dao.PutItem(mdynamodb.ReqPutItem{
		ItemMap:                   itemMap,
		ConditionExpression:       exp.Condition(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
	})
	return
}

func GetUserInfo(username string) (res *pb.UserInfo, err error) {
	var keyB expression.KeyConditionBuilder
	keyB = keyB.And(expression.Key(GsiOnePk).Equal(expression.Value(_gsiOneUsername + username)))
	exp, err := expression.NewBuilder().WithKeyCondition(keyB).Build()
	if err != nil {
		return
	}
	dao := mdynamodb.NewItemDao(TableName)
	out, err := dao.Query(mdynamodb.ReqQueryInput{
		IndexName:                 aws.String(GsiOneName),
		KeyConditionExpression:    exp.KeyCondition(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
	})
	if err != nil {
		return
	}
	err = attributevalue.UnmarshalMap(out.Items[0], &res)
	return
}

func (item *User) ToUpdateBuilder(updateItems []*pb.ExpUpdateItem) (updateBuilder expression.UpdateBuilder, err error) {
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

func (item *User) NameToVal(name string) (vv expression.ValueBuilder, err error) {
	var val interface{}
	switch name {
	case "Version":
		val = item.Version
	case "DeletedAt":
		val = item.UserInfo.DeletedAt
	case "LastLoginAt":
		val = item.UserInfo.LastLoginAt
	default:
		err = derr.NewErrNamtToVal(name)
		return
	}
	return expression.Value(val), nil
}
