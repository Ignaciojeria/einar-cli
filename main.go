package main

import (
	_ "github.com/Ignaciojeria/einar/app/adapter/in/cli"
	_ "github.com/Ignaciojeria/einar/app/adapter/in/controller"
	"github.com/Ignaciojeria/einar/app/shared/archetype/cmd"
)

func main() {
	cmd.RootCmd.Execute()
}
