# sine
The open source centralized currency without value

[![Discord](https://img.shields.io/discord/1487827567695368450?logo=discord&color=%235865f2)](https://discord.gg/JZ6tGpeSeh)
![GitHub last commit](https://img.shields.io/github/last-commit/blacklizardcode/sine?logo=github)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/blacklizardcode/sine/docker-release.yaml?logo=github&label=release)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/blacklizardcode/sine/docker-nightly.yaml?logo=github&label=nightly)

---

## about sine

Sine is a passion project that I wanted to follow 3 rules
1. Fully open source
2. Only provide an official api (because I can't do frontend)
3. Make the currency have no conversion rate

This makes it possible for sine to be a local currency to use within your friend group, company or any group of people for rewards or things like that.

## Documentation

documentation is provided with github [Wiki](https://github.com/blacklizardcode/sine/wiki)

## installation

currently sine is in development and you will have to do installation on your own, I would recommend you do these steps
1. clone the repo
2. make a service file for the project with `go run .` or compile it first and add the binary in the service file
3. setup the database, the docker compose for a database is included, just run `docker compose up -d`

## Contributing

Information about contributing is in the [contributing.md](contributing.md) file

## License

This project is licensed under the gnu affero general public license, more information in the [LICENSE.md](/LICENSE.md) file