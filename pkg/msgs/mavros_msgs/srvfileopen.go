//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

const (
	FileOpenReq_MODE_READ   uint8 = 0
	FileOpenReq_MODE_WRITE  uint8 = 1
	FileOpenReq_MODE_CREATE uint8 = 2
)

type FileOpenReq struct {
	msg.Package     `ros:"mavros_msgs"`
	msg.Definitions `ros:"uint8 MODE_READ=0,uint8 MODE_WRITE=1,uint8 MODE_CREATE=2"`
	FilePath        string
	Mode            uint8
}

type FileOpenRes struct {
	msg.Package `ros:"mavros_msgs"`
	Size        uint32
	Success     bool
	RErrno      int32
}

type FileOpen struct {
	msg.Package `ros:"mavros_msgs"`
	FileOpenReq
	FileOpenRes
}
