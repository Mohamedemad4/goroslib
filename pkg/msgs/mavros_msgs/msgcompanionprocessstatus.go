//nolint:golint
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

const (
	CompanionProcessStatus_MAV_STATE_UNINIT                     uint8 = 0
	CompanionProcessStatus_MAV_STATE_BOOT                       uint8 = 1
	CompanionProcessStatus_MAV_STATE_CALIBRATING                uint8 = 2
	CompanionProcessStatus_MAV_STATE_STANDBY                    uint8 = 3
	CompanionProcessStatus_MAV_STATE_ACTIVE                     uint8 = 4
	CompanionProcessStatus_MAV_STATE_CRITICAL                   uint8 = 5
	CompanionProcessStatus_MAV_STATE_EMERGENCY                  uint8 = 6
	CompanionProcessStatus_MAV_STATE_POWEROFF                   uint8 = 7
	CompanionProcessStatus_MAV_STATE_FLIGHT_TERMINATION         uint8 = 8
	CompanionProcessStatus_MAV_COMP_ID_OBSTACLE_AVOIDANCE       uint8 = 196
	CompanionProcessStatus_MAV_COMP_ID_VISUAL_INERTIAL_ODOMETRY uint8 = 197
)

type CompanionProcessStatus struct {
	msg.Package     `ros:"mavros_msgs"`
	msg.Definitions `ros:"uint8 MAV_STATE_UNINIT=0,uint8 MAV_STATE_BOOT=1,uint8 MAV_STATE_CALIBRATING=2,uint8 MAV_STATE_STANDBY=3,uint8 MAV_STATE_ACTIVE=4,uint8 MAV_STATE_CRITICAL=5,uint8 MAV_STATE_EMERGENCY=6,uint8 MAV_STATE_POWEROFF=7,uint8 MAV_STATE_FLIGHT_TERMINATION=8,uint8 MAV_COMP_ID_OBSTACLE_AVOIDANCE=196,uint8 MAV_COMP_ID_VISUAL_INERTIAL_ODOMETRY=197"`
	Header          std_msgs.Header
	State           uint8
	Component       uint8
}
