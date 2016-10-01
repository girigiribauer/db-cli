package db

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	dc "github.com/fsouza/go-dockerclient"
)

// ContainerOptions for settings
type ContainerOptions struct {
	Incrementation  bool
	Name            string
	DBName          string
	DBUser          string
	DBPass          string
	Filepath        string
	Image           string
	Tag             string
	ContainerPrefix string
	ContainerMax    int
	Directory       string
}

var client *dockerclient

func show(opt *ContainerOptions, driver driver) {
	fmt.Println("===== Show containers =====")

	containerList := client.showContainerList()
	for _, container := range containerList {
		fmt.Println(container)
	}
}

func create(opt *ContainerOptions, driver driver) {
	fmt.Println("===== Create container =====")

	containerName := opt.Name
	if opt.Incrementation {
		containerName = client.findBlankName(opt.ContainerPrefix, opt.ContainerMax)
		if containerName == "" {
			log.Println("create: failed findBlankName()")
			return
		}
	}

	container := client.findContainerByName(containerName)
	if container != nil {
		log.Println("create: already existing container", containerName)
		return
	}

	err := client.downloadImage(driver.imageName())
	if err != nil {
		log.Println("downloadImage: failed download", err)
	}

	publicPortKey := dc.Port(driver.portString())
	publicPortValue := driver.portString()
	if client.isUsingPort(driver.portNumber()) {
		publicPortValue = ""
	}

	isRestore := opt.Filepath != ""

	var healthConfig *dc.HealthConfig
	if isRestore {
		healthConfig = &dc.HealthConfig{
			Test: []string{
				"CMD-SHELL",
				driver.healthcheckCommand(),
			},
			Interval: 1 * time.Second,
			Timeout:  120 * time.Second,
			Retries:  0,
		}
	}

	container, err = client.newContainer(dc.CreateContainerOptions{
		Name: containerName,
		Config: &dc.Config{
			Env:          driver.envString(),
			Healthcheck:  healthConfig,
			Image:        driver.imageName(),
			AttachStderr: true,
			AttachStdout: true,
		},
		HostConfig: &dc.HostConfig{
			PortBindings: map[dc.Port][]dc.PortBinding{
				publicPortKey: []dc.PortBinding{
					dc.PortBinding{
						HostIP:   "",
						HostPort: publicPortValue,
					},
				},
			},
		},
	})
	if err != nil {
		log.Println("create: ", err)
		return
	}

	if isRestore {
		fmt.Println("start restore")

		// there's no way...
		status := make(chan bool)
		go container.myHealthcheck(status)
		select {
		case ok := <-status:
			if !ok {
				log.Println("failed healthcheck")
			}
			break
		}

		fp, reader := openFileReader(opt.Filepath)
		defer fp.Close()

		err = container.execCommand(driver.restoreCommands(), reader, nil)
		if err != nil {
			log.Println("execCommand: ", err)
		}
		fmt.Println("end restore")
	}

	fmt.Println(container)
}

func delete(opt *ContainerOptions, driver driver) {
	fmt.Println("===== delete container =====")

	containerList := []container{}
	containerName := opt.Name

	if containerName == "*" {
		containerList = client.showContainerList()
	} else {
		if opt.Incrementation {
			containerName = client.findFollowingName(opt.ContainerPrefix, opt.ContainerMax)
			if containerName == "" {
				log.Println("delete: no more auto naming containers")
				return
			}
		}

		container := client.findContainerByName(containerName)
		if container == nil {
			log.Println("delete: no existing container", containerName)
			return
		}
		containerList = append(containerList, *container)
	}

	for _, c := range containerList {
		fmt.Println(c)

		err := client.RemoveContainer(dc.RemoveContainerOptions{
			ID:    c.ID,
			Force: true,
		})
		if err != nil {
			log.Println("RemoveContainer: ", err)
			return
		}
	}
}

func dump(opt *ContainerOptions, driver driver) {
	fmt.Println("===== dump container =====")

	containerName := opt.Name
	if opt.Incrementation {
		containerName = client.findFollowingName(opt.ContainerPrefix, opt.ContainerMax)
		if containerName == "" {
			fmt.Println("no target container")
			return
		}
	}

	container := client.findContainerByName(containerName)
	if container == nil {
		log.Println("dump: no existing container", containerName)
		return
	}

	path := opt.Filepath
	if path == "" {
		path = filepath.Join(opt.Directory, fmt.Sprintf("%s.sql", containerName))
	}
	usr, _ := user.Current()
	path = strings.Replace(path, "~", usr.HomeDir, -1)
	path, err := filepath.Abs(path)
	if err != nil {
		log.Println("filepath.Abs: ", err)
		return
	}

	fp, writer := openFileWriter(path)
	defer fp.Close()

	err = container.execCommand(driver.dumpCommands(), nil, writer)
	if err != nil {
		log.Println("execCommand: ", err)
		return
	}

	writer.Flush()

	fmt.Printf("Dump:\t%s\n", path)
	fmt.Println(container)
}

// DB cli entrypoint
func DB(action string, opt *ContainerOptions) {
	var err error
	var driver driver

	client, err = newDockerClient()
	if err != nil {
		fmt.Println("newDockerClient: ", err)
		fmt.Println("running Docker? or installed?")
		return
	}

	switch {
	case strings.Contains(opt.Image, "mariadb"):
		driver = newMariaDBDriver(opt.Image, opt.Tag, opt.DBName, opt.DBUser, opt.DBPass)
	case strings.Contains(opt.Image, "mysql"):
		driver = newMySQLDriver(opt.Image, opt.Tag, opt.DBName, opt.DBUser, opt.DBPass)
	default:
		log.Println("this image is not supported: ", opt.Image)
		return
	}

	switch action {
	case "list":
		show(opt, driver)
	case "delete":
		delete(opt, driver)
	case "dump":
		dump(opt, driver)
	default:
		create(opt, driver)
	}
}
