#!/bin/bash

for arg
do
arr=(${arg//.pug/ })
./jaded $arg > ${arr[0]}.tmpl
done
