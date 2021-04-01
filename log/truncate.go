// +build !windows

package log

import (
	"os"
)

func truncateLog() {
	outlog.Seek(0, os.SEEK_SET)
	err := outlog.Truncate(0)
	if err != nil {
		panic(err)
	}
}
