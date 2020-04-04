// Autogenerated with msg-import, do not edit.
package sensor_msgs

import (
	"github.com/aler9/goroslib/msgs"
	"github.com/aler9/goroslib/msgs/std_msgs"
)

type CameraInfo struct {
	msgs.Package    `ros:"sensor_msgs"`
	Header          std_msgs.Header
	Height          msgs.Uint32
	Width           msgs.Uint32
	DistortionModel msgs.String
	D               []msgs.Float64
	K               [9]msgs.Float64
	R               [9]msgs.Float64
	P               [12]msgs.Float64
	BinningX        msgs.Uint32
	BinningY        msgs.Uint32
	Roi             RegionOfInterest
}
