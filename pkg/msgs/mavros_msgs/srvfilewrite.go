//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type FileWriteReq struct {
	msg.Package `ros:"mavros_msgs"`
	FilePath    string
	Offset      uint64
	Data        []uint8
}

type FileWriteRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
	RErrno      int32
}

type FileWrite struct {
	msg.Package `ros:"mavros_msgs"`
	FileWriteReq
	FileWriteRes
}
