package looker

import (
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

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

func TestFindExactMatch(t *testing.T) {
	strp := func(s string) *string { return &s }
	items := []*string{strp("axb"), nil, strp("A_B"), strp("other")}
	getName := func(s *string) *string { return s }

	tests := map[string]struct {
		name    string
		want    *string
		wantHit bool
	}{
		"exact match":                     {name: "A_B", want: items[2], wantHit: true},
		"case-insensitive match":          {name: "a_b", want: items[2], wantHit: true},
		"wildcard chars are not wildcard": {name: "a_x", wantHit: false},
		"not found":                       {name: "missing", wantHit: false},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			got, ok := findExactMatch(items, tt.name, getName)
			assert.Equal(t, tt.wantHit, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDataSourceSchemaFromResourceSchema(t *testing.T) {
	rs := map[string]*schema.Schema{
		"name":  {Type: schema.TypeString, Required: true, ForceNew: true},
		"flag":  {Type: schema.TypeBool, Optional: true, Default: false},
		"items": {Type: schema.TypeSet, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},
	}

	ds := dataSourceSchemaFromResourceSchema(rs, "name", "lookup key.")

	assert.True(t, ds["name"].Required)
	assert.False(t, ds["name"].Computed)
	assert.Equal(t, "lookup key. Matched case-insensitively.", ds["name"].Description)
	for _, k := range []string{"flag", "items"} {
		assert.True(t, ds[k].Computed, k)
		assert.False(t, ds[k].Optional, k)
		assert.False(t, ds[k].Required, k)
		assert.Nil(t, ds[k].Default, k)
	}
	assert.True(t, rs["name"].ForceNew)
	assert.True(t, rs["items"].Required)
}
