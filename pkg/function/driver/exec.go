package driver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gofunct/pb/pkg/logging"
	"github.com/prometheus/common/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type FunctionHandler func(command *Command) http.HandlerFunc
type FileFunc func(info os.FileInfo)
type Splitter func(data []byte, atEOF bool) (advance int, token []byte, err error)

type Command struct {
	cmd *exec.Cmd
	reader *bufio.Reader
	writer *bufio.Writer
	scanner *bufio.Scanner
	flags 	*pflag.FlagSet
}

func NewCommand(reader io.Reader, writer io.Writer, script string, dir string) *Command {
	if dir == "" {
		dir = "."
	}
	w := bufio.NewWriter(writer)
	r := bufio.NewReader(reader)
	cmd :=  &Command{
		reader: r,
		writer: w,
	}
	cmd.SyncEnv()
	script, err := cmd.ProcessString(script)
	if err != nil {
		log.Warn(err.Error())
	}
	c := &exec.Cmd{
		Path:         "/bin/bash",
		Args: 			[]string{"bash", "-c", script},
		Env:          	os.Environ(),
		Dir:          	dir,
		Stdin:        	reader,
		Stdout:       	writer,
		Stderr:      	writer,
	}

	cmd.cmd = c
	return cmd

}

func (c *Command) Runnable() bool {
	switch  {
	case c.cmd != nil && c.writer != nil:
		return true
	default:
		return false
	}
}

func (c *Command) ScanFrom(reader io.Reader) {
	c.scanner = bufio.NewScanner(reader)
}

func (c *Command) UnmarshalFromConfig(obj interface{}) error {
	return viper.Unmarshal(obj)
}

func (c *Command) Unmarshal(obj interface{}, data []byte) error {
	return json.Unmarshal(data, obj)
}

func (c *Command) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (c *Command) DecodeFromReader(obj interface{}) error  {
	dec := json.NewDecoder(c.reader)
	return dec.Decode(obj)
}

func (c *Command) EncodeFromWriter(obj interface{}) error  {
	enc := json.NewEncoder(c.writer)
	return enc.Encode(obj)
}

// Prompt prompts user for input with default value.
func (c *Command) Prompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

func (c *Command) SyncEnv() {
	env := os.Environ()
	for _, e := range env {
		mapenv := strings.Split(e, "=")
		if err :=  viper.BindEnv(mapenv[0]); err != nil {
			logging.L.Warn(err.Error())
		}
		viper.SetDefault(strings.ToLower(mapenv[0]), mapenv[1])
	}
	for k, v := range viper.AllSettings() {
		if err := os.Setenv(k, v.(string));err != nil {
			logging.L.Warn(err.Error())
		}
	}
}

func (c *Command) WriteString(s string) error {
	_, err := c.writer.WriteString(s)
	return err
}

func (c *Command) Write(b []byte) error {
	_, err := c.writer.Write(b)
	return err
}

func (c *Command) Read(b []byte) error {
	_, err := c.reader.Read(b)
	return err
}

func (c *Command) ReadeBufferSize() int {
	return c.reader.Size()
}

func (c *Command) WriteBufferSize() int {
	return c.writer.Size()
}

func (c *Command) ResetReader() {
	c.reader.Reset(c.reader)
}

func (c *Command) ResetWriter() {
	c.writer.Reset(c.writer)
}

func (c *Command) BufferRemaining() int {
	return c.writer.Available()
}

func (c *Command) ReadToWriter(r io.Reader) {
	c.writer.ReadFrom(r)
}


func (c *Command) ScanFor(s string) {
	for c.scanner.Scan() {
		if strings.Contains(c.scanner.Text(), s) {
			b, _ := json.Marshal(s)
			_, _ = c.reader.Read(b)
		}
	}
}

func (c *Command) ScanAndReplace(replacements ...string) {
	rep := strings.NewReplacer(replacements...)
	for c.scanner.Scan() {
		rep.Replace(c.scanner.Text())
	}
}

func (c *Command) ScanAndReplaceBytes(replacements ...string) {
	rep := strings.NewReplacer(replacements...)
	for c.scanner.Scan() {
		if c.scanner.Err() != nil {
			logging.L.Warn(c.scanner.Err().Error())
			break
		}
		rep.Replace(string(c.scanner.Bytes()))
	}
}


func NewSplitter(base int, bitsize int) Splitter {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanWords(data, atEOF)
		if err == nil && token != nil {
			_, err = strconv.ParseInt(string(token), base, bitsize)
		}
		return advance, token, err
	}
}

func (c *Command) Walk(path string, fileFunc FileFunc) error {
	if er := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileFunc(info)
		return nil
	}); er != nil {
		return er
	}
	return nil
}
