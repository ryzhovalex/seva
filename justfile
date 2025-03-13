set shell := ["nu", "-c"]

run: compile
    @ ./bin/main -shell

compile:
    @ rm -rf bin
    @ mkdir bin
    @ go build -o bin/main
