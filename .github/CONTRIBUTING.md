## Contributing

[legal]: https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license
[license]: ../LICENSE

Hi! Thanks for your interest in contributing to this project!

To encourage active collaboration, pull requests are strongly encouraged, not just bug reports. "Bug reports" may also be sent in the form of a pull request containing a failing test. I'd also love to hear about ideas for new features as issues.

Please do:

* Check existing issues to verify that the bug or feature request has not already been submitted.
* Open an issue if things aren't working as expected.
* Open an issue to propose a significant change.
* open an issue to propose a feature
* Open a pull request to fix a bug.
* Open a pull request to fix documentation about a command.
* Open a pull request for an issue with the help-wanted label and leave a comment claiming it.

Please avoid:

* Opening pull requests for issues marked `needs-design`, `needs-investigation`, `needs-user-input`, or `blocked`.
* Opening pull requests for documentation for a new command specifically. Manual pages are auto-generated from source after every release

## Building the project

Prerequisites:
- Go 1.13+

Build with: `make` or `go build -o bin/tfa ./cmd/tfa/main.go`

Run the new binary as: `./bin/tfa`

Run tests with: make test or go test ./...

## Submitting a pull request

1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and ensure tests pass
1. Submit a pull request


Contributions to this project are made available to public under the [project's open source license][license].
Please note that this project adheres to a [Contributor Code of Conduct](https://github.com/profclems/tfa/tree/trunk/.github/CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

Manual pages are auto-generated from source on every release. You do not need to submit pull requests for documentation specifically; manual pages for commands will automatically get updated after your pull requests gets accepted.

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
