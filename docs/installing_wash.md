---
title: Installing Wash
---
Wash is distributed as a single binary; the only prerequisite is `libfuse`. Here’s how to install it.

* Download the Wash binary for your platform
    * or install it with `brew install puppetlabs/puppet/wash`

* Install `libfuse` if you haven’t already
    * E.g. on MacOS using homebrew: `brew cask install osxfuse`
    * E.g. on CentOS: `yum install fuse fuse-libs`
    * E.g. on Ubuntu: `apt-get install fuse`

* Run `wash version` to ensure that the installation was successful

<!-- TODO: Link to analytics docs -->
**NOTE:** Wash collects anonymous data about how you use it. See the analytics docs for more details.
