package cmd

import (
	"fmt"
	"strings"

	"github.com/shroukelzoghby/gophercises_exercises/task/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		if strings.TrimSpace(task) == "" {
			fmt.Println("Please provide a task description.")
			return
		}
		_, err := db.CreateTask(task)
		must(err)
		fmt.Printf("Added %q to your task list.\n", task)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
