# Contributing

## Security fixes

Security fixes are always welcome, I'm not a security expert so it would be great if you put some of your knowledge into the project

## New features

You can always make new features, but new features aren't always what I had in mind for the project so always first make an issue with the label "feature greenlight" so you don't put too much effort into something that isn't wanted for the project

## Ai

Use ai responsibly, don't vibe code and if you create a pull request and I don't understand it you will have to rewrite it so I can, because if I can understand it, almost everybody can.

## How to contribute
With that out of the way I will now show you how to contribute

1. clone the repo, `git clone https://github.com/blacklizardcode/sine.git`
2. setup a postgresql database, you can just use the database from the docker-compose.yml as long as you comment out the sine service itself
3. set enviroment variables, sine uses enviroment veriables for authentication with the database, most of them can be seen in the docker-compose.yml
4. start the server, `go run .` or `go build . && ./sine`