package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/shroukelzoghby/gophercises_exercises/task/db"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI for managing your TODOs.",
}

// Execute runs the root command.
func Execute() {
	must(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initDB)
}

func initDB() {
	home, err := homedir.Dir()
	must(err)

	dbPath := filepath.Join(home, ".tasks.db")
	must(db.Init(dbPath))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
