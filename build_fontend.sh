#!/usr/bin/env bash
rm -rf public/*
git clone git@github.com:quentin-sommer/gomoku-front.git
cd gomoku-front
yarn install || npm install
npm run build
cp -r build/* ././public/
cd ..
rm -rf gomku-front
