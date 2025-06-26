package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr error
	}{
		{
			name: "valid inputs",
			cfg: Config{
				PipelinesRootPath: "/some/path",
				Username:          "testuser",
				Token:             "testtoken",
			},
			wantErr: nil,
		},
		{
			name: "missing username",
			cfg: Config{
				PipelinesRootPath: "/some/path",
				Token:             "testtoken",
			},
			wantErr: ErrMissingUsername,
		},
		{
			name: "missing token",
			cfg: Config{
				PipelinesRootPath: "/some/path",
				Username:          "testuser",
			},
			wantErr: ErrMissingToken,
		},
		{
			name: "missing both username and token",
			cfg: Config{
				PipelinesRootPath: "/some/path",
			},
			wantErr: ErrMissingUsername,
		},
		{
			name: "empty pipelines root path is invalid",
			cfg: Config{
				Username: "testuser",
				Token:    "testtoken",
			},
			wantErr: ErrMissingPipelinesRootPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
