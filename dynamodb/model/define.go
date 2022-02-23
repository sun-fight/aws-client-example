package model

import (
	"aws-client-example/dynamodb/pb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cast"
)

const (
	TableName = "Tables"

	Pk = "Pk"
	Sk = "Sk"

	GsiOnePk = "GsiOnePk"
	GsiOneSk = "GsiOneSk"

	GsiOneName      = "gsi-one"
	GsiTwoName      = "gsi-two"
	GsiInvertedName = "gsi-inverted"
)

func GetTableName() *string {
	return aws.String(TableName)
}

func GetPk(pkKey string, userID int64) string {
	return pkKey + cast.ToString(userID)
}

func GetPkMap(pkKey string, userID int64) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: pkKey + cast.ToString(userID)}}
}

func GetPkSkMap(pk, sk string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: pk},
		"Sk": &types.AttributeValueMemberS{Value: sk}}
}

func SetToUpdateBuilder(set *pb.ExpUpdateSet, value expression.ValueBuilder, updateBuilder expression.UpdateBuilder) expression.UpdateBuilder {
	name := expression.Name(set.Name)
	var operand expression.OperandBuilder
	switch set.SetValMode {
	case pb.EnumExpUpdateSetValMode_SetValModePlus:
		operand = name.Plus(value)
	case pb.EnumExpUpdateSetValMode_SetValModeMinus:
		operand = name.Minus(value)
	case pb.EnumExpUpdateSetValMode_SetValModeIfNotExists:
		operand = name.IfNotExists(value)
	default:
		operand = value
	}
	return updateBuilder.Set(name, operand)
}
