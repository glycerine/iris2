# Change Log
All notable changes to this project will be documented in this file.

**Upgrading**: Use: `go get -u github.com/go-iris2/iris2` and you are done.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Added
- File storage for sessions
- Leveldb storage for sessions

### Changed
- Fork from kataras/iris to go-iris2/iris2 and rename (`4b71e60`)
- Default routing uses httprouter

### Removed
- Option for gorillamux in favor of user-friendlyness

### Fixed
- Sessions bound to IP-addresses
- Possible session-collisions
