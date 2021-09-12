//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type WaypointSetCurrentReq struct {
	msg.Package `ros:"mavros_msgs"`
	WpSeq       uint16
}

type WaypointSetCurrentRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
}

type WaypointSetCurrent struct {
	msg.Package `ros:"mavros_msgs"`
	WaypointSetCurrentReq
	WaypointSetCurrentRes
}
