var gulp = require('gulp');
var babel = require('gulp-babel');
var plumber = require('gulp-plumber');
var sourcemaps = require("gulp-sourcemaps"); /* source-map */
var rename = require('gulp-regex-rename');

var browserify = require('browserify');
var babelify   = require('babelify');

var source   = require('vinyl-source-stream');


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

gulp.task('browserify', function () {
  return browserify('./src/module_import.es6.js')
        .transform(babelify, {presets: ['es2015']})
        .bundle()
        .pipe(source('bundle.js'))
        .pipe(gulp.dest('./dist'));
});

gulp.task('watch', function () {
  gulp.watch(src, ['babel']);
});

gulp.task('default', ['babel']);
