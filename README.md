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

# Usage

Once you have the binary compiled the usage should be straightforward:

`./drobot /path/to/mysite.make` will generate the status report, similarly to
how [drush](https://www.drush.org/) would do it, but this approach does not
require to have a fully bootstrapped and functional website.
