//nolint:golint
package tf2_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

type TFMessage struct {
	msg.Package `ros:"tf2_msgs"`
	Transforms  []geometry_msgs.TransformStamped
}
