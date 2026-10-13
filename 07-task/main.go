package main

import (
	"github.com/shroukelzoghby/gophercises_exercises/task/cmd"
	"github.com/shroukelzoghby/gophercises_exercises/task/db"
)

func main() {
	defer db.Close()
	cmd.Execute()
}
