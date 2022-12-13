/*
Copyright © 2022 GRISARD Dimitri dimitri.grisard03@gmail.com

*/
package cmd

import (
	"fmt"

	ssh "github.com/DimProject/ucce-cisco/utils"
	"github.com/spf13/cobra"
)

var replicationCmd = &cobra.Command{
	Use:   "replication",
	Short: "Get status replication",
	Run: func(cmd *cobra.Command, args []string) {

		// Storing task in backend calling my-todos REST API
		conn, err := ssh.Connect(user, password, host)
		if err != nil {
			fmt.Print(err.Error())
		}
		err = conn.SendCommands("show status")
		if err != nil {
			fmt.Println(fmt.Errorf("failed to send command: %s", err))
		}
		fmt.Println("Task created with ID:", conn)
	},
}

func init() {
	rootCmd.AddCommand(replicationCmd)

}
