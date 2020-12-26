## tss

A Go port of moreutils/ts and fork of [kevinburke/tss].

[kevinburke/tss]: https://github.com/kevinburke/tss

`tss` is like `ts` from moreutils,
but prints relative durations (with millisecond precision) by default,
and can be shipped as a compiled binary.

Try it out:

```console
$ (sleep 1; echo "hello"; sleep 2; echo "two sec") | tss
   995ms          hello
      3s   2.005s two sec
```

The first column is the amount of time that has elapsed since the program started.
The second column is the amount of time that has elapsed since the last line printed.

[![license](https://img.shields.io/github/license/peaceiris/tss.svg)](https://github.com/peaceiris/tss/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/peaceiris/tss.svg)](https://github.com/peaceiris/tss/releases/latest)
[![GitHub release date](https://img.shields.io/github/release-date/peaceiris/tss.svg)](https://github.com/peaceiris/tss/releases)
[![Release Feed](https://img.shields.io/badge/release-feed-yellow)](https://github.com/peaceiris/tss/releases.atom)
![Test](https://github.com/peaceiris/tss/workflows/CI/badge.svg?branch=main&event=push)
![Code Scanning](https://github.com/peaceiris/tss/workflows/Code%20Scanning/badge.svg?event=push)

[![CodeFactor](https://www.codefactor.io/repository/github/peaceiris/tss/badge)](https://www.codefactor.io/repository/github/peaceiris/tss)
[![Maintainability](https://api.codeclimate.com/v1/badges/5eaad3d1e44d6eb87a95/maintainability)](https://codeclimate.com/github/peaceiris/tss/maintainability)



## Installation

### homebrew

For macOS and Linux.

```sh
brew install peaceiris/tap/tss
```

### Binary

[Releases · peaceiris/tss](https://github.com/peaceiris/tss/releases)
