set shell := ["nu", "-c"]

run: compile
    @ ./Bin/Seva

compile: gentempl buildtailwind
    @ rm -rf Bin
    @ mkdir Bin
    @ go build -o Bin/Seva

gentempl:
    @ cd components
    @ templ generate

buildtailwind:
    @ npx tailwindcss -m -i Setup.css -o Static/Main.css
