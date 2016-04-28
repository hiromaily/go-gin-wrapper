module.exports = call;

function call (fn, callback) {
  var result;

  try {
    result = fn();
  } catch (err) {
    return callback(err);
  }

  callback(undefined, result);
}
