package cmd

import (
	"fmt"
	"log"
	"strconv"
	"task/db"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(completedCmd)
}

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "Marks the tasks from to-do list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var dur time.Duration
		if len(args) == 0 {
			dur = 12 * time.Hour
		} else if len(args) > 1 {
			fmt.Println("Error: only one argument (time in hours) expected")
			return
		} else {
			temp, _ := strconv.Atoi(args[0])
			dur = time.Duration(temp) * time.Hour
		}
		err := db.GetCompletedTasks(dur)
		if err != nil {
			log.Fatalln(err)
		}
	},
}
