package cmd

import (
	"fmt"
	"sync"

	"github.com/coding-yogi/go4grid/grid"
	"github.com/spf13/cobra"
)

var Chrome, Firefox int32

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start up selenium 4 grid hub and nodes",
	Run: func(cmd *cobra.Command, args []string) {

		//flags
		chromeNodes, _ := cmd.Flags().GetInt32("chrome")
		firefoxNodes, _ := cmd.Flags().GetInt32("firefox")
		namespace, _ := cmd.Flags().GetString("namespace")

		fmt.Printf("go4grid: starting grid deployment for %d chrome nodes and %d firefox nodes\n", chromeNodes, firefoxNodes)

		grid := grid.NewGrid(namespace)

		//Hub Deployment
		err := grid.HandleHubDeployment()
		if err != nil {
			return
		}

		//Create Hub Service
		err = grid.HandleHubService()
		if err != nil {
			return
		}

		//Node Deployment
		var wg sync.WaitGroup
		wg.Add(2)
		go triggerNodeDeployment(&grid, "chrome", chromeNodes, &wg)
		go triggerNodeDeployment(&grid, "firefox", firefoxNodes, &wg)
		wg.Wait()
	},
}

func triggerNodeDeployment(grid *grid.Grid, browser string, expReplicas int32, wg *sync.WaitGroup) {
	grid.HandleNodeDeployment(browser, expReplicas)
	wg.Done()
}

func init() {
	startCmd.Flags().Int32Var(&Chrome, "chrome", 1, "number of chrome nodes")
	startCmd.Flags().Int32Var(&Firefox, "firefox", 1, "number of firefox nodes")
	startCmd.Flags().StringVar(&Namespace, "namespace", "default", "kube namespace")
}
