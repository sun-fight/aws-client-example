package uid

import (
	"github.com/bwmarrin/snowflake"
)

type SnowflakeNode int32

const (
	NodeDef SnowflakeNode = iota
)

var nodes map[SnowflakeNode]*snowflake.Node

func Init() {
	snowflake.NodeBits = 0  //节点数量
	snowflake.StepBits = 22 //每毫秒可产多少
	if snowflake.NodeBits+snowflake.StepBits > 22 {
		panic("snowflake NodeBits add StepBits must less than 22")
	}
	nodes = make(map[SnowflakeNode]*snowflake.Node)
	addNode(NodeDef)
}

func addNode(nodeNum SnowflakeNode) {
	var err error
	newNode, err := snowflake.NewNode(int64(nodeNum))
	if err != nil {
		panic("new snowflake node " + err.Error())
	}
	nodes[nodeNum] = newNode
}

func Gen64(nodeNum SnowflakeNode) int64 {
	return nodes[nodeNum].Generate().Int64()
}

func Gen(nodeNum SnowflakeNode) snowflake.ID {
	return nodes[nodeNum].Generate()
}

func GenDef() snowflake.ID {
	return nodes[NodeDef].Generate()
}

func Gen64Def() int64 {
	return nodes[NodeDef].Generate().Int64()
}
