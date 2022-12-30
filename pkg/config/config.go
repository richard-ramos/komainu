package config

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	// DataDir is the file system folder the node should use for any data storage needs.
	DataDir string `validate:"required"`

	// NodeKey is an optional hex-encoded node ID (private key). It will be used for both
	// remote peer identification as well as network traffic encryption.
	NodeKey string

	// LogFile is a folder which contains LogFile
	LogDir string

	// LogFile is filename where exposed logs get written to
	LogFile string

	// LogEnabled enables the logger
	LogEnabled bool `json:"LogEnabled"`

	// LogLevel defines minimum log level. Valid names are "ERROR", "WARN", "INFO", "DEBUG", and "TRACE".
	LogLevel string `validate:"eq=ERROR|eq=WARN|eq=INFO|eq=DEBUG|eq=TRACE"`

	// LogToStderr defines whether logged info should also be output to os.Stderr
	LogToStderr bool
}

func (c *Config) PathDataDir() string {
	dir := os.Getenv("KUMAINU_DIR")
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(c.DataDir)
		if err != nil {
			return c.DataDir
		}
		return dir
	}
	return dir
}
