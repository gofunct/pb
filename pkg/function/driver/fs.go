package driver

import (
	"bytes"
	"github.com/gofunct/pb/pkg/logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// PackageName determines the package name from sourceFile if it is within $GOPATH
func (c *Command) PackageName(sourceFile string) (string, error) {
	if filepath.Ext(sourceFile) != ".go" {
		return "", errors.New("sourcefile must end with .go")
	}
	sourceFile, err := filepath.Abs(sourceFile)
	if err != nil {
		logging.L.Fatal( "Could not convert to absolute path: %s", sourceFile)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("Environment variable GOPATH is not set")
	}
	paths := strings.Split(gopath, string(os.PathListSeparator))
	for _, path := range paths {
		srcDir := filepath.Join(path, "src")
		srcDir, err := filepath.Abs(srcDir)
		if err != nil {
			continue
		}

		//log.Printf("srcDir %s sourceFile %s\n", srcDir, sourceFile)
		rel, err := filepath.Rel(srcDir, sourceFile)
		if err != nil {
			continue
		}
		return filepath.Dir(rel), nil
	}
	return "", errors.New("sourceFile not reachable from GOPATH")
}

// Template reads a go template and writes it to dist given data.
func (c *Command) ProcessAndMove(src string, dest string) {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		logging.L.Fatalf(  "Could not read file %s\n%v\n", src, err)

	}

	tpl := template.New("t")
	tpl, err = tpl.Parse(string(content))
	if err != nil {
		logging.L.Fatalf("Could not parse template %s\n%v\n", src, err)
	}

	f, err := os.Create(dest)
	if err != nil {
		logging.L.Fatalf( "Could not create file for writing %s\n%v\n", dest, err)
	}
	defer f.Close()
	err = tpl.Execute(f, viper.AllSettings())
	if err != nil {
		logging.L.Fatalf("Could not execute template %s\n%v\n", src, err)
	}
}

// StrTemplate reads a go template and writes it to dist given data.
func (c *Command) ProcessString(src string) (string, error) {
	tpl := template.New("t")
	tpl, err := tpl.Parse(src)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, viper.AllSettings())
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}