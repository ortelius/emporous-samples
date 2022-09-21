package main

import (
	"github.com/spf13/cobra"

	"github.com/uor-framework/samples/clients/sample-go/cmd/sample-client/commands"
)

func main() {
	app := commands.NewRootCmd()
	cobra.CheckErr(app.Execute())
}
