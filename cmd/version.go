package cmd

import (
	"fmt"

	"github.com/DanielPickens/keeper/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Keeper",
	Long:  `This command will print Keeper's version number and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		runVersion()
	},
}

func NewVersionCommand() *cobra.Command {
	return versionCmd
}

func runVersion() {
	fmt.Println(fmt.Sprintf("Keeper version %s", version.GetVersion()))
}
