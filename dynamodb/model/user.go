package model

import (
	"aws-client-example/dynamodb/define/derr"
	"aws-client-example/dynamodb/pb"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/sun-fight/aws-client/mdynamodb"
)

const (
	PkUser          = "user#"
	_gsiOneUsername = "user#username#"
)

type User struct {
	*pb.TableUser
}

func NewUserDao() *User {
	return &User{}
}

func GetUserPk(id int64) string {
	return GetPk(PkUser, id)
}

func GetUserNameKey(username string) string {
	return _gsiOneUsername + username
}

func CreateUserInfo(userInfo *pb.TableUser) (err error) {
	userInfo.CreatedAt = int32(time.Now().Unix())
	userInfo.LastLoginAt = userInfo.CreatedAt
	userInfo.Version = 1

	cond := expression.Name(Pk).AttributeNotExists()
	exp, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return
	}
	itemMap, err := attributevalue.MarshalMap(&userInfo)
	if err != nil {
		return
	}
	dao := mdynamodb.NewItemDao(TableName)
	_, err = dao.PutItem(mdynamodb.ReqPutItem{
		ItemMap:                   itemMap,
		ConditionExpression:       exp.Condition(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
	})
	return
}

func GetUserInfo(username string) (res *pb.UserInfo, err error) {
	var keyB expression.KeyConditionBuilder
	keyB = keyB.And(expression.Key(Gsi1Pk).Equal(expression.Value(_gsiOneUsername + username)))
	exp, err := expression.NewBuilder().WithKeyCondition(keyB).Build()
	if err != nil {
		return
	}
	dao := mdynamodb.NewItemDao(TableName)
	out, err := dao.Query(mdynamodb.ReqQueryInput{
		IndexName:                 aws.String(GsiIdx1),
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

func (item *User) UpdateUserInfo(user *pb.TableUser, updateCond *pb.UpdateCondition) (err error) {
	item.TableUser = user

	updateCond.ExpUpdateItems = AddVersion(updateCond.ExpUpdateItems)
	updateBuilder, err := item.ToUpdateBuilder(updateCond.ExpUpdateItems)
	if err != nil {
		return
	}
	condition := expression.AttributeExists(expression.Name(Pk)).
		And(expression.Equal(expression.Name("Version"), expression.Value(item.Version)))
	for _, v := range updateCond.ExpConditions {
		condition, err = AddCondition(v, condition)
		if err != nil {
			return
		}
	}
	exp, err := expression.NewBuilder().
		WithUpdate(updateBuilder).
		WithCondition(condition).
		Build()
	if err != nil {
		return
	}

	dao := mdynamodb.NewItemDao(TableName)
	_, err = dao.UpdateItem(mdynamodb.ReqUpdateItem{
		Key:                       GetPkSkMap(user.Pk, user.Sk),
		UpdateExpression:          exp.Update(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		ConditionExpression:       exp.Condition(),
	})
	return
}

func (item *User) Delete(pk string) (err error) {
	dao := mdynamodb.NewItemDao(TableName)
	_, err = dao.DeleteItem(mdynamodb.ReqDeleteItem{
		Key: GetPkMap(pk),
	})
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
				var value expression.ValueBuilder
				var err2 error
				if set.Value != "" {
					value = ValToBuilder(set.ValType, set.Value)
				} else {
					value, err2 = item.NameToVal(set.Name)
					if err2 != nil {
						err = err2
						return
					}
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
		val = item.TableUser.DeletedAt
	case "LastLoginAt":
		val = time.Now().Unix()
	default:
		err = derr.NewErrNamtToVal(name)
		return
	}
	return expression.Value(val), nil
}
