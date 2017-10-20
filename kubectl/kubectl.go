package kubectl

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"k8s.io/client-go/rest"
)

type KubeCtl struct {
	config    *rest.Config
	namespace string
	debug     bool
}

func NewKubeCtl(config *rest.Config, namespace string, debug bool) *KubeCtl {
	return &KubeCtl{
		config:    config,
		namespace: namespace,
		debug:     debug,
	}
}

func (t *KubeCtl) Run(stdin []byte, args ...string) (string, error) {
	args = append(t.configArgs(), args...)
	if t.debug {
		log.Printf("-----KUBECTL DEBUG-----\nkubectl %s\n%s\n", strings.Join(args, " "), string(stdin))
	}

	cmd := exec.Command("kubectl", args...)

	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}

	out, err := cmd.Output()
	if err != nil {
		errmsg := err.Error()
		exiterr, ok := err.(*exec.ExitError)
		if ok {
			errmsg = fmt.Sprintf("%s: %s", errmsg, string(exiterr.Stderr))
		}
		return "", fmt.Errorf("Kubectl %v failed: %s, %s", args, errmsg, out)
	}

	sout := string(out)
	if t.debug {
		log.Printf("-----KUBECTL OUTPUT-----\n%s\n", sout)
	}
	return string(out), nil
}

func (t *KubeCtl) configArgs() []string {
	args := []string{
		"--namespace", t.namespace,
	}

	cfg := t.config
	if cfg.Host != "" {
		args = append(args, "--server", cfg.Host)
	}
	if cfg.CAFile != "" {
		args = append(args, "--certificate-authority", cfg.CAFile)
	}
	if cfg.CertFile != "" {
		args = append(args, "--client-certificate", cfg.CertFile)
	}
	if cfg.CertFile != "" {
		args = append(args, "--client-key", cfg.KeyFile)
	}
	if cfg.BearerToken != "" {
		args = append(args, "--token", cfg.BearerToken)
	}

	return args
}
