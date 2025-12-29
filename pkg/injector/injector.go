package injector

import "os/exec"

func Inject(cmd *exec.Cmd, secrets map[string]string) error {
	return cmd.Start()
}
