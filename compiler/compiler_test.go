package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

type testCase struct {
	input  string
	output int
}

func TestCorrectness(t *testing.T) {
	cases := []testCase{
		{
			input:  "2",
			output: 2,
		},
		{
			input:  "1+1",
			output: 2,
		},
		{
			input:  "(1+1)*(1*10)",
			output: 20,
		},
		{
			input:  "(1+0+1)/(-2)",
			output: -1,
		},
		{
			input:  "1\n   1   \n(2+3)",
			output: 5,
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			out := strings.Builder{}
			cmd := exec.Cmd{
				Path:   "./run.sh",
				Stdin:  strings.NewReader(tc.input),
				Stdout: &out,
			}
			cmd.Run()
			res := out.String()
			n, err := strconv.Atoi(res[:len(res)-1])
			if err != nil {
				t.Fatal(err)
			}
			if n != tc.output {
				t.Fatal("invalid output", tc.output, n, res)
			}
		})
	}
}
