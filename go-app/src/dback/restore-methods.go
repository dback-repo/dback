package main

import (
	"io/ioutil"
	"sync"
)

//list of saved containers == list of folders in /backup
func getBackupsContainerList() []string {
	var res []string
	files, err := ioutil.ReadDir(`/backup`)
	check(err)
	for _, curFile := range files {
		if curFile.IsDir() {
			res = append(res, curFile.Name())
		}
	}
	return res
}

func restoreContainers(containers []string) {
	var wg sync.WaitGroup
	wg.Add(len(containers))

	for _, curContainer := range containers {
		go restoreContainer(curContainer, &wg)
	}

	wg.Wait()
}

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// //https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
// func Untar(r io.Reader, dst string) error {
// 	tr := tar.NewReader(r)

// 	for {
// 		header, err := tr.Next()

// 		switch {

// 		// if no more files are found return
// 		case err == io.EOF:
// 			return nil

// 		// return any other error
// 		case err != nil:
// 			return err

// 		// if the header is nil, just skip it (not sure how this happens)
// 		case header == nil:
// 			continue
// 		}

// 		// the target location where the dir/file should be created
// 		target := filepath.Join(dst, header.Name)

// 		// the following switch could also be done using fi.Mode(), not sure if there
// 		// a benefit of using one vs. the other.
// 		// fi := header.FileInfo()

// 		// check the file type
// 		switch header.Typeflag {

// 		// if its a dir and it doesn't exist create it
// 		case tar.TypeDir:
// 			if _, err := os.Stat(target); err != nil {
// 				if err := os.MkdirAll(target, 0755); err != nil {
// 					return err
// 				}
// 			}

// 		// if it's a file create it
// 		case tar.TypeReg:
// 			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
// 			if err != nil {
// 				return err
// 			}

// 			// copy over contents
// 			if _, err := io.Copy(f, tr); err != nil {
// 				return err
// 			}

// 			// manually close here after each file operation; defering would cause each file close
// 			// to wait until all operations have completed.
// 			f.Close()
// 		}
// 	}
// }

// func backupMount(c types.Container, m types.MountPoint, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	cli, err := client.NewEnvClient()
// 	check(err)
// 	defer cli.Close()

// 	check(os.MkdirAll(`/backup/`+c.Names[0]+m.Destination, 0664))

// 	reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, m.Destination)
// 	check(err)

// 	lastSlashIdx := strings.LastIndex(m.Destination, `/`)
// 	// if lastSlashIdx > 0 {
// 	// 	lastSlashIdx--
// 	// }
// 	log.Println(`lastSlashIdx`, lastSlashIdx)
// 	destParent := m.Destination[:lastSlashIdx] // /var/www/lynx -> /var/www
// 	if destParent == `` {
// 		destParent = `/`
// 	}
// 	log.Println(`dest`, m.Destination)
// 	log.Println(`destParent`, destParent)

// 	check(Untar(reader, `/backup/`+c.Names[0]+destParent))
// 	log.Println(c.Names[0] + m.Destination)
// }

func restoreContainer(containerName string, wg *sync.WaitGroup) {
	defer wg.Done()

	// cli, err := client.NewEnvClient()
	// check(err)
	// defer cli.Close()

	// if c.State == `running` {
	// 	if len(c.Mounts) > 0 {
	// 		inspect, err := cli.ContainerInspect(context.Background(), c.ID)
	// 		check(err)

	// 		if inspect.HostConfig.AutoRemove == false {
	// 			timeout := time.Minute
	// 			check(cli.ContainerStop(context.Background(), c.ID, &timeout))

	// 			var wgMount sync.WaitGroup
	// 			wgMount.Add(len(c.Mounts))
	// 			for _, curMount := range c.Mounts {
	// 				go backupMount(c, curMount, &wgMount)
	// 			}
	// 			wgMount.Wait()

	// 			check(cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}))
	// 		}

	// 	}
	// }
}
