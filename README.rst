================================================================================
 building debian package for the session of Tokyo Debian Meeting #115 (2014.07)
================================================================================

Purpose
-------

This script is building debian package on Jenkins for the session of Tokyo Debian Meeting #115 (2014.07),
This script executes as follows;

* Clean building from git repository or source packages with pbuilder and cowbuilder.
* Tesging to install/unstall with lintian.
* Signing with GPG Public key.
* Uploading to local archives managed with reprepro.
