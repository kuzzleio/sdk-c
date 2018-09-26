#!/usr/bin/env bash
set -e

mkdir ./build/{lib,include}
cp -fr ./include/*.h ./build/include
cp ./build/*.{so,a}  ./build/lib/
mkdir deploy
tar cfz deploy/kuzzlesdk-c-$ARCH.tar.gz ./build/lib ./build/include
