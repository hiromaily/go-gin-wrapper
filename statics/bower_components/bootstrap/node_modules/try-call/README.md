## try-call

Functional try-catch for cleaner code & [optimization](https://github.com/petkaantonov/bluebird/wiki/Optimization-killers).

## Install

```bash
$ npm install try-call
```

## Usage

```js
var call = require('try-call');

var doc = '{ "foo": 123 }';
var parse = JSON.parse.bind(null, doc);

call(parse, function (error, doc) {
  error
  // => undefined

  doc
  // => { foo: 123}
})
```
