#!/bin/bash

make -C tools

./tools/bin/setup-modules --allow-overwrite
./tools/bin/setup-3p-modules --allow-overwrite
