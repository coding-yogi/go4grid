package grid

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/coding-yogi/go4grid/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Grid struct {
	k8s KubeApi
}

func NewGrid(namespace string) Grid {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		panic("no .kube file found in home directory")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return Grid{
		KubeApi{
			client:    client,
			namespace: namespace,
		},
	}
}

func (grid *Grid) GetHubService() (*corev1.ServiceList, error) {
	return grid.k8s.GetServices("app=" + resources.HubName)
}

func (grid *Grid) GetDeployment(browser string) (*appsv1.DeploymentList, error) {
	var name string

	if browser == "hub" {
		name = "app=" + resources.HubName
	} else {
		name = "app=" + resources.NodeName + "-" + browser
	}

	return grid.k8s.GetDeployments(name)
}

func (grid *Grid) IsHubRunning() (bool, error) {
	fmt.Println("hub: checking if running")
	deployments, err := grid.GetDeployment("hub")
	if err != nil {
		return false, err
	}

	return len(deployments.Items) > 0 && deployments.Items[0].Status.AvailableReplicas > 0, nil
}

func (grid *Grid) IsHubServiceAvailable() (bool, error) {
	fmt.Println("hub: checking service")
	services, err := grid.GetHubService()
	if err != nil {
		return false, err
	}

	return len(services.Items) > 0, nil
}

func (grid *Grid) CreateHubDeployment() error {
	fmt.Println("hub: initiating deployment")

	result, err := grid.k8s.CreateDeployment(resources.HubDeployment())
	if err != nil {
		fmt.Println("hub: deployment failed: " + err.Error())
		return err
	}

	fmt.Printf("hub: deployment created - %q\n", result.GetObjectMeta().GetName())
	return nil
}

func (grid *Grid) CreateHubService() error {
	fmt.Println("hub: creating service")

	result, err := grid.k8s.CreateService(resources.HubService())
	if err != nil {
		fmt.Println("hub: service creation failed: " + err.Error())
		return err
	}

	fmt.Printf("hub: service created - %q\n", result.GetObjectMeta().GetName())
	return nil
}

func (grid *Grid) CreateNodeDeployment(browser string, replicas int32, hubHost, hubPort string) error {
	fmt.Printf("%s: initiating node deployment\n", browser)

	nodeDeployment := resources.NodeDeployment(browser, replicas, hubHost, hubPort)
	result, err := grid.k8s.CreateDeployment(nodeDeployment)
	if err != nil {
		fmt.Printf("%s: deployment failed: %s\n", browser, err.Error())
		return err
	}

	fmt.Printf("%s: deployment created %q\n", browser, result.GetObjectMeta().GetName())
	return nil
}

func (grid *Grid) DeleteDeployment(browser string) {
	fmt.Printf("%s: cleaning up deployment\n", browser)

	var name string
	if browser == "hub" {
		name = resources.HubName
	} else {
		name = resources.NodeName + "-" + browser
	}

	fmt.Printf("%s: deleting deployment %s\n", browser, name)
	deployments, _ := grid.GetDeployment(browser)

	if len(deployments.Items) > 0 {
		err := grid.k8s.DeleteDeployment(name)
		if err != nil {
			fmt.Printf("%s: deleting deployment failed: %s\n", browser, err.Error())
		}
	} else {
		fmt.Printf("%s: deployment not found\n", browser)
	}
}

func (grid *Grid) DeleteHubService() {
	fmt.Println("hub: cleaning up service")

	services, _ := grid.GetHubService()

	if len(services.Items) > 0 {
		err := grid.k8s.DeleteService(resources.HubName)
		if err != nil {
			fmt.Println("hub: deleting service failed: " + err.Error())
		}
	} else {
		fmt.Println("hub: service not found")
	}
}

func (grid *Grid) UpdateDeployment(deployment *appsv1.Deployment) error {
	_, err := grid.k8s.UpdateDeployment(deployment)

	if err != nil {
		return err
	}

	return nil
}

func (grid *Grid) WaitFor(browser string, replicas int32) {
	fmt.Printf("%s: waiting for pods to come up ...\n", browser)

	var activePods int32 = 0
	timer := 0

	for activePods != replicas && timer < 20 {
		time.Sleep(3 * time.Second)
		timer++

		deployments, _ := grid.GetDeployment(browser)
		activePods = deployments.Items[0].Status.AvailableReplicas
	}

	if timer == 20 {
		fmt.Printf("%s: timed out waiting for nodes to come up. Check state using status command\n", browser)
	}
}

func (grid *Grid) HandleHubDeployment() error {
	//Check is hub is running
	isHubRunning, err := grid.IsHubRunning()
	if err != nil {
		fmt.Println("hub: unable to get deployment status due to error " + err.Error())
		return err
	}

	if isHubRunning {
		fmt.Println("hub: already running. Skipping deployment")
	} else {
		//Initiate Hub deployment
		err = grid.CreateHubDeployment()
		if err != nil {
			fmt.Println("hub: deployment failed " + err.Error())
			return err
		}

		grid.WaitFor("hub", 1) //Wait for a minute for hub to come up , else exit
	}

	return nil
}

func (grid *Grid) HandleNodeDeployment(browser string, expReplicas int32) error {

	nodeDeployments, err := grid.GetDeployment(browser)
	if err != nil {
		fmt.Printf("%s: unable to get node deployment: %s\n", browser, err.Error())
		return err
	}

	if len(nodeDeployments.Items) == 0 {
		fmt.Printf("%s: no existing deployment, initiating new deployment\n", browser)

		//New deployment
		err := grid.CreateNodeDeployment(browser, expReplicas, resources.HubName, "4444")
		if err != nil {
			fmt.Printf("%s: node deployment failed %s\n", browser, err.Error())
			return err
		}
	} else {
		fmt.Printf("%s: existing deployment found\n", browser)

		//Get current stats
		deployment := nodeDeployments.Items[0]
		deploymentStatus := deployment.Status
		fmt.Printf("%s: required pods: %d, ready pods: %d, upcoming pods: %d\n", browser, expReplicas, deploymentStatus.ReadyReplicas, deploymentStatus.UnavailableReplicas)

		scalefactor := expReplicas - (deploymentStatus.ReadyReplicas + deploymentStatus.UnavailableReplicas)

		if scalefactor == 0 {
			fmt.Printf("%s: no more replicas required, wait for any upcoming pods to be available\n", browser)
		} else {
			fmt.Printf("%s: need to scale replicas by: %d, updating deployment\n", browser, scalefactor)
			//Update deployment
			deployment.Spec.Replicas = &expReplicas
			grid.UpdateDeployment(&deployment)
		}
	}

	//Wait for nodes to come up
	grid.WaitFor(browser, expReplicas)
	return nil
}

func (grid *Grid) HandleHubService() error {
	isHubServiceAvailable, err := grid.IsHubServiceAvailable()
	if err != nil {
		fmt.Println("hub: unable to get service status due to error " + err.Error())
		return err
	}

	if isHubServiceAvailable {
		fmt.Println("hub: service already available.")
	} else {
		err = grid.CreateHubService()
		if err != nil {
			fmt.Println("hub: service creation failed " + err.Error())
			return err
		}
	}

	return nil
}
