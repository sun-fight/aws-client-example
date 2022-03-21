package model

import "errors"

var ErrUpdateItemNoSet = errors.New("set updateItem to update")
var ErrUpdateItemOperationMode = errors.New("operationMode must set")
var ErrConditionMode = errors.New("conditionMode not set")
var ErrLogicalMode = errors.New("logicalMode not set")

type ErrNameToVal struct {
	Name string
}

func NewErrNameToVal(name string) *ErrNameToVal {
	return &ErrNameToVal{
		Name: name,
	}
}

func (e *ErrNameToVal) Error() string {
	return "can't find name:" + e.Name
}

type ErrParamsValid struct {
	Name string
}

func NewErrParamsValid(name string) *ErrParamsValid {
	return &ErrParamsValid{
		Name: name,
	}
}

func (e *ErrParamsValid) Error() string {
	return "params valid:" + e.Name
}
