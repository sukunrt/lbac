package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestGetNumber(t *testing.T) {
	type result struct {
		i   int
		err error
	}
	cases := []struct {
		s        string
		expected []result
	}{
		{
			s: "10 20 30",
			expected: []result{
				{i: 10},
				{i: 20},
				{i: 30},
				{i: 0, err: io.EOF}},
		},
		{
			s: "10 20p30",
			expected: []result{
				{i: 10},
				{i: 20},
				{i: 0, err: errors.New("NAN")},
				{i: 30},
				{i: 0, err: io.EOF}},
		},
		{
			s: "10pp",
			expected: []result{
				{i: 10},
				{i: 0, err: errors.New("NAN")},
				{i: 0, err: errors.New("NAN")},
				{i: 0, err: io.EOF}},
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			s := newScanner(bufio.NewScanner(strings.NewReader(tc.s)))
			for _, res := range tc.expected {
				n, err := getNumber(s)
				s1, s2 := fmt.Sprintf("%s", err), fmt.Sprintf("%s", res.err)
				if n != res.i || s1 != s2 {
					t.Fatalf("expected: %d %s got: %d %s", res.i, res.err, n, err)
				}
				if res.err != nil {
					s.Pop()
				}
			}
		})
	}
}
