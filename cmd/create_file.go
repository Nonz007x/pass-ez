/*
Copyright © 2024 Nonz007x <nathanon.sp00@gmail.com>

*/

package cmd

import (
	// "fmt"
	"github.com/spf13/cobra"
	// "github.com/spf13/myapp/utils"
)

var createfileCmd = &cobra.Command{
	Use:   "create [filename]",
	Short: "Create a file",
	Long: `Create a new JSON file with the specified filename.`,
	Args:  cobra.ExactArgs(1),
	Run: createFile,
}

func init() {
	rootCmd.AddCommand(createfileCmd)

}

func createFile(cmd *cobra.Command, args []string) {
	// fileName := args[0]

	// err := utils.CreateFile(fileName, []byte("dummy"))
	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }

	// fmt.Printf("File '%s' created successfully!\n", fileName)
}
