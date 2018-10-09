#!/bin/bash

set -e

if [ "$TRAVIS_PULL_REQUEST" = "false" ] && ([[ "$TRAVIS_BRANCH" == *"-dev" ]] || [ "$TRAVIS_BRANCH" = "master" ]); then

  git config --global user.email "support@kuzzle.io"
  git config --global user.name "Travis CI"
  echo "clone"
  git clone --quiet --branch=${TRAVIS_BRANCH} https://${GH_TOKEN}@github.com/${TRAVIS_REPO_SLUG} travis-repository-copy > /dev/null

  cd travis-repository-copy
  ./align-submodules.sh

  git commit -am "Travis CI - [ci skip] - align submodules to $TRAVIS_BRANCH" > /dev/null 2>&1
  git push origin ${TRAVIS_BRANCH}

fi
