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

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Gracefully shutdown UCCE appliance",
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
		if host == "cusp" {
			err = conn.SendCommands("shutdown")
			if err != nil {
				fmt.Println(fmt.Errorf("failed to send command: %s", err))
			}
		} else {
			err = conn.SendCommands("shutdown\n")
			if err != nil {
				fmt.Println(fmt.Errorf("failed to send command: %s", err))
			}
		}
		fmt.Println("Task created with ID:", conn)
	},
}

func init() {
	rootCmd.AddCommand(shutdownCmd)
}
