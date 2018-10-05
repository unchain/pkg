package xexec

import (
	"fmt"
	"os"
	"os/exec"

	"strings"

	"bufio"
	"bytes"

	"io"

	"github.com/pkg/errors"
)

type Options struct {
	streamOutput bool
	dir          string
	args         []string
	fmtArgs      []interface{}
	host         string
}

type OptionFunc func(*Options)

// Run is an experimental run, does not work if it contains spaces inside command arguments or paths
func Run(command string, optFuncs ...OptionFunc) ([]byte, error) {
	options := &Options{}
	for _, optFunc := range optFuncs {
		optFunc(options)
	}

	if len(options.fmtArgs) != 0 {
		command = fmt.Sprintf(command, options.fmtArgs...)
	}

	if len(options.args) != 0 {
		command += " " + strings.Join(options.args, " ")
	}

	if options.host != "" {
		command = fmt.Sprintf(`
ssh %s <<'ENDSSH'
	%s	
ENDSSH`, options.host, command)
	}

	cmd := exec.Command("sh", "-c", command)

	if options.dir != "" {
		cmd.Dir = options.dir
	}

	//fmt.Printf("Executing %s\n", command)

	return run(cmd, options)
}

func WithFormat(fmtArgs ...interface{}) OptionFunc {
	return func(o *Options) {
		o.fmtArgs = fmtArgs
	}
}

func WithSSH(host string) OptionFunc {
	return func(o *Options) {
		o.host = host
	}
}

func WithDir(dir string) OptionFunc {
	return func(o *Options) {
		o.dir = dir
	}
}

func WithArgs(args ...string) OptionFunc {
	return func(o *Options) {
		o.args = args
	}
}

func StreamOutput() OptionFunc {
	return func(o *Options) {
		o.streamOutput = true
	}
}

func run(cmd *exec.Cmd, options *Options) ([]byte, error) {
	// To fix a path issue caused by automatic path conversions in some windows shells, see https://github.com/docker/toolbox/issues/282
	cmd.Env = getEnv()

	outPipe, err := cmd.StdoutPipe()

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	errPipe, err := cmd.StderrPipe()

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	err = cmd.Start()

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	outW := bufio.NewWriter(&outBuf)
	errW := bufio.NewWriter(&errBuf)

	if options.streamOutput {
		go io.Copy(io.MultiWriter(outW, os.Stdout), outPipe)
		go io.Copy(io.MultiWriter(errW, os.Stderr), errPipe)
	} else {
		go io.Copy(outW, outPipe)
		go io.Copy(errW, errPipe)
	}

	err = cmd.Wait()

	if err != nil {
		return nil, errors.Wrap(err, errBuf.String())
	}

	err = outW.Flush()

	if err != nil {
		return nil, errors.Wrap(err, errBuf.String())
	}

	if !cmd.ProcessState.Success() {
		//log.Printf("Failed to exec command\n")
		return nil, errors.Wrap(err, errBuf.String())
	}

	//log.Printf("Success\n")
	return outBuf.Bytes(), nil
}

func getEnv() []string {
	return append(os.Environ(),
		`MSYS_NO_PATHCONV=1`,
		//`VBOX_MSI_INSTALL_PATH=C:\Program Files\Oracle\VirtualBox\`,
		//`DOCKER_TLS_VERIFY=1`,
		//`DOCKER_HOST=tcp://192.168.99.100:2376`,
		//`DOCKER_MACHINE_NAME=default`,
		//`COMPOSE_CONVERT_WINDOWS_PATHS=false`,
		//`DOCKER_CERT_PATH=C:\Users\enikolov\.docker\machine\machines\default`,
		//`NO_PROXY=192.168.99.100`,
	)
}
