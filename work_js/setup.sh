#!/bin/sh

###############################################################################
## when setup environment, don't use setup.sh. it's for first configuration'
#npm install

## it requires jq command
#brew install jq

## global
# babel
#npm install -g babel-cli
#babel -V

# gulp
#npm install -g gulp
#npm install -g gulp-babel

# eslint
#npm install -g eslint
#eslint -v
###############################################################################

###############################################################################
## Environment Valiable
###############################################################################
BUILD_MODE=2       #1:Browserify (watchify) + babelify, 2:gulp + gulp-babel

###############################################################################
## for projects
###############################################################################
# create package json
npm init -y

# ES2015
npm install --save-dev babel-preset-es2015

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
  "presets": ["es2015"]
}
EOF


# create src directories
mkdir dist
#touch src/hiromaily.es6.js


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

    # setup gulpfile.js
    touch gulpfile.js
    cat <<EOF > gulpfile.js
var gulp = require('gulp');
var babel = require('gulp-babel');
var plumber = require('gulp-plumber');
var sourcemaps = require("gulp-sourcemaps"); /* source-map */
var rename = require('gulp-regex-rename')

var src = ['src/*.js', 'src/**/*.js'];

gulp.task('babel', function () {
  return gulp.src(src)
    .pipe(plumber())
    .pipe(sourcemaps.init()) /* source-map */
    .pipe(babel())
    .pipe(rename(/\.es6\.js$/, '.js'))
    .pipe(sourcemaps.write(".")) /* source-map */
    .pipe(gulp.dest('./dist'));
});

gulp.task('watch', function () {
  gulp.watch(src, ['babel']);
});

gulp.task('default', ['babel']);
EOF

    ## auto compile
    #gulp watch

fi
