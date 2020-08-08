package cmd

import (
	"fmt"
	"log"
	"strconv"
	"task/db"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks the tasks from to-do list as complete",
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
				fmt.Println("Invalid task number:", id)
				continue
			} else {
				err := db.CompleteTask(tasks[id-1].Key)
				if err != nil {
					fmt.Printf("Error marking task %d as complete: %s\n", id, err)
				}
				fmt.Printf("You have completed the \"%s\" task\n", tasks[id-1].Value)
			}
		}
	},
}
