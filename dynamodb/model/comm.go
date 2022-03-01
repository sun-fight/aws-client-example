package model

import (
	"aws-client-example/dynamodb/define/derr"
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

	GsiPk     = "Pk"
	GsiSk     = "Sk"
	Gsi1Pk    = "Gsi1Pk"
	Gsi1Sk    = "Gsi1Sk"
	Gsi2Pk    = "Gsi2Pk"
	Gsi2Sk    = "Gsi2Sk"
	Gsi3Pk    = "Gsi3Pk"
	Gsi3Sk    = "Gsi3Sk"
	GsiSortPk = "GsiSortPk"

	GsiIdx1 = "GSI1"
	GsiIdx2 = "GSI2"
	GsiIdx3 = "GSI3"
)

func GetPkExp() (exp expression.Expression) {
	cond := expression.Name(Pk).AttributeNotExists()
	exp, _ = expression.NewBuilder().WithCondition(cond).Build()
	return
}

func GetTableName() *string {
	return aws.String(TableName)
}

func GetPk(pkKey string, userID int64) string {
	return pkKey + cast.ToString(userID)
}

func GetPkMap(pk string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: pk}}
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

func AddCondition(v *pb.ExpCondition, condition expression.ConditionBuilder) (conditionRes expression.ConditionBuilder, err error) {
	var where expression.ConditionBuilder
	switch v.ConditionMode {
	case pb.EnumExpConditionMode_ConditionModeEqual:
		where = expression.Name(v.Name).Equal(expression.Value(ValToInterface(v.ValType, v.Value)))
	default:
		err = derr.ErrConditionMode
		return
	}
	conditionRes = condition
	switch v.LogicalMode {
	case pb.EnumExpLogicalMode_LogicalModeAnd:
		conditionRes = conditionRes.And(where)
	case pb.EnumExpLogicalMode_LogicalModeOr:
	case pb.EnumExpLogicalMode_LogicalModeNot:
	default:
		err = derr.ErrLogicalMode
		return
	}
	return
}

func ValToInterface(valType pb.EnumExpValType, val interface{}) interface{} {
	switch valType {
	case pb.EnumExpValType_ValTypeI64:
		return cast.ToInt64(val)
	default:
		return val
	}
}

func ValToBuilder(valType pb.EnumExpValType, val interface{}) expression.ValueBuilder {
	switch valType {
	case pb.EnumExpValType_ValTypeI64:
		val = cast.ToInt64(val)
	default:
	}
	return expression.Value(val)
}

func AddVersion(items []*pb.ExpUpdateItem) []*pb.ExpUpdateItem {
	return append(items, &pb.ExpUpdateItem{
		OperationMode: pb.EnumExpUpdateOperationMode_OperationModeSet,
		ExpUpdateSets: []*pb.ExpUpdateSet{
			{
				Name:       "Version",
				SetValMode: pb.EnumExpUpdateSetValMode_SetValModePlus,
				Value:      "1",
				ValType:    pb.EnumExpValType_ValTypeI64,
			},
		},
	})
}
