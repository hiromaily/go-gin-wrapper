#!/bin/sh

###############################################################################
## Environment Valiable
###############################################################################
EXEC_ONCE=0         #0:off, 1:exec (it's ok just once execution')
BUILD_MODE=0        #0:off, 1:Browserify (watchify) + babelify,
                    #       2:gulp + gulp-babel
REACT_FLG=1

###############################################################################
## Initialize environment
###############################################################################
if [ $EXEC_ONCE -eq 1 ]; then

    ## when setup environment, don't use setup.sh. it's for first configuration'
    npm install

    ## it requires jq command
    brew install jq

    # babel
    npm install -g babel-cli
    #babel -V

    # gulp
    npm install -g gulp
    npm install -g gulp-babel

    # eslint
    npm install -g eslint
    npm install -g eslint-plugin-react
    #eslint -v
fi


###############################################################################
## for projects
###############################################################################
#if [ $BUILD_MODE -ge 1 ]; then
if [ $BUILD_MODE -ne 0 ]; then

    # create package json
    npm init -y

    # ES2015
    npm install --save-dev babel-preset-es2015
    npm install --save-dev babel-plugin-transform-es2015-modules-commonjs

    # setup eslint
    touch .eslintrc.json
    cat <<EOF > .eslintrc.json
{
    "extends": ["eslint:recommended"],
    "plugins": [],
    "parserOptions": {
        "ecmaVersion": 6,
        "sourceType": "module",
        "ecmaFeatures": {
            "jsx": true
        }
    },
    "env": {
        "browser": true,
        "es6": true
    },
    "globals": {
        "$": false
    },
    "rules": {
        "no-console":0,
        "semi" : [2 , "never"]
    }
}
EOF


    # setup babel
    touch .babelrc
    cat <<EOF > .babelrc
{
  "presets": ["react", "es2015"],
  "plugins": ["transform-es2015-modules-commonjs"]
}
EOF

    # create src directories
    mkdir dist
    #touch src/hiromaily.es6.js
fi

#####################################
# Browserify (watchify) + babelify
#####################################
if [ $BUILD_MODE -eq 1 ]; then

    # watchify is Browserify with watching tool
    npm install --save-dev watchify

    # tools added babel features for Browserify
    npm install --save-dev babelify

    # for output source map
    npm install --save-dev exorcist

    # rewrite package.json
    cp -r package.json package.json.bk

    WATCH="watchify -t babelify ./src/hiromaily.es6.js -o 'exorcist ./dist/hiromaily.js.map > ./dist/hiromaily.js' -d"
    LINT="eslint src/*.js"

    cat package.json |
    jq --arg WATCH "$WATCH" --arg LINT "$LINT" 'to_entries |
        map(if .key == "scripts"
            then . + {"value":
                        {
                    "watch": $WATCH,
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


    ## auto compile
    #npm run watch

    ## lint
    #npm run lint


#####################################
# gulp + gulp-babel
#####################################
elif [ $BUILD_MODE -eq 2 ]; then
    # gulp
    npm install --save-dev gulp
    npm install --save-dev gulp-babel gulp-plumber
    npm install --save-dev gulp-sourcemaps
    #npm install --save-dev gulp-exec
    #npm install --save-dev gulp-rename
    npm install --save-dev gulp-regex-rename

    # browserify
    npm install --save-dev browserify
    npm install --save-dev babelify
    npm install --save-dev vinyl-source-stream

    # others
    npm install --save-dev superagent
    
    # setup gulpfile.js
    touch gulpfile.babel.js
    cat <<EOF > gulpfile.bebel.js
import gulp from 'gulp'
import babel from 'gulp-babel'
import plumber from 'gulp-plumber'
import sourcemaps from 'gulp-sourcemaps'
import rename from 'gulp-regex-rename'

import browserify from "browserify"
import babelify from "babelify"
import source from "vinyl-source-stream"


var src = ['src/*.js', 'src/**/*.js']
var out = './dist'


gulp.task('babel', () => {
  return gulp.src(src)
    .pipe(plumber())
    .pipe(sourcemaps.init()) /* source-map */
    .pipe(babel())
    .pipe(rename(/\.es6\.js$/, '.js'))
    .pipe(sourcemaps.write(".")) /* source-map */
    .pipe(gulp.dest(out))
})

gulp.task('browserify', function () {
  return browserify('./src/module_import.es6.js')
        .transform(babelify, {presets: ['es2015']})
        .bundle()
        .pipe(source('bundle.js'))
        .pipe(gulp.dest(out))
})

gulp.task('watch', function () {
  gulp.watch(src, ['babel'])
})

gulp.task('default', ['babel'])
EOF

    ## auto compile
    #gulp watch
    #gulp browserify

fi


#####################################
# REACT
#####################################
if [ $REACT_FLG -eq 1 ]; then
    npm install -g react-tools
    npm install --save react react-dom
    #npm install --save react-redux
    npm install --save-dev babel-preset-react
fi
