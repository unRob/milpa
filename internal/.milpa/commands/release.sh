#!/usr/bin/env bash

cores="$(sysctl -n hw.physicalcpu 2>/dev/null || grep -c ^processor /proc/cpuinfo)"

cd "$MILPA_ROOT" || _fail "unknown root"
# make clean
make -j"$cores" dist/release || _fail "Could not complete compa build"
