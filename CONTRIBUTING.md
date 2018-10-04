# How to contribute to the C SDK

Here are a few rules and guidelines to follow if you want to contribute to the C SDK and, more importantly, if you want to see your pull requests accepted by Kuzzle team.

## Tools

We use git submodules to link the sdk-go and the sdk-c.  
When you are developing a new functionality that had implications on the other SDK, you should align all your submodules on your development branch.  
You can use `align-submodules.sh` script to achieve this. (Eg: `./align-submodules.sh 1-dev` to align all submodules on `1-dev` branch)


To build the SDK, you can use this Docker image to build the SDK:  
```
docker run --rm -it -v "$(pwd)":/mnt kuzzleio/sdk-cross:gcc make
```
