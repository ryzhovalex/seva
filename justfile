set shell := ["nu", "-c"]

run: compile
    @ ./Bin/Seva

compile:
    @ rm -rf Bin
    @ mkdir Bin
    @ go build -o Bin/Seva
