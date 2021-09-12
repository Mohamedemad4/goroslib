//nolint:golint,lll
package mavros_msgs

import (
	"github.com/aler9/goroslib/pkg/msg"
)

type FileListReq struct {
	msg.Package `ros:"mavros_msgs"`
	DirPath     string
}

type FileListRes struct {
	msg.Package `ros:"mavros_msgs"`
	List        []FileEntry
	Success     bool
	RErrno      int32
}

type FileList struct {
	msg.Package `ros:"mavros_msgs"`
	FileListReq
	FileListRes
}
