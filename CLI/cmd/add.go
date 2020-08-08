package cmd

import (
	"fmt"
	"log"
	"strings"
	"task/db"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Complete and remove a task from your to-do list",
	Run: func(cmd *cobra.Command, args []string) {
		taskName := strings.Join(args, " ")
		err := db.AddTask(taskName)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf("Added %s to your task list", taskName)
		}
	},
}
