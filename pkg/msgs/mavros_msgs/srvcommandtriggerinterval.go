//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type CommandTriggerIntervalReq struct {
	msg.Package     `ros:"mavros_msgs"`
	CycleTime       float32
	IntegrationTime float32
}

type CommandTriggerIntervalRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
	Result      uint8
}

type CommandTriggerInterval struct {
	msg.Package `ros:"mavros_msgs"`
	CommandTriggerIntervalReq
	CommandTriggerIntervalRes
}
