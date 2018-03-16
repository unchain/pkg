package xconfig

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"io"

	multierror "github.com/hashicorp/go-multierror"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/unchainio/pkg/xerrors"
)

func LoadConfig(cfg interface{}, path ...string) error {
	return Load(cfg, FromPaths(path...))
}

type Options struct {
	viper *viper.Viper
}

type remoteConfig struct {
	provider, endpoint, path string
}

type OptionFunc func(*Options) error

func Load(cfg interface{}, optFuncs ...OptionFunc) error {
	var err error
	var warns error

	opts := &Options{
		viper: viper.New(),
	}

	for _, optFunc := range optFuncs {
		warns = multierror.Append(warns, optFunc(opts))
	}

	err = opts.viper.Unmarshal(cfg)

	if err != nil {
		return multierror.Append(warns, errors.Wrapf(err, "failed to unmarshal config files"))
	}

	if warns != nil {
		return errors.Wrapf(xerrors.ToWarn(warns), "failed to load configs from some sources")
	}

	return nil
}

func FromPaths(paths ...string) OptionFunc {
	return func(o *Options) error {
		return MergeInConfigs(o.viper, paths)
	}
}

func FromReaders(cfgType string, readers ...io.Reader) OptionFunc {
	return func(o *Options) error {
		return MergeInReaders(o.viper, cfgType, readers)
	}
}

func FromRemote(provider, endpoint, path string) OptionFunc {
	return func(o *Options) error {
		err := o.viper.AddRemoteProvider(provider, endpoint, path)

		if err != nil {
			return errors.Wrapf(err, "failed to connect to remote config provider")
		}

		err = o.viper.ReadRemoteConfig()

		if err != nil {
			return errors.Wrapf(err, "failed to load config from remote provider")
		}

		return nil
	}
}

// MergeInConfigs merges the viper configs found in several readers into a single one
func MergeInReaders(v *viper.Viper, cfgType string, readers []io.Reader) error {
	var errs error

	v.SetConfigType(cfgType)

	for i, reader := range readers {
		err := v.MergeConfig(reader)

		if err != nil {
			msg := fmt.Sprintf("failed to load config from reader with index: `%d`", i)
			errs = multierror.Append(errs, errors.WithMessage(err, msg))
		}
	}

	if errs != nil {
		return errors.Wrap(xerrors.ToWarn(errs), "failed to load configs from some readers")
	}

	return nil
}

// MergeInConfigs merges the viper configs found in several paths into a single one
func MergeInConfigs(v *viper.Viper, paths []string) error {
	var errs error

	for _, path := range paths {
		v.SetConfigFile(absPathify(path))
		err := v.MergeInConfig()

		if err != nil {
			msg := fmt.Sprintf("failed to load config from path: `%s`", path)
			errs = multierror.Append(errs, errors.WithMessage(err, msg))
		}
	}

	if errs != nil {
		return errors.Wrap(xerrors.ToWarn(errs), "failed to load configs from some paths")
	}

	return nil
}

func absPathify(inPath string) string {
	home, err := homedir.Dir()

	if err != nil {
		return ""
	}

	if strings.HasPrefix(inPath, "$HOME") {
		inPath = home + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)

	if err != nil {
		return ""
	}

	return filepath.Clean(p)
}
