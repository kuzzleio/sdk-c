#!/usr/bin/env bash
set -e

mkdir ./build/{lib,include}
cp -fr ./include/*.h ./build/include
cp ./build/*.{so,a}  ./build/lib/
mkdir deploy
cd build
tar cfz ../deploy/kuzzlesdk-c-$ARCH.tar.gz lib include

