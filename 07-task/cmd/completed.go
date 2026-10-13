package cmd

import (
	"fmt"

	"github.com/shroukelzoghby/gophercises_exercises/task/db"
	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List tasks completed today",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.TasksCompletedToday()
		must(err)
		if len(tasks) == 0 {
			fmt.Println("You have not completed any tasks today.")
			return
		}
		fmt.Println("You have finished the following tasks today:")
		for _, task := range tasks {
			fmt.Printf("- %s\n", task.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
