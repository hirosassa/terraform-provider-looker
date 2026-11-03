package looker

import (
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func TestHash(t *testing.T) {
	tests := map[string]struct {
		input   interface{}
		wantRes string
	}{
		"normal string": {
			input:   "hello",
			wantRes: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		"empty string": {
			input:   "",
			wantRes: "",
		},
		"nil": {
			input:   nil,
			wantRes: "",
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			actual := hash(tt.input)
			assert.Equal(t, tt.wantRes, actual)
		})
	}
}

func TestExpandStringListFromSet(t *testing.T) {
	tests := map[string]struct {
		input    []interface{}
		wantList []string
	}{
		"normal strings": {
			input:    []interface{}{"c", "a", "b"},
			wantList: []string{"a", "b", "c"},
		},
		"empty set": {
			input:    []interface{}{},
			wantList: nil,
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			set := schema.NewSet(schema.HashString, tt.input)
			actual := expandStringListFromSet(set)
			sort.Strings(actual)
			assert.Equal(t, tt.wantList, actual)
		})
	}
}

func TestFlattenStringListToSet(t *testing.T) {
	tests := map[string]struct {
		input    []string
		wantList []string
	}{
		"normal strings": {
			input:    []string{"c", "a", "b"},
			wantList: []string{"a", "b", "c"},
		},
		"empty slice": {
			input:    []string{},
			wantList: []string{},
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			result := flattenStringListToSet(tt.input)
			actual := make([]string, 0, result.Len())
			for _, v := range result.List() {
				actual = append(actual, v.(string))
			}
			sort.Strings(actual)
			assert.Equal(t, tt.wantList, actual)
		})
	}
}

func TestFlattenStringList(t *testing.T) {
	tests := map[string]struct {
		input   []string
		wantRes []interface{}
	}{
		"normal strings": {
			input:   []string{"a", "b", "c"},
			wantRes: []interface{}{"a", "b", "c"},
		},
		"empty slice": {
			input:   []string{},
			wantRes: []interface{}{},
		},
		"nil slice": {
			input:   nil,
			wantRes: []interface{}{},
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			actual := flattenStringList(tt.input)
			assert.Equal(t, tt.wantRes, actual)
		})
	}
}
