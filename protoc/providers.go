package protoc

import (
	"github.com/gofunct/pb/protoc/build"
	"sync"

	"github.com/google/wire"

	"github.com/izumin5210/gex"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

var (
	gexCfg   *gex.Config
	gexCfgMu sync.Mutex

	toolRepo   tool.Repository
	toolRepoMu sync.Mutex
)

func ProvideGexConfig(
	fs afero.Fs,
	execer exec.Interface,
	io build.IO,
	rootDir RootDir,
) *gex.Config {
	gexCfgMu.Lock()
	defer gexCfgMu.Unlock()
	if gexCfg == nil {
		gexCfg = &gex.Config{
			OutWriter:  io.Out(),
			ErrWriter:  io.Err(),
			InReader:   io.In(),
			FS:         fs,
			Execer:     execer,
			WorkingDir: rootDir.String(),
			Verbose:    build.IsVerbose() || build.IsDebug(),
			Logger:     zap.NewStdLog(zap.L()),
		}
	}
	return gexCfg
}

func ProvideToolRepository(gexCfg *gex.Config) (tool.Repository, error) {
	toolRepoMu.Lock()
	defer toolRepoMu.Unlock()
	if toolRepo == nil {
		var err error
		toolRepo, err = gexCfg.Create()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return toolRepo, nil
}

// WrapperSet is a provider set that includes gex things and Wrapper instance.
var WrapperSet = wire.NewSet(
	ProvideGexConfig,
	ProvideToolRepository,
	NewWrapper,
)
