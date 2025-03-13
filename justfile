set shell := ["nu", "-c"]

run: compile
    @ ./bin/main

compile:
    @ rm -rf bin
    @ mkdir bin
    @ go build -o bin/main
