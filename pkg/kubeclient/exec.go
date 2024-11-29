package kubeclient

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ExecOptions struct {
	Namespace   string
	Pod         string
	Container   string
	Command     []string
	Interactive bool
}

func Exec(clientset *kubernetes.Clientset, options ExecOptions) error {
	ctx := context.Background()

	// Build kubeconfig
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	// Generate a session ID for logging
	sessionID := uuid.New().String()

	// Prepare the exec request
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(options.Pod).
		Namespace(options.Namespace).
		SubResource("exec").
		Param("container", options.Container).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", fmt.Sprintf("%t", options.Interactive))

	for _, cmd := range options.Command {
		req.Param("command", cmd)
	}

	// Establish SPDY executor
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to initialize SPDY executor: %w", err)
	}

	// Create the logger
	logger, err := createLogger()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Log the start of the session
	sessionLog := LogEntry{
		SessionID: sessionID,
		User:      os.Getenv("USER"), // Replace with actual user detection
		Namespace: options.Namespace,
		Pod:       options.Pod,
		Container: options.Container,
		Command:   strings.Join(options.Command, " "),
		Timestamp: time.Now().Format(time.RFC3339),
	}
	logNewSession(logger, sessionLog)
	logCommand(logger, sessionLog)

	// Wrap stdout and stderr to log their output with correct metadata
	stdoutWriter := io.MultiWriter(os.Stdout, loggerWriter{logger: logger, entry: sessionLog, stream: "stdout"})
	stderrWriter := io.MultiWriter(os.Stderr, loggerWriter{logger: logger, entry: sessionLog, stream: "stderr"})

	// Stream the interactive session
	return exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  io.TeeReader(bufio.NewReader(os.Stdin), loggerWriter{logger: logger, entry: sessionLog, stream: "stdin"}),
		Stdout: stdoutWriter,
		Stderr: stderrWriter,
		Tty:    options.Interactive,
	})
}

// loggerWriter is a custom writer that logs input or output
type loggerWriter struct {
	logger    zerolog.Logger
	entry     LogEntry
	sessionID string
	stream    string // "stdin", "stdout", or "stderr"
}

func (lw loggerWriter) Write(p []byte) (int, error) {
	lw.logger.Info().
		Timestamp().
		Str("session_id", lw.entry.SessionID).
		Str("user", lw.entry.User).
		Str("namespace", lw.entry.Namespace).
		Str("pod", lw.entry.Pod).
		Str("container", lw.entry.Container).
		Str(lw.stream, string(p)). // Log the content of the stream
		Msg(fmt.Sprintf("%s stream activity", strings.ToUpper(lw.stream)))
	return len(p), nil
}
