#!/usr/bin/env bash
rm -rf public/*

git clone git@github.com:quentin-sommer/gomoku-front.git \
  && cd gomoku-front \
  && yarn install \
  && echo "Building application..." \
  && npm run build &>/dev/null \
  && cp -r build/* ./../public/ \
  && cd .. \
  && rm -rf gomoku-front \
  && echo "Application built!"
