# Contributing

Thank you for your interest in this project!

We use *GitHub Issues* for trackings issues and features. You can make a contribution by:

1. Reporting an issue or making a feature request [here](https://github.com/cloud-ca/terraform-provider-cloudca/issues).
2. Contributing code to this project by fixing an issue or adding a new feature (see below).

Before contributing a new feature, please discuss its suitability with the project maintainers in an issue first. Thanks!

## Contribution Process

1. Fork and *git clone* [cca](https://github.com/cloud-ca/terraform-provider-cloudca).
2. Create a new *git branch* from the master branch where you develop your changes.
3. Create a [Pull Request](https://help.github.com/articles/about-pull-requests/) for your contribution by following the instructions in the pull request template [here](https://github.com/cloud-ca/terraform-provider-cloudca/pull).
4. Perform a code review with the project maintainers on the pull request. We may suggest changes, improvements or alternatives.
5. Once approved, your code will be merged into `master`.

### Requirements

Pull requests have to meet the following requirements:

1. **Tests**: Code changes need to be tested with code and tests being located in the same folder. Make sure that your tests pass using `make test`.

2. **Documentation**: Pull requests need to update the [documentation](https://github.com/cloud-ca/terraform-provider-cloudca/tree/master/README.md) together with the code change.

3. **Commits**: Commits should be as small as possible while ensuring that each commit compiles and passes tests independently. [Write good commit messages](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html). If needed, [squash your commits](https://davidwalsh.name/squash-commits-git) prior to submission.

4. **Code Style**: Use [gofmt](https://blog.golang.org/go-fmt-your-code) to format your code. If useful, include code comments to support your intentions.

## Additional Resources

- [Golang Basics: Writing Unit Tests (Alex Ellis)](https://blog.alexellis.io/golang-writing-unit-tests/)
- [Advanced Testing in Go (Mitchell Hashimoto)](https://about.sourcegraph.com/go/advanced-testing-in-go/)
