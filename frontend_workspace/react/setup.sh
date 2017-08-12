#!/bin/sh

###############################################################################
## Environment Valiable
###############################################################################



###############################################################################
## Gulp Environment
###############################################################################
#https://github.com/hkusu/react-babel-browserify-gulp-sample
#http://qiita.com/hkusu/items/e068bba0ae036b447754

#React
#babel
#gulp

#http://127.0.0.1:8000

###############################################################################
# Run
###############################################################################
#gulp watch
#gulp browserify
#gulp webserver

###############################################################################
## Setup
###############################################################################
npm init

npm install -g gulp
npm install -g gulp-babel
npm install -g eslint
npm install -g eslint-plugin-react
npm install -g gulp-beautify

npm i -S react react-dom flux

npm i -D gulp
npm i -D gulp-babel gulp-plumber
npm i -D gulp-uglify

npm i -D babel-core babel-loader
npm i -D babel-preset-es2015 babel-preset-react

npm i -D browserify
npm i -D browserify-shim
npm i -D babelify
npm i -D gulp-webserver
npm i -D vinyl-source-stream
npm i -D vinyl-buffer

npm i -D jquery

# .babelrc
touch .babelrc
cat <<EOF > .babelrc
{
  "presets": ["react", "es2015"]
}
EOF

# rewrite package.json
cp -r package.json package.json.bk

"build": ""

BUILD="browserify --debug --transform babelify main.js --outfile bundle.js"
LINT="eslint src/*.js"

cat package.json |
jq --arg BUILD "$BUILD" --arg LINT "$LINT" 'to_entries |
    map(if .key == "scripts"
        then . + {"value":
                    {
                      "build": $BUILD,
                      "lint": $LINT
                    }
                 }
        else .
        end
    ) | from_entries' |
jq 'del(.main) | del(.keywords) | del(.author) | del(.license)' >> tmp.json

rm -rf package.json
mv tmp.json package.json
rm -rf tmp.json

# gulpfile.js
touch gulpfile.js
cat <<EOF > gulpfile.js
var gulp = require('gulp');
var browserify = require('browserify');
var babelify = require('babelify');
var source = require('vinyl-source-stream');
var webserver = require('gulp-webserver');

gulp.task('browserify', function() {
  browserify('./app.jsx', { debug: true })
    .transform(babelify)
    .bundle()
    .on("error", function (err) { console.log("Error : " + err.message); })
    .pipe(source('bundle.js'))
    .pipe(gulp.dest('./'))
});

gulp.task('watch', function() {
  gulp.watch('./*.jsx', ['browserify'])
});

gulp.task('webserver', function() {
  gulp.src('./')
    .pipe(webserver({
      host: '127.0.0.1',
      livereload: true
    })
  );
});

gulp.task('default', ['browserify', 'watch', 'webserver']);
EOF
