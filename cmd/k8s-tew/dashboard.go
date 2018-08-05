package main

import (
	"fmt"
	"os"
	"time"

	"github.com/darxkies/k8s-tew/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var openBrowser bool

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Retrieves and shows the dashboard token",
	Long:  "Retrieves and shows the dashboard token",
	Run: func(cmd *cobra.Command, args []string) {
		if error := Bootstrap(false); error != nil {
			log.WithFields(log.Fields{"error": error}).Error("dashboard failed")

			os.Exit(-1)
		}

		kubectlCommand := fmt.Sprintf("%s --kubeconfig %s", _config.GetFullLocalAssetFilename(utils.KUBECTL_BINARY), _config.GetFullLocalAssetFilename(utils.ADMIN_KUBECONFIG))
		dashboardKeyCommand := fmt.Sprintf("%s -n kube-system describe secret $(%s -n kube-system get secret | grep admin-user | awk '{print $1}') | grep token: | awk '{print $2}'", kubectlCommand, kubectlCommand)

		output, error := utils.RunCommandWithOutput(dashboardKeyCommand)
		if error != nil {
			log.WithFields(log.Fields{"error": error}).Error("dashboard failed")

			os.Exit(-2)
		}

		fmt.Printf(output)

		if openBrowser {
			fmt.Printf("\nOpening web browser...\n")

			time.Sleep(3 * time.Second)

			ip, error := _config.GetWorkerIP()
			if error != nil {
				log.WithFields(log.Fields{"error": error}).Error("failed to get worker ip address")

				os.Exit(-3)
			}

			if _, error := utils.RunCommandWithOutput(fmt.Sprintf("xdg-open https://%s:%d", ip, _config.Config.DashboardPort)); error != nil {
				log.WithFields(log.Fields{"error": error}).Error("failed to open the web browser")

				os.Exit(-4)
			}

		}
	},
}

func init() {
	dashboardCmd.Flags().BoolVarP(&openBrowser, "open-browser", "o", false, "Open the web browser with a delay of 3 seconds")
	RootCmd.AddCommand(dashboardCmd)
}