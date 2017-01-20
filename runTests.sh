#!/bin/bash
for test in command config hosts; do
    go test -cover "./${test}"
done
