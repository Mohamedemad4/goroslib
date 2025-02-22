//nolint:golint
package vision_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type Detection3DArray struct {
	msg.Package `ros:"vision_msgs"`
	Header      std_msgs.Header
	Detections  []Detection3D
}
