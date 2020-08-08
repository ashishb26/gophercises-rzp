package cmd

import (
	"fmt"
	"log"
	"task/db"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the tasks in your to-do list",
	Run: func(cmd *cobra.Command, args []string) {
		taskList, err := db.ListTasks()
		if err != nil {
			log.Fatalln(err)
		}

		if len(taskList) == 0 {
			fmt.Println("You have no tasks reamining")
		} else {
			for i, task := range taskList {
				fmt.Printf("%d. %s\n", i+1, task.Value)
			}
		}
	},
}
