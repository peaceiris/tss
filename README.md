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
