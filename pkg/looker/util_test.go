package looker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifferentStringSlices(t *testing.T) {
	tests := map[string]struct {
		a       []string
		b       []string
		wantRes bool
	}{
		"same slice": {
			a:       []string{"a", "b", "c"},
			b:       []string{"a", "b", "c"},
			wantRes: false,
		},
		"different slice": {
			a:       []string{"a", "b", "c"},
			b:       []string{"a", "b", "d"},
			wantRes: true,
		},
		"one slice is empty": {
			a:       []string{},
			b:       []string{"a", "b", "c"},
			wantRes: true,
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			actual := differentStringSlices(tt.a, tt.b)
			assert.Equal(t, tt.wantRes, actual)
		})
	}
}

func TestBuildTwoPartID(t *testing.T) {
	tests := map[string]struct {
		a       string
		b       string
		wantRes string
	}{
		"normal string": {
			a:       "abc",
			b:       "def",
			wantRes: "abc:def",
		},
		"both empty string": {
			a:       "",
			b:       "",
			wantRes: ":",
		},
		"first string is empty": {
			a:       "",
			b:       "def",
			wantRes: ":def",
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			actual := buildTwoPartID(&tt.a, &tt.b)
			assert.Equal(t, tt.wantRes, actual)
		})
	}
}

func TestParseTwoPartID(t *testing.T) {
	tests := map[string]struct {
		id       string
		wantRes1 string
		wantRes2 string
		wantErr  bool
	}{
		"normal input": {
			id:       "123:456",
			wantRes1: "123",
			wantRes2: "456",
			wantErr:  false,
		},
		"no colon contained": {
			id:       "123456",
			wantRes1: "",
			wantRes2: "",
			wantErr:  true,
		},
		"first part only": {
			id:       "123:",
			wantRes1: "123",
			wantRes2: "",
			wantErr:  false,
		},
		"second part only": {
			id:       ":456",
			wantRes1: "",
			wantRes2: "456",
			wantErr:  false,
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			a := assert.New(t)
			actualRes1, actualRes2, actualErr := parseTwoPartID(tt.id)
			if tt.wantErr {
				a.Error(actualErr)
			} else {
				a.NoError(actualErr)
				a.Equal(tt.wantRes1, actualRes1)
				a.Equal(tt.wantRes2, actualRes2)
			}
		})
	}
}
