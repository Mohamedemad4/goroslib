//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

const (
	CommandVtolTransitionReq_STATE_MC uint8 = 3
	CommandVtolTransitionReq_STATE_FW uint8 = 4
)

type CommandVtolTransitionReq struct {
	msg.Package     `ros:"mavros_msgs"`
	msg.Definitions `ros:"uint8 STATE_MC=3,uint8 STATE_FW=4"`
	Header          std_msgs.Header
	State           uint8
}

type CommandVtolTransitionRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
	Result      uint8
}

type CommandVtolTransition struct {
	msg.Package `ros:"mavros_msgs"`
	CommandVtolTransitionReq
	CommandVtolTransitionRes
}
