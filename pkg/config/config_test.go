package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg  Config
		wantErr error
	}{
		{
			name: "valid inputs",
			cfg: Config{
				RootPath: "/some/path",
				Username: "testuser",
				Token:    "testtoken",
			},
			wantErr: nil,
		},
		{
			name: "missing username",
			cfg: Config{
				RootPath: "/some/path",
				Token:    "testtoken",
			},
			wantErr: ErrMissingUsername,
		},
		{
			name: "missing token",
			cfg: Config{
				RootPath: "/some/path",
				Username: "testuser",
			},
			wantErr: ErrMissingToken,
		},
		{
			name: "missing both username and token",
			cfg: Config{
				RootPath: "/some/path",
			},
			wantErr: ErrMissingUsername,
		},
		{
			name: "empty root path is valid",
			cfg: Config{
				Username: "testuser",
				Token:    "testtoken",
			},
			wantErr: nil,
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
