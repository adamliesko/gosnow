# Gosnow

A Go library for handling the [API Blueprint](https://apiblueprint.org) format.

## Usage

Setup the inner drafter directory with: 
`git submodule update --init --recursive`

Install the drafter dylib with `make install`

Run the pure c tests with `make test`

## Issues

linking to the `libdrafter.dylib` is currently done by linking the library to the global scope in `/usr/local/lib/`. It would be much preferred to have the dylib found locally.

This is the runtime error that occurs when the global library is not present
```
dyld: Library not loaded: /usr/local/lib/libdrafter.dylib
  Referenced from: /Users/ataylor/rastech/gosnow/./test
  Reason: image not found
```
