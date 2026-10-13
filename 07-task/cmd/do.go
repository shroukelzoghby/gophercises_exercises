package cmd

import (
	"fmt"
	"strconv"

	"github.com/shroukelzoghby/gophercises_exercises/task/db"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task number to complete.")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Failed to parse the argument:", args[0])
			return
		}

		tasks, err := db.AllTasks()
		must(err)
		if id <= 0 || id > len(tasks) {
			fmt.Println("Invalid task number.")
			return
		}

		task := tasks[id-1]
		text, err := db.CompleteTask(task.Key)
		must(err)
		fmt.Printf("You have completed the %q task.\n", text)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
