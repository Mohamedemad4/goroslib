//nolint:golint
package geographic_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type GeoPointStamped struct {
	msg.Package `ros:"geographic_msgs"`
	Header      std_msgs.Header
	Position    GeoPoint
}
