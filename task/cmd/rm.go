package cmd

import (
	"fmt"
	"log"
	"strconv"
	"task/db"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes the tasks from the to-do list",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			key, err := strconv.Atoi(arg)
			if err != nil {
				log.Fatalln(err)
			} else {
				ids = append(ids, key)
			}
		}
		tasks, err := db.ListTasks()
		if err != nil {
			log.Fatalln(err)
		}
		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Task", id, " does not exist")
				continue
			} else {
				err := db.DeleteTask(tasks[id-1].Key)
				if err != nil {
					fmt.Printf("Error removing task %d: %s\n", id, err)
				}
				fmt.Printf("Removed the task \"%s\" \n", tasks[id-1].Value)
			}
		}
	},
}
