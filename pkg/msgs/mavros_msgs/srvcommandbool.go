//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type CommandBoolReq struct {
	msg.Package `ros:"mavros_msgs"`
	Value       bool
}

type CommandBoolRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
	Result      uint8
}

type CommandBool struct {
	msg.Package `ros:"mavros_msgs"`
	CommandBoolReq
	CommandBoolRes
}
