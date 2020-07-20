# Hivemind

<img align="right" width="224" height="74" title="Hivemind logo"
     src="https://cdn.rawgit.com/DarthSim/hivemind/master/logo.svg">

[![Build Status](https://travis-ci.org/DarthSim/hivemind.svg?branch=master)](https://travis-ci.org/DarthSim/hivemind)

Hivemind is a process manager for Procfile-based applications. At the moment, it supports Linux, FreeBSD, and macOS.

Procfile is a simple format to specify types of processes your application provides (such as web application server, background queue process, front-end builder) and commands to run those processes. It can significantly simplify process management for developers and is used by popular Platforms-as-a-Service, such as Heroku and Deis. You can learn more about the `Procfile` format [here](https://devcenter.heroku.com/articles/procfile).

There are some good Procfile-based process management tools, including [foreman](https://github.com/ddollar/foreman) by David Dollar, which started it all. The problem with most of those tools is that processes you want to manage start to think they are logging their output into a file, and that can lead to all sorts of problems: severe lagging, losing or breaking colored output. Tools can also add vanity information (unneeded timestamps in logs). Hivemind was created to fix those problems once and for all.

See this article for a good intro and all the juicy details! [Introducing
Overmind and Hivemind](https://evilmartians.com/chronicles/introducing-overmind-and-hivemind)

<a href="https://evilmartians.com/?utm_source=hivemind">
<img src="https://evilmartians.com/badges/sponsored-by-evil-martians.svg" alt="Sponsored by Evil Martians" width="236" height="54">
</a>

## Enter Hivemind

Hivemind uses `pty` to capture process output. That fixes any problem with log clipping, delays, and TTY colors other process management tools may have.

**If you would like a process management tool with a lot of features, including [tmux](https://tmux.github.io/) support, restarting and killing individual processes and advanced configuration, you should take a look at Hivemind's big brother — [Overmind](https://github.com/DarthSim/overmind)!**

## Installation

#### With Homebrew (macOS)

```bash
brew install hivemind
```

#### Download the latest Hivemind release binary

You can download the latest release [here](https://github.com/DarthSim/hivemind/releases/latest).

#### From Source

You need Go 1.11 or later to build the project.

```bash
$ GO111MODULE=on go get -u -f github.com/DarthSim/hivemind
```

**Note:** You can update Hivemind the same way.

## Usage

Hivemind works with a `Procfile`. It may look like this:

```Procfile
web: bin/rails server
worker: bundle exec sidekiq
assets: gulp watch
```

To get started, you just need to run Hivemind from your working directory containing `Procfile`.

```bash
$ hivemind
```

If `Procfile` isn't located in your working directory, or named it non-standard as `Procfile.dev`, you can specify the path to it: [**Fun Fact:** Name of the `Procfile` is arbitrary and can be anything, although it is a best practice to name it as `Procfile` for sanity]

```bash
$ hivemind path/to/your/Procfile
$ hivemind path/to/your/Procfile.dev
```

Run `hivemind --help` to see other options. Note that every Hivemind option can be set with corresponding environment variable.

### Environment

If you need to set specific environment variables before running a `Procfile`, you can specify them in the `.env` file in the current working directory. The file should contain `variable=value` pairs, one per line:

```
PATH=$PATH:/additional/path
PORT=3000
HIVEMIND_TITLE=my_awsome_app
```

## Author

Sergey "DarthSim" Aleksandrovich

Highly inspired by [Foreman](https://github.com/ddollar/foreman).

Many thanks to @antiflasher for the awesome logo.

## License

Hivemind is licensed under the MIT license.

See LICENSE for the full license text.
