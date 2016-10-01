package db

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	dc "github.com/fsouza/go-dockerclient"
)

type container struct {
	dc.APIContainers
	client *dockerclient
	Name   string
}

func newContainer(rawContainer dc.APIContainers, client *dockerclient) container {
	c := container{
		rawContainer,
		client,
		strings.TrimPrefix(rawContainer.Names[0], "/"),
	}
	return c
}

func (c *container) hasLabel(client *dockerclient, label string) bool {
	_, ok := c.Labels[label]

	return ok
}

func (c *container) execCommand(cmd []string, reader io.Reader, writer io.Writer) error {
	stdinFlag := false
	if reader != nil {
		stdinFlag = true
	}

	stdoutFlag := false
	if writer != nil {
		stdoutFlag = true
	}

	execObject, err := client.CreateExec(dc.CreateExecOptions{
		AttachStdin:  stdinFlag,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          stdoutFlag,
		Cmd:          cmd,
		Container:    c.Name,
	})
	if err != nil {
		log.Println("CreateExec: ", err)
		return err
	}

	err = client.StartExec(execObject.ID, dc.StartExecOptions{
		InputStream:  reader,
		OutputStream: writer,
		ErrorStream:  os.Stderr,
	})
	if err != nil {
		log.Println("StartExec: ", err)
		return err
	}

	return nil
}

func (c *container) myHealthcheck(status chan bool) {
	healthyCounter := 1
	for {
		c.refresh()
		fmt.Printf("healthchecking... %s\n", c.Status)
		if strings.Contains(c.Status, "(healthy)") {
			if healthyCounter <= 0 {
				status <- true
				break
			}
			healthyCounter--
		} else if strings.Contains(c.Status, "(unhealthy)") {
			status <- false
			break
		}

		time.Sleep(5 * time.Second)
	}
}

func (c *container) refresh() {
	tmp := c.client.findContainerByID(c.ID)
	c.APIContainers = tmp.APIContainers
}

func (c container) String() string {
	c.refresh()
	output := ""

	output += fmt.Sprintf("Name:\t%s\n", c.Name)
	output += fmt.Sprintf("ID:\t%s\n", c.ID)
	output += fmt.Sprintf("Image:\t%s\n", c.Image)
	output += fmt.Sprintf("Status:\t%s\n", c.Status)
	publicPorts := []string{}
	for _, port := range c.Ports {
		publicPorts = append(publicPorts, fmt.Sprintf("%d", port.PublicPort))
	}
	output += fmt.Sprintf("Ports:\t%s\n", strings.Join(publicPorts, ","))

	return output
}
