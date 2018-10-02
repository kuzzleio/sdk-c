(http://codecov.io/github/kuzzleio/sdk-cpp/coverage.svg?branch=master)](http://codecov.io/github/kuzzleio/sdk-cpp?branch=master)

Official Kuzzle C SDK
======

## About Kuzzle

A backend software, self-hostable and ready to use to power modern apps.

You can access the Kuzzle repository on [Github](https://github.com/kuzzleio/kuzzle)

## SDK Documentation

The complete SDK documentation is available [here](http://docs.kuzzle.io/sdk-reference/)

## Protocol used

The C SDK implements the websocket protocol.

### Build

Execute the following snippet to clone the GIT repository, and build the SDK. It will then be available in the "build/" directory

```sh
git clone --recursive git@github.com:kuzzleio/sdk-c.git
cd sdk-c
git submodule update --init --recursive
make
```

### Installation

You can find prebuilt SDK's for each version and architecture:

arm64: https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-c-aarch64-1.0.0.tar.gz

arm32: https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-c-armhf-1.0.0.tar.gz

amd64: https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-c-amd64-1.0.0.tar.gz

x86:  https://dl.kuzzle.io/sdk/c/master/kuzzlesdk-c-x86-1.0.0.tar.gz

You should now have the SDK in the build directory.

