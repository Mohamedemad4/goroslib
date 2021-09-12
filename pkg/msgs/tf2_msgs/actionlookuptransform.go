//nolint:golint,lll
package tf2_msgs

import (
	"time"

	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

type LookupTransformActionGoal struct {
	msg.Package `ros:"tf2_msgs"`
	TargetFrame string
	SourceFrame string
	SourceTime  time.Time
	Timeout     time.Duration
	TargetTime  time.Time
	FixedFrame  string
	Advanced    bool
}

type LookupTransformActionResult struct {
	msg.Package `ros:"tf2_msgs"`
	Transform   geometry_msgs.TransformStamped
	Error       TF2Error
}

type LookupTransformActionFeedback struct {
	msg.Package `ros:"tf2_msgs"`
}

type LookupTransformAction struct {
	msg.Package `ros:"tf2_msgs"`
	LookupTransformActionGoal
	LookupTransformActionResult
	LookupTransformActionFeedback
}
