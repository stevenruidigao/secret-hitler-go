#!/bin/bash

yarn webpack --mode production --config webpack/webpack.config.prod.js ; gofmt -w *.go && gofmt -w */*.go && go build && ./secrethitler.io
