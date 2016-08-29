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
        .on("error", function (err) { console.log("Error : " + err.message)})
        .pipe(source('bundle.js'))
        .pipe(gulp.dest(out))
})

gulp.task('watch', function () {
  gulp.watch(src, ['babel'])
})

gulp.task('default', ['babel'])
