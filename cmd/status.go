/*
Copyright © 2022 GRISARD Dimitri dimitri.grisard03@gmail.com

*/
package cmd

import (
	"fmt"

	ssh "github.com/DimProject/ucce-cisco/utils"
	"github.com/spf13/cobra"
)

var (
	host, user, password, typeOf string
	port                         uint16
)

// type Parameters struct {
// 	Host     string
// 	User     string
// 	Password string
// }

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Create a todo",
	Long:  `This command will create todo`,
	Run: func(cmd *cobra.Command, args []string) {

		// Storing task in backend calling my-todos REST API
		fmt.Println("---- begin conn.Connect ----")
		conn, _ := ssh.Connect(user, password, host)
		fmt.Println("---- end conn.Connect ----")
		fmt.Println("---- begin conn.SendCommands ----")
		err := conn.SendCommands("show status")
		if err != nil {
			fmt.Println(fmt.Errorf("failed to send command: %s", err))
		}
		fmt.Println("Task created with ID:", conn)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVarP(&host, "host", "a", "", "Hostname or IP address targeted")
	statusCmd.Flags().Uint16VarP(&port, "port", "p", 22, "Ssh port used")
	statusCmd.Flags().StringVarP(&user, "user", "u", "", "User used to login")
	statusCmd.Flags().StringVarP(&password, "password", "x", "", "Password used to login")
	statusCmd.Flags().StringVarP(&typeOf, "typeOf", "t", "", "Type of UCCE Instance (Finesse, Cuic...)")
}
