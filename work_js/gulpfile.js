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
