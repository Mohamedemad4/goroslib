//nolint:golint
package vision_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

type ObjectHypothesisWithPose struct {
	msg.Package `ros:"vision_msgs"`
	Id          int64
	Score       float64
	Pose        geometry_msgs.PoseWithCovariance
}
