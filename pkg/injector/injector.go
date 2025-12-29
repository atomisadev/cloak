package injector

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func RunCommand(command []string, secrets map[string]string) error {
	if len(command) == 0 {
		return fmt.Errorf("injector: no command provided")
	}

	binary := command[0]
	args := command[1:]
	cmd := exec.Command(binary, args...)

	currentEnv := os.Environ()
	secretEnv := make([]string, 0, len(secrets))
	for k, v := range secrets {
		secretEnv = append(secretEnv, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = append(currentEnv, secretEnv...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("injector: failed to start command '%s': %w", binary, err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// goroutine to forward signals to child
	go func() {
		for sig := range sigChan {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
			}
		}
	}()

	err := cmd.Wait()

	signal.Stop(sigChan)
	close(sigChan)

	return err
}
