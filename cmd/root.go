// Copyright © 2019 Coleman Word <coleman.word@gofunct.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/gofunct/pb/cmd/load"
	"github.com/gofunct/pb/cmd/walk"
	"github.com/gofunct/pb/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (

	// Used for flags.
	in, out, pkg string
	rootCmd                            = &cobra.Command{
		Use:   "pb",
		Short: "A generator for protobuf based applications ",
		Long: `pb is a protocol buffers utility tool based on docker\n
Images:
colemanword/gcloud
colemanword/protoc
colemanword/tools
colemanword/source
colemanword/templates
`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute:%s/n", err)
	}
}

func init() {
	config.Init(rootCmd.Flags())
	config.Init(rootCmd.PersistentFlags())
	{
		rootCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "use Viper for configuration")
		rootCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "use Viper for configuration")
		rootCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "", "use Viper for configuration")
	}

	{
		viper.SetDefault("input", ".")
		viper.SetDefault("output", ".")
	}

	{
		rootCmd.AddCommand(walk.RootCmd)
		rootCmd.AddCommand(load.RootCmd)
	}
}

func ifErrr (err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}