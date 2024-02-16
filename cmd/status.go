/*
Copyright © 2022 GRISARD Dimitri dimitri.grisard03@gmail.com

*/
package cmd

import (
	"fmt"
	"log"

	ssh "github.com/dim-ops/go-ucce/utils"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status of every services",
	Run: func(cmd *cobra.Command, args []string) {

		allowedType = []string{"cuic", "finesse", "vvb"}

		err := checkUcceType(allowedType, typeOf)
		if err != nil {
			log.Fatal(err.Error())
		}

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
	rootCmd.AddCommand(statusCmd)
}
