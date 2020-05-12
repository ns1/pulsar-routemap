Release process
===============

> N.B. This section is for maintainers.

Prerequisites: 

* We use [GoReleaser](https://goreleaser.com/) to perform releases so make sure
that's installed and setup.
* You'll need to create an Access Token for Github. 
  * You create your token here: https://github.com/settings/tokens
  * Make sure it has at least **repo** access.
  * Once you have the token, add it to your environment:  `export GITHUB_TOKEN=xyz`


You can do a dry run like this:

```sh
$ goreleaser --skip-publish --snapshot
```

To do an actual release:

```sh
# Make sure everything is on master. Then create an push a new tag for your release
$ git tag -a x.y.z -m "Creating x.y.z release: blah blah"
$ git push origin x.y.z

# Now do the release!
$ goreleaser
```
