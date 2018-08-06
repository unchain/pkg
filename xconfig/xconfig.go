package xconfig

import (
	"fmt"
	"io"

	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/unchainio/pkg/xlogger"
	"github.com/unchainio/pkg/xpath"
	"github.com/fsnotify/fsnotify"
	"os"
)

func Load(cfg interface{}, optFuncs ...OptionFunc) error {
	var err error
	opts := &Options{}

	for _, optFunc := range optFuncs {
		err = optFunc(opts)

		if err != nil {
			return err
		}
	}

	if opts.viper == nil {
		opts.viper = viper.New()
	}

	if opts.pathFlag != nil {
		cfgPath := ""
		flag.StringVar(&cfgPath, opts.pathFlag.Name, opts.pathFlag.DefValue, opts.pathFlag.Usage)
		flag.Parse()

		if cfgPath != "" {
			opts.paths = []string{cfgPath}
		}
	}

	var log *xlogger.Logger

	if opts.verbose {
		log, err = xlogger.New(&xlogger.Config{
			Level:  "info",
			Format: "text",
		})
	} else {
		log, err = xlogger.New(&xlogger.Config{
			Level: "panic",
		})
	}

	if opts.watch {
		if len(opts.paths) < 1 {
			return errors.New("can only watch config changes when path flag is set")
		}
		go func() {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				panic(err)
			}
			defer watcher.Close()

			dir, err := os.Getwd()

			for _, path := range opts.paths {
				done := make(chan bool)

				go func() {
					for {
						select {
						case event := <-watcher.Events:
							if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
								panic("CONFIG UPDATED - FORCING RESTART")
							}
						}
					}
				}()
				err = watcher.Add(fmt.Sprintf(dir + "/" + path))
				if err != nil {
					panic(err)
				}
				<- done
			}

		}()
	}

	log.Printf("Loading config from %+v", opts.paths)

	if len(opts.paths) != 0 {
		err = MergeInConfigs(opts.viper, opts.paths, log)

		if err != nil {
			return err
		}
	}

	if opts.cfgType != "" && len(opts.readers) != 0 {
		err = MergeInReaders(opts.viper, opts.cfgType, opts.readers, log)

		if err != nil {
			return err
		}
	}

	if opts.remote != nil {
		err = opts.viper.AddRemoteProvider(opts.remote.provider, opts.remote.endpoint, opts.remote.path)

		if err != nil {
			return errors.Wrapf(err, "failed to connect to remote config provider")
		}

		err = opts.viper.ReadRemoteConfig()

		if err != nil {
			return errors.Wrapf(err, "failed to load config from remote provider")
		}
	}

	if opts.automaticEnv {
		opts.viper.AutomaticEnv()
	}

	err = opts.viper.Unmarshal(cfg)

	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal config files")
	}

	return nil
}

type OptionFunc func(*Options) error

type Options struct {
	verbose bool
	watch bool

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

func Verbose(flag bool) OptionFunc {
	return func(o *Options) error {
		o.verbose = flag
		return nil
	}
}

func WithWatcher() OptionFunc {
	return func(o *Options) error {
		o.watch = true

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

		usage := fmt.Sprintf("Path to the config file, defaults to %s.", defValue)

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
func MergeInConfigs(v *viper.Viper, paths []string, log *xlogger.Logger) error {
	for _, path := range paths {
		v.SetConfigFile(xpath.Abs(path))
		err := v.MergeInConfig()

		if err != nil {
			log.Warnf("failed to load config from path: `%s`, error was: %+v", path, err)
		}
	}

	return nil
}

// MergeInConfigs merges the viper configs found in several readers into a single one
func MergeInReaders(v *viper.Viper, cfgType string, readers []io.Reader, log *xlogger.Logger) error {
	v.SetConfigType(cfgType)

	for i, reader := range readers {
		err := v.MergeConfig(reader)

		if err != nil {
			log.Warnf("failed to load config from reader with index: `%d`, error was: %+v", i, err)
		}
	}

	return nil
}
