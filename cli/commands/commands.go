package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// RootCmd sets task command config
	RootCmd = &cobra.Command{
		Use: "resume",
		Run: func(cmd *cobra.Command, args []string) {
			versionCommand(cmd, args)
			cmd.Usage()
		},
	}
)

// Run runs CLI action
func Run() {
	if err := RootCmd.Execute(); err != nil {
		Exit(fmt.Errorf("failed to run: %s", err.Error()), 1)
	}
}

// Exit finishs requests
func Exit(err error, codes ...int) {
	var code int
	if len(codes) > 0 {
		code = codes[0]
	} else {
		code = 2
	}
	fmt.Println(err)
	os.Exit(code)
}

// IsExist checks to exist file or directory
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
