package cmd

import (
	"fmt"
	"strconv"

	"github.com/shroukelzoghby/gophercises_exercises/task/db"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a task from your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task number to delete.")
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
		must(db.DeleteTask(task.Key))
		fmt.Printf("You have deleted the %q task.\n", task.Value)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
