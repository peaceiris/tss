## tss

Annotate stdin with timestamps per line.
A Go port of [moreutils](https://joeyh.name/code/moreutils/)/ts and fork of [kevinburke/tss].

[kevinburke/tss]: https://github.com/kevinburke/tss

`tss` is like `ts` from moreutils,
but prints relative durations (with millisecond precision) by default,
and can be shipped as a compiled binary.

Try it out:

```console
$ (sleep 1; echo "hello"; sleep 2; echo "two sec") | tss
   1.00s    1.00s hello
   3.01s    2.01s two sec
```

The first column is the amount of time that has elapsed since the program started.
The second column is the amount of time that has elapsed since the last line printed.

[![license](https://img.shields.io/github/license/peaceiris/tss.svg)](https://github.com/peaceiris/tss/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/peaceiris/tss.svg)](https://github.com/peaceiris/tss/releases/latest)
[![GitHub release date](https://img.shields.io/github/release-date/peaceiris/tss.svg)](https://github.com/peaceiris/tss/releases)
[![Release Feed](https://img.shields.io/badge/release-feed-yellow)](https://github.com/peaceiris/tss/releases.atom)
![Test](https://github.com/peaceiris/tss/workflows/CI/badge.svg?branch=main&event=push)
![Code Scanning](https://github.com/peaceiris/tss/workflows/Code%20Scanning/badge.svg?event=push)

[![Go Report Card](https://goreportcard.com/badge/github.com/peaceiris/tss)](https://goreportcard.com/report/github.com/peaceiris/tss)
[![CodeFactor](https://www.codefactor.io/repository/github/peaceiris/tss/badge)](https://www.codefactor.io/repository/github/peaceiris/tss)
[![Maintainability](https://api.codeclimate.com/v1/badges/5eaad3d1e44d6eb87a95/maintainability)](https://codeclimate.com/github/peaceiris/tss/maintainability)
[![codecov](https://codecov.io/gh/peaceiris/tss/branch/main/graph/badge.svg?token=4119ASAR7K)](https://codecov.io/gh/peaceiris/tss)



## Installation

### Homebrew

For macOS and Linux.

```sh
brew install peaceiris/tap/tss
```

[peaceiris/homebrew-tap/tss.rb](https://github.com/peaceiris/homebrew-tap/blob/main/Formula/tss.rb)

### Binary

[Releases · peaceiris/tss](https://github.com/peaceiris/tss/releases)
