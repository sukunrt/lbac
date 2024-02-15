#!/bin/bash
input=$(cat)
go run . <<< "$input" > ../t.s
cd ..
clang main.c t.s
./a.out
