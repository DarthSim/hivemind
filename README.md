# Hivemind

Hivemind is a tool for running processes of development environment. At the moment it supports Linux, FreeBSD and Mac OS X.

<a href="https://evilmartians.com/?utm_source=hivemind">
<img src="https://evilmartians.com/badges/sponsored-by-evil-martians.svg" alt="Sponsored by Evil Martians" width="236" height="54">
</a>

## Installation

You need Go to build the project.

### Using `go get`

```bash
$ go get -u github.com/DarthSim/hivemind
```

### Using `make`

```bash
git clone https://github.com/DarthSim/hivemind.git
cd hivemind
make install
```

## Usage

Hivemind works with a Procfile.

```Procfile
web: bin/rails server
worker: bundle exec sidekiq
assets: gulp watch
```

To get started you just need to run Hivemind from your working directory containing Procfile.

```bash
$ hivemind
```

If Procfile isn't located in your working directory, you can specify it:

```bash
$ hivemind path/to/your/Procfile
```

Run `hivemind --help` to see other options.

## Author

Sergey "DarthSim" Aleksandrovich

## License

Hivemind is licensed under the MIT license.

See LICENSE for the full license text.
