//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type FileTruncateReq struct {
	msg.Package `ros:"mavros_msgs"`
	FilePath    string
	Length      uint64
}

type FileTruncateRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
	RErrno      int32
}

type FileTruncate struct {
	msg.Package `ros:"mavros_msgs"`
	FileTruncateReq
	FileTruncateRes
}
