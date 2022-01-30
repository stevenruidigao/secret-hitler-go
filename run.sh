#!/bin/bash

gofmt -w *.go && gofmt -w */*.go && go build && ./secrethitler.io
