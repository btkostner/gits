/**
* test/src/helpers.js
* Tests the ability of helper functions
*/

const test = require('ava')

const helpers = require('../lib/helpers')

test('able to make and remove a folder', (t) => {
  t.notThrows(
    helpers.mkdirp('/tmp/gits/handler/testing/nested/thing')
    .then(() => helpers.rmr('/tmp/gits/handler/testing'))
  )
})
