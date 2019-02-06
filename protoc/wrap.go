package protoc

import (
	"context"
	"github.com/gofunct/pb/protoc/build"
	"os"
	"path/filepath"

	"github.com/izumin5210/gex/pkg/tool"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

// Wrapper can execute protoc commands for current project's proto files.
type Wrapper interface {
	Exec(context.Context) error
}

type wrapperImpl struct {
	cfg      *Config
	build    afero.Fs
	ui       UI
	execer   exec.Interface
	toolRepo tool.Repository
	rootDir  RootDir
}

// NewWrapper creates a new Wrapper instance.
func NewWrapper(cfg *Config, build afero.Fs, execer exec.Interface, ui UI, toolRepo tool.Repository, rootDir RootDir) Wrapper {
	return &wrapperImpl{
		cfg:      cfg,
		build:    build,
		ui:       ui,
		execer:   execer,
		toolRepo: toolRepo,
		rootDir:  rootDir,
	}
}

func (e *wrapperImpl) Exec(ctx context.Context) (err error) {
	e.ui.Section("Execute protoc")

	e.ui.Subsection("Install plugins")
	err = errors.WithStack(e.installPlugins(ctx))
	if err != nil {
		return
	}

	e.ui.Subsection("Execute protoc")
	err = errors.WithStack(e.execProtocAll(ctx))

	return
}

func (e *wrapperImpl) installPlugins(ctx context.Context) error {
	return errors.WithStack(e.toolRepo.BuildAll(ctx))
}

func (e *wrapperImpl) execProtocAll(ctx context.Context) error {
	protoFiles, err := e.cfg.ProtoFiles(e.build, e.rootDir.String())
	if err != nil {
		return errors.WithStack(err)
	}

	var errs []error
	for _, path := range protoFiles {
		err = e.execProtoc(ctx, path)
		relPath, _ := filepath.Rel(e.rootDir.String(), path)
		if err == nil {
			e.ui.ItemSuccess(relPath)
		} else {
			zap.L().Error("failed to execute protoc", zap.Error(err))
			errs = append(errs, err)
			e.ui.ItemFailure(relPath, err)
		}
	}

	if len(errs) > 0 {
		return errors.New("failed to execute protoc")
	}

	return nil
}

func (e *wrapperImpl) execProtoc(ctx context.Context, protoPath string) error {
	outDir, err := e.cfg.OutDirOf(e.rootDir.String(), protoPath)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = build.CreateDirIfNotExists(e.build, outDir); err != nil {
		return errors.WithStack(err)
	}

	cmds, err := e.cfg.Commands(e.rootDir.String(), protoPath)
	if err != nil {
		return errors.WithStack(err)
	}

	path := e.rootDir.BinDir().String() + string(filepath.ListSeparator) + os.Getenv("PATH")
	env := append(os.Environ(), "PATH="+path)

	for _, args := range cmds {
		cmd := e.execer.CommandContext(ctx, args[0], args[1:]...)
		cmd.SetEnv(env)
		cmd.SetDir(e.rootDir.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Wrapf(err, "failed to execute command: %v\n%s", args, string(out))
		}
	}

	return nil
}
