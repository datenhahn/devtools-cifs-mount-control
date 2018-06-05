#!/bin/bash

./build.sh

docker run --rm -v $PWD:/code -ti snapcore/snapcraft bash -c "cd /code;snapcraft clean;snapcraft"
