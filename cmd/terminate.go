package cmd

import (
	"fmt"
	"sync"

	"github.com/coding-yogi/go4grid/grid"
	"github.com/spf13/cobra"
)

var terminateCmd = &cobra.Command{
	Use:   "terminate",
	Short: "cleans up selenium 4 grid hub and nodes",
	Run: func(cmd *cobra.Command, args []string) {

		namespace, _ := cmd.Flags().GetString(NAMESPACE)
		fmt.Println("go4grid: terminating all grid components")

		grid := grid.NewGrid(namespace)
		var wg sync.WaitGroup
		wg.Add(4)
		go triggerDeploymentDeletion(&grid, CHROME, &wg)
		go triggerDeploymentDeletion(&grid, FIREFOX, &wg)
		go triggerDeploymentDeletion(&grid, HUB, &wg)
		go triggerHubServiceDeletion(&grid, &wg)
		wg.Wait()
	},
}

func init() {
	terminateCmd.Flags().StringVar(&Namespace, NAMESPACE, "default", "kube namespace")
}

func triggerDeploymentDeletion(grid *grid.Grid, name string, wg *sync.WaitGroup) {
	grid.DeleteDeployment(name)
	wg.Done()
}

func triggerHubServiceDeletion(grid *grid.Grid, wg *sync.WaitGroup) {
	grid.DeleteHubService()
	wg.Done()
}
