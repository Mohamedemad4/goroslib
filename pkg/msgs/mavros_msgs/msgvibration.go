//nolint:golint
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type Vibration struct {
	msg.Package `ros:"mavros_msgs"`
	Header      std_msgs.Header
	Vibration   geometry_msgs.Vector3
	Clipping    [3]float32
}
