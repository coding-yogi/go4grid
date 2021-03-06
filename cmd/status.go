package cmd

import (
	"fmt"
	"os"

	"github.com/coding-yogi/go4grid/grid"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "gets the current state of selenium 4 grid",
	Run: func(cmd *cobra.Command, args []string) {

		namespace, _ := cmd.Flags().GetString(NAMESPACE)

		grid := grid.NewGrid(namespace)
		browsers := [3]string{HUB, CHROME, FIREFOX}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Namespace", "Pods", "Created", "Image"})

		for _, browser := range browsers {

			deployments, _ := grid.GetDeployment(browser)

			if len(deployments.Items) > 0 {
				deployment := deployments.Items[0]
				name := deployment.GetName()
				requiredReplicas := deployment.Status.Replicas
				availableReplicas := deployment.Status.AvailableReplicas
				createdate, _ := deployment.ObjectMeta.CreationTimestamp.MarshalText()
				image := deployment.Spec.Template.Spec.Containers[0].Image

				table.Append([]string{name, namespace, fmt.Sprintf("%d/%d", availableReplicas, requiredReplicas), string(createdate), image})
			}
		}

		table.Render()
	},
}

func init() {
	statusCmd.Flags().StringVar(&Namespace, NAMESPACE, "default", "kube namespace")
}
