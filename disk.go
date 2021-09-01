package goutil

import (
	"fmt"
	"os"

	humanize "github.com/dustin/go-humanize"
	"github.com/minio/minio/pkg/disk"
)

func GetDiskAvailable() (uint64, error) {
	path, _ := os.Getwd()
	di, err := disk.GetInfo(path)
	if err != nil {
		return 0, err
	}
	return di.Free, nil
}

func GetPercentDiskAvailable() (float64, error) {
	path, _ := os.Getwd()
	di, err := disk.GetInfo(path)
	if err != nil {
		return 0, err
	}
	return (float64(di.Total-di.Free) / float64(di.Total)) * 100, nil
}

func GetDiskInfoString() string {
	path, _ := os.Getwd()
	di, err := disk.GetInfo(path)
	if err != nil {
		return fmt.Sprint("GetDiskInfoString err:", err)
	}
	percentage := (float64(di.Total-di.Free) / float64(di.Total)) * 100
	return fmt.Sprintf("%s of %s disk space used (%0.2f%%)\n",
		humanize.Bytes(di.Total-di.Free),
		humanize.Bytes(di.Total),
		percentage,
	)
}
