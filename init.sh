#!/bin/bash

./zwallet register --verbose --devserver

./zwallet faucet --methodName pour --input "{Pay day}" --verbose --devserver
./zwallet faucet --methodName pour --input "{Pay day}" --verbose --devserver
./zwallet faucet --methodName pour --input "{Pay day}" --verbose --devserver
./zwallet faucet --methodName pour --input "{Pay day}" --verbose --devserver
./zwallet faucet --methodName pour --input "{Pay day}" --verbose --devserver
./zwallet faucet --methodName pour --input "{Pay day}" --verbose --devserver
./zwallet faucet --methodName pour --input "{Pay day}" --verbose --devserver

./zwallet getbalance --verbose --devserver
