(http://codecov.io/github/kuzzleio/sdk-cpp/coverage.svg?branch=master)](http://codecov.io/github/kuzzleio/sdk-cpp?branch=master)

Official Kuzzle C SDK
======

## About Kuzzle

A backend software, self-hostable and ready to use to power modern apps.

You can access the Kuzzle repository on [Github](https://github.com/kuzzleio/kuzzle)

## SDK Documentation

The complete SDK documentation is available [here](http://docs.kuzzle.io/sdk-reference/)

## Protocol used

The C SDK implement the websocket protocol.

## Build

### Clone the depot

```sh
git clone --recursive git@github.com:kuzzleio/sdk-c.git
cd sdk-go
git submodule update --init --recursive
make
```

You should know have the SDK in the build directory.

### Installation

You can find prebuilt SDK's for each version and architecture:

arm64: https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-cpp-aarch64-1.0.0.tar.gz

arm32: https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-cpp-armhf-1.0.0.tar.gz

amd64: https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-cpp-amd64-1.0.0.tar.gz

i386:  https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-cpp-i386-1.0.0.tar.gz
