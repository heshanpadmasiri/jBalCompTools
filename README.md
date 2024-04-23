# JBalCompTools
A collection of various utility functions I use regularly unified under a easy to use cli tool. For a list of currently implemented as well as planned features checkout [TODO.md](TODO.md) file
## Install
```sh
go install github.com/heshanpadmasiri/jBalCompTools@latest
```
## Configuration
This program expect a configuration file at `~/.config/jBalCompTools/config.toml` with fallowing data
```toml
sourcePath = "PATH TO ballerina-lang repo"
version = "version to the jBallerina you are building"
```
