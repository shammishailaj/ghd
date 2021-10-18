/*
   Copyright © 2020 Shammi Shailaj <shammi.shailaj@healthians.com>

   Licensed under the HLT License, Version 0.0.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://dogo.healthians.com/licenses/LICENSE-0.0.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package cmd

import (
	"github.com/phayes/hookserve/hookserve"
	"github.com/shammishailaj/ghd/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// cleanCmd represents the cleanCmd command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Git Webhook Server",
	Long:  `Used to run Git Webhook Server`,
	Run: func(cmd *cobra.Command, args []string) {
		server := hookserve.NewServer()
		serverPortStr := os.Getenv("GITHUB_WEBHOOK_SERVER_PORT")
		serverPort, serverPortErr := strconv.Atoi(serverPortStr)
		if serverPortErr != nil {
			log.Errorf("Error Converting port number %s from string to integer. Reason: %s", serverPortStr, serverPortErr)
			serverPort = utils.GITHUB_WEBHOOK_SERVER_PORT_DEFAULT
		}
		server.Port = serverPort
		server.Secret = os.Getenv("GITHUB_WEBHOOK_SECRET")
		log.Infof("Starting Github Webhook Server on port %d", serverPort)
		server.GoListenAndServe()

		// Everytime the server receives a webhook event, print the results
		for event := range server.Events {
			log.Infof("Action: %s", event.Action)
			log.Infof("BaseBranch: %s", event.BaseBranch)
			log.Infof("BaseOwner: %s", event.BaseOwner)
			log.Infof("BaseRepo: %s", event.BaseRepo)
			log.Infof("Branch: %s", event.Branch)
			log.Infof("Commit: %s", event.Commit)
			log.Infof("Owner: %s", event.Owner)
			log.Infof("Repo: %s", event.Repo)
			log.Infof("Type: %s", event.Type)
			log.Infof("event.String() %s", event.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
