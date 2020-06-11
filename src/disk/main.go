package main

import (
	"log"

	"github.com/shirou/gopsutil/disk"
)

func main() {
	diskStat, err := disk.Usage("/")
	log.Println(diskStat, err)
}
