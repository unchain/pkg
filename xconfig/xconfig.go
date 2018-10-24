package xconfig

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/unchainio/pkg/xpath"
)

type Info struct {
	Paths []string
}

func Load(cfg interface{}, optFuncs ...OptionFunc) error {
	var err error
	opts := &Options{}

	for _, optFunc := range optFuncs {
		err = optFunc(opts)

		if err != nil {
			return err
		}
	}

	if opts.info == nil {
		opts.info = new(Info)
	}

	if opts.viper == nil {
		opts.viper = viper.New()
	}

	if opts.pathFlag != nil {
		cfgPath := ""
		//fs := flag.NewFlagSet("", flag.ExitOnError)
		fs := flag.CommandLine
		fs.StringVar(&cfgPath, opts.pathFlag.Name, opts.pathFlag.DefValue, opts.pathFlag.Usage)
		FlagsFromEnv(fs)
		fs.Parse(os.Args[1:])

		if cfgPath != "" {
			opts.paths = []string{cfgPath}
		}
	}

	if opts.watchFn != nil {
		opts.viper.OnConfigChange(opts.watchFn)
	}

	var errs error

	if len(opts.paths) != 0 {
		opts.info.Paths = opts.paths

		err = MergeInConfigs(opts.viper, opts.paths)

		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	if opts.cfgType != "" && len(opts.readers) != 0 {
		err = MergeInReaders(opts.viper, opts.cfgType, opts.readers)

		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	if opts.remote != nil {
		err = opts.viper.AddRemoteProvider(opts.remote.provider, opts.remote.endpoint, opts.remote.path)

		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "failed to connect to remote config provider"))
		}

		err = opts.viper.ReadRemoteConfig()

		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "failed to load config from remote provider"))
		}
	}

	if opts.automaticEnv {
		opts.viper.AutomaticEnv()
	}

	err = opts.viper.Unmarshal(cfg)

	if err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "failed to unmarshal config files"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

type OptionFunc func(*Options) error

type Options struct {
	info     *Info
	watchFn  func(in fsnotify.Event)
	pathFlag *flag.Flag
	paths    []string

	cfgType string
	readers []io.Reader

	remote *remoteConfig

	automaticEnv bool

	viper *viper.Viper
}

type remoteConfig struct {
	provider, endpoint, path string
}

func GetInfo(info *Info) OptionFunc {
	return func(o *Options) error {
		o.info = info
		return nil
	}
}

func WithWatcher(watchFn func(in fsnotify.Event)) OptionFunc {
	return func(o *Options) error {
		o.watchFn = watchFn

		return nil
	}
}

func WithViper(v *viper.Viper) OptionFunc {
	return func(o *Options) error {
		o.viper = v
		return nil
	}
}

func FromPathFlag(name string, defValue string) OptionFunc {
	return func(o *Options) error {
		if len(o.paths) != 0 {
			return errors.New("FromPathFlag and FromPaths are incompatible")
		}

		usage := fmt.Sprintf("Path to the config file")

		o.pathFlag = &flag.Flag{
			Name:     name,
			DefValue: defValue,
			Usage:    usage,
		}

		return nil
	}
}

func FromPaths(paths ...string) OptionFunc {
	return func(o *Options) error {
		if o.pathFlag != nil {
			return errors.New("FromPathFlag and FromPaths are incompatible")
		}

		o.paths = paths

		return nil
	}
}

func FromReaders(cfgType string, readers ...io.Reader) OptionFunc {
	return func(o *Options) error {
		o.cfgType = cfgType
		o.readers = readers

		return nil
	}
}

func FromRemote(provider, endpoint, path string) OptionFunc {
	return func(o *Options) error {
		o.remote = &remoteConfig{
			provider: provider,
			endpoint: endpoint,
			path:     path,
		}

		return nil
	}
}

func FromEnv() OptionFunc {
	return func(o *Options) error {
		o.automaticEnv = true

		return nil
	}
}

// MergeInConfigs merges the viper configs found in several paths into a single one
func MergeInConfigs(v *viper.Viper, paths []string) error {
	var errs error

	for _, path := range paths {
		v.SetConfigFile(xpath.Abs(path))
		err := v.MergeInConfig()

		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "failed to load config from path: `%s`", path))
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

// MergeInConfigs merges the viper configs found in several readers into a single one
func MergeInReaders(v *viper.Viper, cfgType string, readers []io.Reader) error {
	var errs error
	v.SetConfigType(cfgType)

	for i, reader := range readers {
		err := v.MergeConfig(reader)

		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "failed to load config from reader with index: `%d``", i))
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

func FlagsFromEnv(fs *flag.FlagSet) {
	flags := map[string]struct{}{}
	fs.Visit(func(f *flag.Flag) {
		flags[f.Name] = struct{}{}
	})

	fs.VisitAll(func(f *flag.Flag) {
		envVar := strings.ToUpper(f.Name)
		envVar = strings.Replace(envVar, "-", "_", -1)
		val := os.Getenv(envVar)

		if val != "" {
			if _, ok := flags[f.Name]; !ok {
				fs.Set(f.Name, val)
			}
		}

		f.Usage = fmt.Sprintf("%s [%s]", f.Usage, envVar)
	})
}
