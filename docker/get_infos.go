package docker

import (
	"errors"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

func GetAllContainers() []container.Summary {
	containers, err := dc.ContainerList(dcCtx, container.ListOptions{All: true})
	if err != nil {
		lib.PrintError("Failed to list containers:", err)
	}
	return containers
}

func GetSSMContainers() (fec []Container, bec []Container) {
	containers := GetAllContainers()

	// 遍历容器列表并解析镜像名称
	for _, c := range containers {
		imageNameParts := strings.Split(c.Image, ":")
		if len(imageNameParts) != 2 {
			continue
		}

		ssmContainer := Container{
			Name:            c.Names[0][1:],
			DockerContainer: c,
			ImageName:       imageNameParts[0],
			Version:         imageNameParts[1],
		}

		// 继承原来的 private
		if len(c.Ports) > 0 && c.Ports[0].IP == "127.0.0.1" {
			ssmContainer.Private = true
		}

		switch imageNameParts[0] {
		case vars.DockerNameFE:
			ssmContainer.ContainerType = vars.ContainerTypeFE
			fec = append(fec, ssmContainer)
		case vars.DockerNameBE:
			ssmContainer.ContainerType = vars.ContainerTypeBE
			bec = append(bec, ssmContainer)
		}
	}
	return
}

func GetContainerByName(n string) (Container, bool) {
	// 检查是否存在正在运行的名字为n的容器
	fel, bel := GetSSMContainers()

	for _, fec := range fel {
		if fec.Name == n {
			return fec, true
		}
	}

	for _, bec := range bel {
		if bec.Name == n {
			return bec, true
		}
	}

	return Container{}, false
}

func (c *Container) GetPortInfo() (info PortInfo, error error) {
	if c.DockerContainer.State != "running" {
		error = errors.New("container not found or not running")
		return
	}

	info.Public = strconv.Itoa(int(c.DockerContainer.Ports[0].PublicPort))
	info.Private = strconv.Itoa(int(c.DockerContainer.Ports[0].PrivatePort))
	info.Type = c.DockerContainer.Ports[0].Type
	return
}
