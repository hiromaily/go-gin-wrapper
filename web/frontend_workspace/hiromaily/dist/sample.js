'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var ES6 = 'ECMAScript2015';

// function log
var log = function log() {
    console.log(ES6);
    console.log("ECMAScript2015");
};

// class Hy

var Hy = function () {
    function Hy(x, y) {
        _classCallCheck(this, Hy);

        //public
        this.x = x;
        this.y = y;
    }

    //public


    _createClass(Hy, null, [{
        key: 'distance',
        value: function distance(a, b) {
            var dx = a.x - b.x;
            var dy = a.y - b.y;

            return Math.sqrt(dx * dx + dy * dy);
        }
    }]);

    return Hy;
}();

//main()


function main() {
    alert(11);

    log();

    //fmt.Printf的な
    var name = 'Harry';
    console.log('Hello, ' + name);

    //class
    var p1 = new Hy(5, 5);
    var p2 = new Hy(10, 10);

    console.log(Hy.distance(p1, p2));
}

main();
//# sourceMappingURL=sample.js.map
