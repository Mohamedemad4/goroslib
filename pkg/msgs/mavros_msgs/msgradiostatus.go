//nolint:golint
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type RadioStatus struct {
	msg.Package `ros:"mavros_msgs"`
	Header      std_msgs.Header
	Rssi        uint8
	Remrssi     uint8
	Txbuf       uint8
	Noise       uint8
	Remnoise    uint8
	Rxerrors    uint16
	Fixed       uint16
	RssiDbm     float32
	RemrssiDbm  float32
}
