//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type FileRemoveDirReq struct {
	msg.Package `ros:"mavros_msgs"`
	DirPath     string
}

type FileRemoveDirRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
	RErrno      int32
}

type FileRemoveDir struct {
	msg.Package `ros:"mavros_msgs"`
	FileRemoveDirReq
	FileRemoveDirRes
}
