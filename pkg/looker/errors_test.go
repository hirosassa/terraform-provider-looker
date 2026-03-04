package looker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapSDKError(t *testing.T) {
	tests := map[string]struct {
		err          error
		operation    string
		resourceType string
		format       string
		args         []interface{}
		wantErr      bool
		wantMsg      string
	}{
		"nil error returns nil": {
			err:          nil,
			operation:    "CreateUser",
			resourceType: "user",
			format:       "%s",
			args:         []interface{}{"test@example.com"},
			wantErr:      false,
		},
		"no format includes operation and resource type": {
			err:          errors.New("api error"),
			operation:    "AllUsers",
			resourceType: "user",
			format:       "",
			wantErr:      true,
			wantMsg:      `AllUsers failed for user: api error`,
		},
		"single arg format includes quoted identifier": {
			err:          errors.New("not found"),
			operation:    "CreateConnection",
			resourceType: "connection",
			format:       "%s",
			args:         []interface{}{"my-db"},
			wantErr:      true,
			wantMsg:      `CreateConnection failed for connection "my-db": not found`,
		},
		"multiple args format includes composite identifier": {
			err:          errors.New("not found"),
			operation:    "UpdateUser",
			resourceType: "user",
			format:       "email=%s, id=%s",
			args:         []interface{}{"a@example.com", "42"},
			wantErr:      true,
			wantMsg:      `UpdateUser failed for user "email=a@example.com, id=42": not found`,
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			actual := wrapSDKError(tt.err, tt.operation, tt.resourceType, tt.format, tt.args...)
			if tt.wantErr {
				assert.EqualError(t, actual, tt.wantMsg)
			} else {
				assert.NoError(t, actual)
			}
		})
	}
}
