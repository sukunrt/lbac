#!/bin/bash
read -p "input: " input
cd compiler
go run . <<< $input > ../t.s
cd ..
clang main.c t.s
./a.out
