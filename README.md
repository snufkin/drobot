[ ![Codeship Status for snufkin/drobot](https://app.codeship.com/projects/a3f473f0-3fcd-0136-942a-6a4ec1b9b8d8/status?branch=master)](https://app.codeship.com/projects/290926)
[![Go Report Card](https://goreportcard.com/badge/github.com/snufkin/drobot)](https://goreportcard.com/report/github.com/snufkin/drobot)

# Overview
The purpose of this utility is to check what components are outdated in a Drupal
 installation. It requires a manifest file, which contains the list of modules
 and themes along with the version information, and compares it to the release
 information published on drupal.org. Currently only supports drush makefiles,
 but composer.lock files (for Drupal 8+) are also planned.

# Installing from scratch

Make sure you have a working Go environment. There should be no specific
version requirements. [Read the Golang install documentation](https://golang.org/doc/install)
on how to do this. Once ready issue the following command:

```$ go get github.com/snufkin/drobot```

you can build the binary using a simple `go build` command.

# Dependencies

The following packages are required for `.make` and `composer.lock` support, respectively:
1. github.com/snufkin/go-composerlockparser
1. github.com/snufkin/go-drushmakeparser

The dependencies should automatically downloaded when you `go get` drobot.

# Usage

Once you have the binary compiled the usage should be straightforward:

```
$ ./drobot parse /path/to/mysite.make
$ ./drobot parse /path/to/mysite/compoers.lock
```

Both of these commands will generate a status report, similarly to how 
[drush](https://www.drush.org/) would do it, when using `drush up`. The main
difference is that this approach does not have to connect to a fully bootstrapped
and functional website, which means, that you don't have to move heavy databases
for a simple monitoring task.

# Planned features

The report output currently is hardcoded, but it would be nice to support various
formats, such as CSV, JSON and the tab delimited format as is currently.
