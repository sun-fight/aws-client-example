package model

import (
	"aws-client-example/dynamodb/define/derr"
	"aws-client-example/dynamodb/pb"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sun-fight/aws-client/mdynamodb"
)

const (
	PkOauth = "oauth#"
)

type Oauth struct {
	*pb.TableOauth
}

func NewOauthDao() *Oauth {
	return &Oauth{}
}

func GetOauthPk(username string) string {
	return PkOauth + username
}

func GetOauthPkMap(id int64) map[string]types.AttributeValue {
	return GetPkSkMap(PkOauth, id)
}

func CreateOauth(oauth *pb.TableOauth) (err error) {
	createdAt := int32(time.Now().Unix())
	oauth.CreatedAt = createdAt
	oauth.Version = 1

	cond := expression.Name(Pk).AttributeNotExists()
	exp, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return
	}

	oauthItemMap, err := attributevalue.MarshalMap(&oauth)
	if err != nil {
		return
	}
	userInfo := &pb.UserInfo{
		Pk:          GetUserPk(oauth.UserID),
		Sk:          GetUserPk(oauth.UserID),
		UserID:      oauth.UserID,
		CreatedAt:   createdAt,
		LastLoginAt: createdAt,
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
