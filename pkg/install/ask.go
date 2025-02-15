package install

import (
	"os/exec"
	"strings"

	"github.com/rancher/os2/pkg/config"
	"github.com/rancher/os2/pkg/questions"
	"github.com/rancher/os2/pkg/util"
)

func Ask(cfg *config.Config) error {
	if err := AskInstallDevice(cfg); err != nil {
		return err
	}

	if err := AskConfigURL(cfg); err != nil {
		return err
	}

	if cfg.RancherOS.Install.ConfigURL == "" {
		if err := AskGithub(cfg); err != nil {
			return err
		}

		if err := AskPassword(cfg); err != nil {
			return err
		}
	}

	return nil
}

func AskInstallDevice(cfg *config.Config) error {
	if cfg.RancherOS.Install.Device != "" {
		return nil
	}

	output, err := exec.Command("/bin/sh", "-c", "lsblk -r -o NAME,TYPE | grep -w disk | grep -v fd0 | awk '{print $1}'").CombinedOutput()
	if err != nil {
		return err
	}
	fields := strings.Fields(string(output))
	i, err := questions.PromptFormattedOptions("Installation target. Device will be formatted", -1, fields...)
	if err != nil {
		return err
	}

	cfg.RancherOS.Install.Device = "/dev/" + fields[i]
	return nil
}

func isServer(cfg *config.Config) (bool, bool, error) {
	opts := []string{"server", "agent", "none"}
	i, err := questions.PromptFormattedOptions("Run as server or agent?", 0, opts...)
	if err != nil {
		return false, false, err
	}

	return i == 0, i == 1, nil
}

func AskPassword(cfg *config.Config) error {
	if cfg.RancherOS.Install.Silent || cfg.RancherOS.Install.Password != "" {
		return nil
	}

	var (
		ok   = false
		err  error
		pass string
	)

	for !ok {
		pass, ok, err = util.PromptPassword()
		if err != nil {
			return err
		}
	}

	if pass != "" {
		pass, err = util.GetEncryptedPasswd(pass)
		if err != nil {
			return err
		}
	}

	cfg.RancherOS.Install.Password = pass
	return nil
}

func AskGithub(cfg *config.Config) error {
	if len(cfg.SSHAuthorizedKeys) > 0 || cfg.RancherOS.Install.Password != "" || cfg.RancherOS.Install.Silent {
		return nil
	}

	ok, err := questions.PromptBool("Authorize GitHub users to root SSH?", false)
	if !ok || err != nil {
		return err
	}

	str, err := questions.Prompt("Comma separated list of GitHub users to authorize: ", "")
	if err != nil {
		return err
	}

	for _, s := range strings.Split(str, ",") {
		cfg.SSHAuthorizedKeys = append(cfg.SSHAuthorizedKeys, "github:"+strings.TrimSpace(s))
	}

	return nil
}

func AskConfigURL(cfg *config.Config) error {
	if cfg.RancherOS.Install.ConfigURL != "" || cfg.RancherOS.Install.Silent {
		return nil
	}

	ok, err := questions.PromptBool("Configure system using an cloud-config file?", false)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	str, err := questions.Prompt("cloud-config file location (file path or http URL): ", "")
	if err != nil {
		return err
	}

	cfg.RancherOS.Install.ConfigURL = str
	return nil
}
