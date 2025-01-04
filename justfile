set shell := ["nu", "-c"]

run: compile
    @ ./bin/server

compile: gentempl buildtailwind
    @ rm -rf bin
    @ mkdir bin
    @ go build -o bin/server

gentempl:
    @ cd components
    @ templ generate

buildtailwind:
    @ npx tailwindcss -m -i Setup.css -o Static/Main.css
