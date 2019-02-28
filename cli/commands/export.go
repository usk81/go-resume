package commands

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/usk81/go-resume/shema"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "converts from JSON-Resume and exports to file",
		Long:  "converts from JSON-Resume and exports to file",
		Run:   exportCommand,
	}
)

func init() {
	RootCmd.AddCommand(exportCmd)
}

func exportCommand(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		Exit(errors.New("variables are not enough; resume export (template directory path) (output directory path) (json resume file path)"), 1)
	}

	if err := exportAction(args[0], args[1], args[2]); err != nil {
		Exit(err, 1)
	}
}

func exportAction(from string, to string, jp string) error {
	var r shema.Resume
	err := getResumeFromFile(jp, &r)
	if err != nil {
		return err
	}
	if err = shema.Validation(r); err != nil {
		return err
	}
	return OutputHTML(r, from, to)
}

func getResumeFromFile(fp string, r interface{}) (err error) {
	if !IsExist(fp) {
		return fmt.Errorf("%s doesn't exist", fp)
	}
	f, err := os.Open(fp)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to open file; %s", fp))
	}

	if err = json.NewDecoder(f).Decode(&r); err != nil {
		f.Close()
		return err
	}
	f.Close()
	return nil
}

// OutputHTML outputs HTML file and copy assets
func OutputHTML(r shema.Resume, src string, dst string) (err error) {
	rt := regexp.MustCompile(`^resume(\.html)*\.(template|tpl).*$`)

	if d, err := isDir(src); !d || err != nil {
		if err != nil {
			return err
		}
		return errors.Errorf("%s is not directory", src)
	}

	dst = filepath.Join(dst, "resume")
	if err = os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	out := filepath.Join(dst, "index.html")
	if isExist(out) {
		return errors.Errorf("%s already exists", out)
	}
	f, err := os.Create(out)
	if err != nil {
		return errors.Wrapf(err, "fail to open HTML file")
	}
	defer f.Close()

	sd, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, s := range sd {
		srcPath := filepath.Join(src, s.Name())
		b, err := isDir(srcPath)
		if err != nil {
			return err
		}
		if b {
			if err := CopyDir(srcPath, filepath.Join(dst, s.Name())); err != nil {
				return err
			}
		} else {
			if rt.MatchString(s.Name()) {
				if err = CreateHTML(f, r, srcPath); err != nil {
					return err
				}
			} else {
				if err := CopyDir(srcPath, filepath.Join(dst, s.Name())); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// CreateHTML creates HTML data
func CreateHTML(w io.Writer, r shema.Resume, tp string) error {
	tpl := filepath.Base(tp)

	fs := template.FuncMap{
		"commalist": func(args []string) string {
			return strings.Join(args, ",")
		},
	}
	tmpl, err := template.New(tpl).Funcs(fs).ParseFiles(tp)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, r)
}

func isDir(p string) (bool, error) {
	if !IsExist(p) {
		return false, errors.Errorf("directory is not found : %s", p)
	}
	src, err := os.Open(p)
	if err != nil {
		return false, errors.Wrapf(err, "fail to access : %s", p)
	}
	defer src.Close()

	fi, err := src.Stat()
	if err != nil {
		return false, errors.Wrapf(err, "fail to get file information : %s", p)
	}
	return fi.IsDir(), nil
}

func isExist(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

// Copy a file `from` path to `to` path
func Copy(src, dst string) error {
	if !IsExist(src) {
		return fmt.Errorf("source file does not exist")
	}
	if IsExist(dst) {
		return fmt.Errorf("destination already exists")
	}
	sf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("fail to open source file")
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("fail to create source file")
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}

// CopyDir copies a directory source path to destination path
func CopyDir(src, dst string) error {
	walk := func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, srcPath[len(src):])
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}
		return Copy(srcPath, dstPath)
	}
	return filepath.Walk(src, walk)
}
