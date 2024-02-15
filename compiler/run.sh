#!/bin/bash
read -p "input: " input
go run . <<< $input > ../t.s
cd ..
clang main.c t.s
./a.out
