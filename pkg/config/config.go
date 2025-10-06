package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// Validation errors
var (
	ErrMissingUsername          = errors.New("username must be provided")
	ErrMissingToken             = errors.New("token must be provided")
	ErrMissingPipelinesRootPath = errors.New("pipelines root path must be provided")
	ErrMissingNamespace         = errors.New("namespace must be provided")
)

// Default timeout for the action
const DefaultTimeout = 1 * time.Minute

// Config is the inputs used to configure the action
type Config struct {
	// FleetManagementURL is the URL of the Fleet Management API
	FleetManagementURL string

	// PipelinesRootPath is the path to start recursing from when looking for pipelines
	PipelinesRootPath string

	// Username is the username when authenticating to Fleet Management
	Username string

	// Token is the token to use when authenticating to Fleet Management
	Token string

	// Namespace is the namespace to sync the pipelines to
	Namespace string

	// GlobalMatcher is an optional matcher to add to all pipelines
	GlobalMatcher string

	// Timeout is the maximum time to wait for operations to complete
	Timeout time.Duration

	// Verbose enables detailed logging output
	Verbose bool

	// DryRun runs without making actual changes
	DryRun bool
}

// NewFromEnv creates a new Config from GitHub Action environment variables
func NewFromEnv() (*Config, error) {
	cfg := &Config{
		FleetManagementURL: os.Getenv("INPUT_FM_URL"),
		PipelinesRootPath:  os.Getenv("INPUT_PIPELINES_ROOT_PATH"),
		Username:           os.Getenv("INPUT_FM_USERNAME"),
		Token:              os.Getenv("INPUT_FM_TOKEN"),
		Namespace:          os.Getenv("INPUT_NAMESPACE"),
		GlobalMatcher:      os.Getenv("INPUT_GLOBAL_MATCHER"),
		Timeout:            DefaultTimeout,
	}

	// Parse timeout if provided
	if timeoutStr := os.Getenv("INPUT_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout value '%s': %w", timeoutStr, err)
		}
		cfg.Timeout = timeout
	}

	// Parse verbose flag
	if verboseStr := os.Getenv("INPUT_VERBOSE"); verboseStr != "" {
		cfg.Verbose = verboseStr == "true"
	}

	// Parse dry-run flag
	if dryRunStr := os.Getenv("INPUT_DRY_RUN"); dryRunStr != "" {
		cfg.DryRun = dryRunStr == "true"
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the provided config
func (c *Config) Validate() error {
	if c.Username == "" {
		return ErrMissingUsername
	}
	if c.Token == "" {
		return ErrMissingToken
	}
	if c.PipelinesRootPath == "" {
		return ErrMissingPipelinesRootPath
	}
	if c.Namespace == "" {
		return ErrMissingNamespace
	}
	return nil
}

// SetupLogging configures the global logger based on the verbose flag
func (c *Config) SetupLogging() {
	logLevel := slog.LevelInfo
	if c.Verbose {
		logLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}
