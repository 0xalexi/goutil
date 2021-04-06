package log

import (
	"os"
)

func truncateLog() {
	outlog.Close()
	err := os.Truncate(Basename+".log", 0)
	if err != nil {
		panic(err)
	}
	outlog, err = os.OpenFile(Basename+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
}
