package main

import (
	"archive/tar"
	"context"
	"io"

	//	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func Untar(r io.Reader, dst string) error {

	// gzr, err := gzip.NewReader(r)
	// if err != nil {
	// 	return err
	// }
	// defer gzr.Close()

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func backupMount(c types.Container, m types.MountPoint, wg *sync.WaitGroup) {
	defer wg.Done()
	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	check(os.MkdirAll(`/backup/`+c.Names[0], 0664))

	reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, m.Destination)
	check(err)

	check(Untar(reader, `/backup/`+c.Names[0]))
}

func backupContainer(c types.Container, wg *sync.WaitGroup) {
	defer wg.Done()

	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	// log.Println(c.State)
	// log.Println(c.Status)

	if c.State == `running` {

		timeout := time.Minute
		check(cli.ContainerStop(context.Background(), c.ID, &timeout))

		var wgMount sync.WaitGroup
		wgMount.Add(len(c.Mounts))
		for _, curMount := range c.Mounts {
			go backupMount(c, curMount, &wgMount)
		}
		wgMount.Wait()

		check(cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}))
	}

}
