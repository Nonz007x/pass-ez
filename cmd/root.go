/*
Copyright © 2024 Nonz007x <nathanon.sp00@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/myapp/utils"
)

var rootCmd = &cobra.Command{
	Use:   "passez",
	Short: "Manage your passwords ez",
	Long: `Simple and easy, just like your mom!`,
	Run: rootCommand,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	
	rootCmd.Flags().BoolP("ez", "e", false, "Open easy mode")
}

func rootCommand(cmd *cobra.Command, args []string) {
	
	easyMode, err := cmd.Flags().GetBool("ez")
	if err != nil {
		cmd.Println("Error retrieving flag:", err)
		return
	}

	if easyMode {
		utils.EasyMode()
	} else {
		cmd.Help()
	}
}