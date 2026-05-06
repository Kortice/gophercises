package main

import (
	"fmt"
	"os"

	"github.com/Kotrice/gophercises/task/cmd"
	"github.com/Kotrice/gophercises/task/db"
)

func main() {
	filePath := "./tasks.db"
	must(db.Init(filePath))
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
