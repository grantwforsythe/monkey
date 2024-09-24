# Monkeylang
[![codecov](https://codecov.io/gh/grantwforsythe/monkey/graph/badge.svg?token=B6KW6CCWHY)](https://codecov.io/gh/grantwforsythe/monkey)
[![CI](https://github.com/grantwforsythe/monkeylang/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/grantwforsythe/monkeylang/actions/workflows/ci.yml)

## Get Started
```sh
docker build . -t monkey
docker run -i monkey
```

## Example
```

let name = "Monkey";
let age = 1;
let inspirations = ["Scheme", "Lisp", "JavaScript", "Clojure"];
let book = {
  "title": "Writing A Compiler In Go",
  "author": "Thorsten Ball",
  "prequel": "Writing An Interpreter In Go"
};

let printBookName = fn(book) {
    let title = book["title"];
    let author = book["author"];
    puts(author + " - " + title);
};

printBookName(book);
// => prints: "Thorsten Ball - Writing A Compiler In Go"

let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      return 1;
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};

let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(rest(arr), push(accumulated, f(first(arr))));
    }
  };

  iter(arr, []);
};

let numbers = [1, 1 + 1, 4 - 1, 2 * 2, 2 + 3, 12 / 2];
map(numbers, fibonacci);
// => returns: [1, 1, 2, 3, 5, 8]```

## TODO
- [] Web app to play with the interpreter
- [] Compiler
- [] LSP
