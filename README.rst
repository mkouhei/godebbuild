=====================================
 building debian package for Jenkins
=====================================

Purpose
-------

This tool is building debian package on Jenkins
This tool executes as follows;

* Clean building from git repository or source packages with pbuilder and cowbuilder.
* Tesging to install/unstall
* Signing with GPG Public key.
* Uploading to local archives managed with reprepro with lintian test.

Requirements
------------

Debian packages
~~~~~~~~~~~~~~~

* devscripts
* dpkg-dev
* pbuilder
* cowbuilder
* git-buildpackage
* piuparts
* dput
* openssh-client
