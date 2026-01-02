package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Run initializes the container namespaces and re-execs the binary.
func Run() {
	startTime := time.Now()

	if len(os.Args) < 3 {
		fmt.Println("Usage: nanok run <command> [args]")
		os.Exit(1)
	}

	fmt.Printf("Parent: Running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child", fmt.Sprintf("%d", startTime.UnixNano())}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running child process: %v\n", err)
		os.Exit(1)
	}
}

// Child is the entry point for the re-executed binary inside the namespaces.
func Child() {
	if len(os.Args) < 4 {
		fmt.Println("Invalid child arguments")
		os.Exit(1)
	}

	startTimeNano := int64(0)
	fmt.Sscanf(os.Args[2], "%d", &startTimeNano)
	startTime := time.Unix(0, startTimeNano)

	command := os.Args[3]
	args := os.Args[4:]

	fmt.Printf("Child: Running %v as PID %d\n", os.Args[3:], os.Getpid())

	// Phase 1 doesn't include pivot_root yet, but we'll need it in Phase 2.
	// For now, we just log the startup time and exec the command.

	elapsed := time.Since(startTime)
	fmt.Printf("Startup time: %vms\n", elapsed.Milliseconds())

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
