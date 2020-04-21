package docker

import (
	"log"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func MustNewDockerClient() string {
	res, err := NewDockerClient()
	check(err)

	return res
}

func NewDockerClient() (string, error) {
	return ``, nil
}

func CheckWeAreInDocker() {
	log.Println(`TODO: check we are in docker`)
}
