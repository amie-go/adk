package xruntime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncNameDetailed(t *testing.T) {
	var tests = []struct {
		name     string
		expected []string
	}{
		{"", []string{"", "", ""}},
		{"/", []string{"", "", ""}},
		{".", []string{"", "", ""}},
		{"/.", []string{"", "", ""}},
		{"testing.tRunner", []string{"", "testing", "tRunner"}},
		{"github.com/amie-go/adk/errgen.TestStacktrace.func2.1", []string{"github.com/amie-go/adk", "errgen", "TestStacktrace.func2.1"}},
	}

	for _, v := range tests {
		moduleName, pkgName, funcName := FuncNameDetailed(v.name)
		assert.Equal(t, v.expected[0], moduleName)
		assert.Equal(t, v.expected[1], pkgName)
		assert.Equal(t, v.expected[2], funcName)
	}
}

func TestPackageInfo(t *testing.T) {
	assert.Equal(t, "github.com/amie-go/adk", moduleName)
	assert.Equal(t, "common", packageName)
}
