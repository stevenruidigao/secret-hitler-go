#!/bin/bash

yarn webpack --mode development webpack/webpack.config.dev.js && gofmt -w *.go && gofmt -w */*.go && go build && ./secrethitler.io
