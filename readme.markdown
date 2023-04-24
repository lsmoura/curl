# curl

This is an attempt to implement the well-known app curl, as a drop-in replacement.
It's an exercise to learn more Go and programming patterns.

It's implemented as a module, so it can be imported by other projects, and it has a
command line tool that uses the module that can be installed as a binary on any system 
that go can compile to.

## request

In the spirit of the official `http` package, this module implements a `Request` type
that can be used to make requests.

## author

* [Sergio Moura](https://sergio.moura.ca)
