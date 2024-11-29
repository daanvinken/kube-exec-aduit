package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"kube-exec-audit/pkg/kubeclient"
)

// Command flags
var (
	namespace   string
	pod         string
	container   string
	command     []string
	interactive bool
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command in a Kubernetes container",
	Long: `Executes a command in a Kubernetes container. You can specify the namespace, pod, and container,
as well as whether you want to start an interactive session.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("No command provided to execute. Use --help for usage details.")
		}

		command = args
		executeCommand()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Define flags and configuration settings
	execCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the pod (required)")
	execCmd.Flags().StringVarP(&pod, "pod", "p", "", "Name of the pod (required)")
	execCmd.Flags().StringVarP(&container, "container", "c", "", "Name of the container (optional)")
	execCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Start an interactive session (default false)")

	// Mark required flags
	_ = execCmd.MarkFlagRequired("namespace")
	_ = execCmd.MarkFlagRequired("pod")
}

func executeCommand() {
	fmt.Printf("Preparing to execute command in namespace: %s, pod: %s, container: %s\n", namespace, pod, container)

	// Initialize Kubernetes client
	clientset := kubeclient.NewClient()

	// Prepare exec options
	options := kubeclient.ExecOptions{
		Namespace:   namespace,
		Pod:         pod,
		Container:   container,
		Command:     command,
		Interactive: interactive,
	}

	// Perform the exec operation
	if err := kubeclient.Exec(clientset, options); err != nil {
		log.Fatalf("Exec operation failed: %v", err)
	}

	fmt.Println("Command executed successfully.")
}
