//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type FileRenameReq struct {
	msg.Package `ros:"mavros_msgs"`
	OldPath     string
	NewPath     string
}

type FileRenameRes struct {
	msg.Package `ros:"mavros_msgs"`
	Success     bool
	RErrno      int32
}

type FileRename struct {
	msg.Package `ros:"mavros_msgs"`
	FileRenameReq
	FileRenameRes
}
