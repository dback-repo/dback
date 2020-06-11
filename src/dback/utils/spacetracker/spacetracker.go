// the package for catch min free space left on device after tracking started
package spacetracker

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/shirou/gopsutil/disk"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type SpaceTracker struct {
	period        time.Duration
	path          string
	StartSpace    uint64
	MinSpaceBytes uint64
}

func NewSpaceTracker(period time.Duration) *SpaceTracker {
	path := getwd()
	res := SpaceTracker{period: period, path: path, StartSpace: getCurrentSpace(path)}
	res.MinSpaceBytes = math.MaxUint64

	go res.trackSpace()

	return &res
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + "\r\n" + err.Error())
	}
}

func (t *SpaceTracker) trackSpace() {
	for {
		currentSpace := getCurrentSpace(t.path)

		if currentSpace < t.MinSpaceBytes {
			t.MinSpaceBytes = currentSpace
		}

		time.Sleep(t.period)
	}
}

func getwd() string {
	res, err := os.Getwd()
	check(err, `cannot lookup current directory`)

	return res
}

func getCurrentSpace(path string) uint64 {
	diskStat, err := disk.Usage(path)
	check(err, `cannot get disk usage`)

	return diskStat.Free
}

func (t *SpaceTracker) MinSpaceMB() int {
	return int(t.MinSpaceBytes / MB)
}

func (t *SpaceTracker) UsedSpaceMB() int {
	return int((t.StartSpace - t.MinSpaceBytes) / MB)
}

func (t *SpaceTracker) PrintReport() {
	log.Println(`Minimal disk space: `, t.MinSpaceMB(), `MB`)
	log.Println(`Used space: `, t.UsedSpaceMB(), `MB`)
}
