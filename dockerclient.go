package db

import (
	"fmt"
	"log"
	"strings"

	dc "github.com/fsouza/go-dockerclient"
)

const (
	identityContainerLabel = "DBCLI_CREATED"
)

var dockerclientEndpoint string

type dockerclient struct {
	*dc.Client
}

func init() {
	dockerclientEndpoint = "unix:///var/run/docker.sock"
}

func newDockerClient() (*dockerclient, error) {
	client, err := dc.NewClient(dockerclientEndpoint)
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, err
	}

	return &dockerclient{client}, err
}

func validateContainerName(s string) string {
	return strings.Map(func(r rune) rune {
		if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '_' {
			return r
		}
		return -1
	}, s)
}

func (client *dockerclient) newContainer(opt dc.CreateContainerOptions) (*container, error) {
	opt.Config.Labels = map[string]string{
		identityContainerLabel: "1",
	}

	target, err := client.CreateContainer(opt)
	if err != nil {
		log.Println("CreateContainer: ", err)
		return nil, err
	}

	err = client.StartContainer(target.ID, nil)
	if err != nil {
		log.Println("StartContainer: ", err)
		return nil, err
	}

	container := client.findContainerByID(target.ID)

	return container, nil
}

func (client *dockerclient) downloadImage(name string) error {
	result, err := client.ListImages(dc.ListImagesOptions{
		All:    true,
		Filter: name,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	if len(result) > 0 {
		return nil
	}

	fmt.Printf("Pulling %s image...\n", name)

	// without auth on PullImage
	names := strings.Split(name, ":")
	err = client.PullImage(dc.PullImageOptions{
		Repository: names[0],
		Tag:        names[1],
	}, dc.AuthConfiguration{})
	if err != nil {
		log.Println("downloadImage", err)
		return err
	}

	return nil
}

func (client *dockerclient) showContainersAll() map[string]container {
	return client.showContainers(nil)
}

func (client *dockerclient) showContainersByName(value string) map[string]container {
	filters := map[string][]string{}

	filters["name"] = []string{value}

	return client.showContainers(filters)
}

func (client *dockerclient) showContainersByID(value string) map[string]container {
	filters := map[string][]string{}

	filters["id"] = []string{value}

	return client.showContainers(filters)
}

func (client *dockerclient) showContainers(filter map[string][]string) map[string]container {
	list, err := client.ListContainers(dc.ListContainersOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		log.Println("showContainers: ", err)
		return nil
	}

	containers := map[string]container{}
	for _, rawContainer := range list {
		container := newContainer(rawContainer, client)

		if container.hasLabel(client, identityContainerLabel) {
			containers[container.Name] = container
		}
	}

	return containers
}

func (client *dockerclient) showPorts() map[int64]struct{} {
	containers := client.showContainersAll()

	ports := make(map[int64]struct{})

	for _, container := range containers {
		for _, portConfig := range container.Ports {
			ports[portConfig.PublicPort] = struct{}{}
		}
	}
	return ports
}

func (client *dockerclient) isUsingPort(n int64) bool {
	ports := client.showPorts()
	_, ok := ports[n]
	return ok
}

func (client *dockerclient) findBlankName(prefix string, max int) string {
	name := ""
	prefix = validateContainerName(prefix)

	containers := client.showContainersByName(prefix + "[0-9]+")

	for i := 0; i < max; i++ {
		_, ok := containers[fmt.Sprintf("%s%d", prefix, i)]
		if !ok {
			name = fmt.Sprintf("%s%d", prefix, i)
			break
		}
	}

	return name
}

func (client *dockerclient) findFollowingName(prefix string, max int) string {
	name := ""
	prefix = validateContainerName(prefix)

	containers := client.showContainersByName(prefix + "[0-9]+")

	for i := 0; i < max; i++ {
		_, ok := containers[fmt.Sprintf("%s%d", prefix, i)]
		if ok {
			name = fmt.Sprintf("%s%d", prefix, i)
			break
		}
	}

	return name
}

func (client *dockerclient) showContainerList() []container {
	containers := client.showContainersAll()

	list := []container{}

	for _, c := range containers {
		list = append(list, c)
	}

	return list
}

func (client *dockerclient) findContainerByName(name string) *container {
	containers := client.showContainersAll()
	container, ok := containers[name]
	if ok {
		return &container
	}

	return nil
}

func (client *dockerclient) findContainerByID(id string) *container {
	var container *container
	containers := client.showContainersByID(id)

	for _, c := range containers {
		container = &c
		break
	}

	return container
}

func (client *dockerclient) showContainerPortList(name string) []int64 {
	ports := []int64{}

	container := client.findContainerByName(name)
	if container == nil {
		log.Println("showContainerPortList: no containers")
		return []int64{}
	}

	for _, portConfig := range container.Ports {
		ports = append(ports, portConfig.PublicPort)
	}

	return ports
}
