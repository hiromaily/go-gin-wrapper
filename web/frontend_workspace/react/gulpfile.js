var gulp = require('gulp');
var browserify = require('browserify');
var babelify = require('babelify');
var uglify = require('gulp-uglify');
var source = require('vinyl-source-stream');
var buffer = require('vinyl-buffer');
var webserver = require('gulp-webserver');
//var beautify = require('gulp-beautify');

var inDir = './app/src/'
var outDir = './app/dist/'
var inFiles = ['apilist'];
//var inFile = './app/src/index.js';
//var outFile = './app/dist/index.bundle.js';

//browserify for release version
gulp.task('release', function() {
  process.env.NODE_ENV = 'production';
  inFiles.forEach(function(file,i,ar){
    browserify(inDir+file+".js", { debug: false }) //debug: true is for sourcemap
      .transform(babelify)
      .transform('browserify-shim', { global: true })
      .bundle()
      .on("error", function (err) { console.log("Error : " + err.message); })
      .pipe(source(outDir+file+".bundle.js"))
      .pipe(buffer())
      .pipe(uglify())
      .pipe(gulp.dest('./'))
  });
});

//browserify for rdev version
gulp.task('dev', function() {
  //_.each srcname, (inFiles) ->
  inFiles.forEach(function(file,i,ar){
    browserify(inDir+file+".js", { debug: true }) //debug: true is for sourcemap
      .transform(babelify)
      .transform('browserify-shim', { global: true })
      .bundle()
      .on("error", function (err) { console.log("Error : " + err.message); })
      .pipe(source(outDir+file+".bundle.js"))
      .pipe(gulp.dest('./'))
  });
});

//watch
gulp.task('watch', function() {
  //gulp.watch(/.jsx?$/, ['browserify'])
  //gulp.watch('./*.jsx', ['browserify'])
  gulp.watch('**/*.jsx', ['dev'])
});

//beautify
//gulp.task('beautify', function() {
//  gulp.src(['./app/src/*.js', './app/components/**/*.js'])
//    .pipe(beautify({indentSize: 2}))
//    .pipe(gulp.dest('./public/'))
//});

//webserver
gulp.task('web', function() {
  gulp.src('app')
    .pipe(webserver({
      host: '127.0.0.1',
      livereload: true
    })
  );
});

gulp.task('do', ['dev', 'watch', 'web']);
gulp.task('default', ['release', 'watch', 'web']);


//cp /Users/hy/work/go/src/github.com/hiromaily/go-gin-wrapper/frontend_workspace/react/app/dist/apilist.bundle.js \
///Users/hy/work/go/src/github.com/hiromaily/go-gin-wrapper/statics/js/
