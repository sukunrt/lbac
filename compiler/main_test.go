package main

import (
	"fmt"
	"maps"
	"testing"
)

func TestNewScope(t *testing.T) {
	makeArgs := func(args ...string) []token {
		targs := make([]token, len(args))
		for i, a := range args {
			targs[i] = token{T: Identifier, V: a}
		}
		return targs
	}

	cases := []struct {
		args   []string
		result map[string]int
	}{
		{
			args: []string{"a", "b", "c"},
			result: map[string]int{
				"c": 2,
				"b": 3,
				"a": 4,
			},
		},
		{
			args: []string{"a"},
			result: map[string]int{
				"a": 2,
			},
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			scope := newScope(makeArgs(tc.args...))
			if !maps.Equal(scope, tc.result) {
				t.Fatalf("expected maps to be equal\n%+v\n%+v", scope, tc.result)
			}
		})
	}
}
