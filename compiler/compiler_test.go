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
		{
			input: `
			x=5
			y=15
			z=2
			(x+y-10)/z + 1
			`,
			output: 6,
		},
		{
			input: `
			x=5
			y=15
			z=2
			x=15
			(x+y-10)/z + 1
			`,
			output: 11,
		},
		{
			input: `
			x=5
			y=15
			IF x-y
				y=2+2
			ELSE
				y=2+3
			ENDIF
			y
			`,
			output: 4,
		},
		{
			input: `
			x=5
			y=5
			IF x-y
				y=2+2
			ELSE
				y=2+3
			ENDIF
			y
			`,
			output: 5,
		},
		{
			input: `
			x=100
			y=5
			WHILE x-y
				x = x-1
			ENDWHILE
			x
			`,
			output: 5,
		},
		{
			input: `
			x=100
			y=5
			WHILE x > y
				x = x-1
			ENDWHILE
			x
			`,
			output: 5,
		},
		{
			input: `
			x=1024
			y=2
			WHILE x > y
				y = y * 2
			ENDWHILE
			y
			`,
			output: 1024,
		},
		{
			input: `
			x=100
			y=1
			z=10
			FN hello(x y z)
				x+y+z
			ENDFN
			CALL hello(x y z)
			`,
			output: 111,
		},
		{
			input: `
			x=10
			y=1
			FN sum(st end)
				z = 0
				WHILE st < end
					z = z + st
					st = st + 1
				ENDWHILE
				z
			ENDFN
			CALL sum(y x)
			`,
			output: 45,
		},
		{
			input: `
			FN fib(i)
				f = 0
				IF i < 2
					f=i
				ELSE
					x = i-1
					y = i-2
					a = CALL fib(x)
					b = CALL fib(y)
					f = a+b
				ENDIF
				f
			ENDFN
			z = CALL fib(6)
			z
			`,
			output: 8,
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
