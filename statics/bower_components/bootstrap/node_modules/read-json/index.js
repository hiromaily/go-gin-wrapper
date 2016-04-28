var fs = require("fs");
var call = require("try-call");

module.exports = async;
module.exports.sync = sync;

function async (filename, options, callback) {
  if(arguments.length == 2){
    callback = options;
    options = {};
  }

  fs.readFile(filename, options, function(error, bf){
    if(error) return callback(error);
    call(parse.bind(null, bf), callback);
  });
}

function sync (filename, options) {
  return parse(fs.readFileSync(filename, options));
}

function parse (bf) {
  return JSON.parse(bf.toString().replace(/^\ufeff/g, ''));
}
