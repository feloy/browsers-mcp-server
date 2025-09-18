package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type Logger struct {
	logFile   string
	logLevel  string
	logWriter *os.File
}

func (o *Logger) initLoggerFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.logFile, "log-file", "", "Write Logs to file, no logs by default")
	_ = cmd.Flags().MarkHidden("log-file")
	cmd.Flags().StringVar(&o.logLevel, "log-level", "warn", "Log level to use (debug, info, warn, error), warn by default")
	_ = cmd.Flags().MarkHidden("log-level")
}

func (o *Logger) initLogger() {
	log.SetReportCaller(true)

	if o.logFile == "" {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.FatalLevel)
		return
	}
	logLevel, err := log.ParseLevel(o.logLevel)
	if err != nil {
		log.Error("invalid log level", "error", err)
		return
	}
	o.logWriter, err = os.Create(o.logFile)
	if err != nil {
		log.Error("failed to create log file", "error", err)
		return
	}
	log.SetOutput(o.logWriter)
	log.SetLevel(logLevel)
}

func (o *Logger) disposeLogger() {
	if o.logWriter != nil {
		_ = o.logWriter.Close()
	}
}
