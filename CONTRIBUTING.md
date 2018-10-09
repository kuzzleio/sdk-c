# How to contribute to the C SDK

Here are a few rules and guidelines to follow if you want to contribute to the Java SDK and, more importantly, if you want to see your pull requests accepted by the  Kuzzle team.

## Tools

This SDK inherits from the following repositories, linked as git submodules: sdk-go.  
Whenever significant changes are applied to the parent SDKs, you need to align the linked submodules accordingly.
You can use the `align-submodules.sh` script to achieve this. (e.g. `./align-submodules.sh 1-dev` to align all submodules on `1-dev` branch)


You can use this Docker image to build the SDK:  
```
docker run --rm -it -v "$(pwd)":/mnt kuzzleio/sdk-cross:gcc make
```
