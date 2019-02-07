package main

import (
	"fmt"
	"github.com/gofunct/gofs"
	"github.com/gofunct/gofs/assetfs"
	"github.com/gofunct/stencil"
	"github.com/mgutz/str"
)

func tasks(p *gofs.Project) {
	p.Func("test", nil, func(c *gofs.Context) {
		c.Run("go test")
	})

	p.Func("test", stencil.S{"build"}, func(c *stencil.Context) {
		c.Run("go test")
	})

	p.Func("dist", stencil.S{"test", "lint"}, nil)

	p.Func("install", nil, func(c *stencil.Context) {
		c.Run("go get github.com/golang/lint/golint")
		// Run("go get github.com/mgutz/goa")
		c.Run("go get github.com/robertkrimen/godocdown/godocdown")
	})

	p.Func("lint", nil, func(c *stencil.Context) {
		c.Run("golint .")
		c.Run("gofmt -w -s .")
		c.Run("go vet .")
	})

	p.Func("readme", stencil.S{"install"}, func(c *stencil.Context) {
		c.Run("godocdown -o README.md")
		c.Pipe(
			assetfs.Load("./README.md"),
			gofs.Str(str.ReplaceF("--", "\n[stencil](https://github.com/stencil.com)\n", 1)),
			gofs.Write(),
		)
	})

	p.Func("build", nil, func(c *stencil.Context) {
		c.Run("go install", stencil.M{"$in": "cmd/gostencil"})
	})

	p.Func("interactive", nil, func(c *stencil.Context) {
		c.Bash(`
			echo name?
			read name
			echo hello $name
		`)
	})

	p.Func("whoami", nil, func(c *stencil.Context) {
		c.Run("whoami")
	})

	pass := 0
	p.Func("err2", nil, func(*stencil.Context) {
		if pass == 2 {
			stencil.Halt("oh oh")
		}
	})

	p.Func("err", stencil.S{"err2"}, func(*stencil.Context) {
		pass++
		if pass == 1 {
			return
		}
		stencil.Halt("foo err")
	}).Src("test/*.txt")

	p.Func("hello", nil, func(c *stencil.Context) {
		name := c.Args.AsString("default value", "name", "n")
		fmt.Println("Hello", name)
	}).Src("*.hello").Wait(3000)

	p.Func("server", nil, func(c *stencil.Context) {
		c.Start("main.go", stencil.M{"$in": "cmd/example"})
	}).Src("cmd/example/**/*.go")

	p.Func("change-package", nil, func(c *stencil.Context) {
		// works on mac
		c.Run(`find . -name "*.go" -print | xargs sed -i "" 's|gopkg.in/gostencil.v1|gopkg.in/gostencil.v2|g'`)
		// maybe linux?
		//Run(`find . -name "*.go" -print | xargs sed -i 's|gopkg.in/gostencil.v1|gopkg.in/gostencil.v2|g'`)
	})
}

func main() {
	stencil.Stencil(tasks)
}
