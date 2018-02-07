package xexec

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	
	"strings"
	
	"bufio"
	"bytes"
	
	"github.com/pkg/errors"
)

type Options struct {
	dir  string
	args []string
}

type OptionFunc func(*Options)


// Run is an experimental run, does not work if it contains spaces inside command arguments or paths
func Run(command string, optFuncs ...OptionFunc) ([]byte, error) {
	options := &Options{}
	
	for _, optFunc := range optFuncs {
		optFunc(options)
	}
	
	if len(options.args) != 0 {
		command += " " + strings.Join(options.args, " ")
	}
	
	cmd := exec.Command("sh", "-c", command)
	
	if options.dir != "" {
		cmd.Dir = options.dir
	}
	
	log.Printf("Executing %s\n", command)
	
	return run(cmd)
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


func run(cmd *exec.Cmd) ([]byte, error) {
	// To fix a path issue caused by automatic path conversions in some windows shells, see https://github.com/docker/toolbox/issues/282
	cmd.Env = getEnv()
	
	stdout, err := cmd.StdoutPipe()
	
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	
	stderr, err := cmd.StderrPipe()
	
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	
	err = cmd.Start()
	
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	
	var b bytes.Buffer
	
	out := bufio.NewWriter(&b)
	
	go copyPipe(stdout, out, os.Stdout)
	go copyPipe(stderr, os.Stderr)
	
	err = cmd.Wait()
	
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	
	err = out.Flush()
	
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	
	if !cmd.ProcessState.Success() {
		//log.Printf("Failed to exec command\n")
		return nil, errors.New("failed to exec command")
	}
	
	//log.Printf("Success\n")
	return b.Bytes(), nil
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

func copyPipe(src io.ReadCloser, dsts ...io.Writer) {
	var err error
	
	buff := make([]byte, 100)
	
	var n int
	for err == nil {
		n, err = src.Read(buff)
		if n > 0 {
			for _, dst := range dsts {
				fmt.Fprintf(dst, "%s", string(buff[:n]))
			}
		}
	}
}
