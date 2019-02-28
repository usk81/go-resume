package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version    = "0.0.1"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Show version",
		Run:   versionCommand,
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func versionCommand(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}
