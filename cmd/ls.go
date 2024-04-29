package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"sub-store-manager-cli/docker"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all sub-store docker containers",
	Run: func(cmd *cobra.Command, args []string) {
		listAllSSMContainer()
	},
}

func listAllSSMContainer() {
	fel, bel := docker.GetSSMContainers()

	if len(fel) == 0 && len(bel) == 0 {
		fmt.Println("No Sub-Store Manager Docker Containers found")
		return
	}

	fmt.Println("Sub-Store Manager Docker Containers:")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Type", "ID", "Version", "Port", "Status", "Name")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, c := range append(fel, bel...) {
		var portStr string
		if p, e := c.GetPortInfo(); e != nil {
			portStr = "none"
		} else {
			ip := "-"
			if len(c.DockerContainer.Ports) > 0 {
				ip = c.DockerContainer.Ports[0].IP
			}
			portStr = fmt.Sprintf("%s:%s->%s/%s", ip, p.Public, p.Private, p.Type)
		}
		tbl.AddRow(c.ContainerType, c.DockerContainer.ID[:24], c.Version, portStr, c.DockerContainer.State, c.Name)
	}

	tbl.Print()
}
