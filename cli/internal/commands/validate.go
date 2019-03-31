package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/usk81/go-resume/shema"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "validates JSON-Resume",
		Long:  "validates JSON-Resume",
		Run:   validateCommand,
	}
)

func init() {
	RootCmd.AddCommand(validateCmd)
}

func validateCommand(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "resume.json"
	}

	if err := validateAction(path); err != nil {
		Exit(err, 1)
	}
}

func validateAction(path string) error {
	if !IsExist(path) {
		return fmt.Errorf("%s doesn't exist", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to open file; %s", path))
	}
	defer f.Close()

	var r shema.Resume
	if err = json.NewDecoder(f).Decode(&r); err != nil {
		return err
	}

	if es := shema.Validation(&r); es != nil {
		if _, ok := es.(*validator.InvalidValidationError); ok {
			fmt.Println(es)
			return nil
		}
		fmt.Println("validation error:")
		for _, e := range es.(validator.ValidationErrors) {
			fmt.Printf("key: %s, tag: %s, value: %+v", e.Namespace(), e.Tag(), e.Value())
			fmt.Println()
		}
	}

	return nil
}
