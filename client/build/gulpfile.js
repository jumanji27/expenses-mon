// var gulp = require('gulp'),
//   ts = require('gulp-typescript'),
//   run_sequence = require('run-sequence'),
//   concat = require('gulp-concat'),
//   clean = require('gulp-clean'),
//   jade = require('gulp-jade'),
//   less = require('gulp-less'),
//   rename = require('gulp-rename'),
//   through = require('through2'),
//   path = require('path');


// var jade_modify = function() {
//   return through.obj(function (file, enc, callback) {
//     if (!file.isBuffer()) {
//       this.push(file);
//       callback();
//       return;
//     }

//     file_name =
//       file.path
//         .substring(
//           file.path.indexOf('views/')
//         )
//         .replace('views/', '')
//         .replace('.js', '')
//         .replace(/\//g, '_');

//     var contents =
//       file.contents
//         .toString()
//         .replace('function template(locals) {', 'function tmpl_' + file_name + ' (locals) {');
//     file.contents = new Buffer(contents);
//     this.push(file);

//     callback();
//   });
// }


// gulp.task('clean_public', function() {
//   return gulp.src(
//       ['../public/**/*', '!../public/storage', '!../public/storage/**/*.jpg'],
//       {read: false}
//     )
//     .pipe(
//       clean({force: true})
//     );
// });

// gulp.task('concat_vendor', function() {
//   return gulp.src(
//       [
//         'bower_components/jquery/jquery.js',
//         'bower_components/underscore/underscore.js',
//         'bower_components/backbone/backbone.js',
//         'bower_components/jade/jade.js',
//         'bower_components/jade/runtime.js'
//       ]
//     )
//     .pipe(concat('vendor.js'))
//     .pipe(gulp.dest('tmp'));
// });

// gulp.task('compile_typescript', function () {
//   return gulp.src(
//       ['../src/collections/**/*.ts', '../src/models/**/*.ts', '../src/views/**/*.ts', '../src/router.ts']
//     )
//     .pipe(
//       ts({out: 'ts.js'})
//     )
//     .js.pipe(gulp.dest('tmp'));
// });

// gulp.task('compile_jade', function() {
//   return gulp.src('../src/**/*.jade')
//     .pipe(
//       jade({client: true})
//     )
//     .pipe(jade_modify())
//     .pipe(gulp.dest('tmp/jade'))
// });

// gulp.task('concat_js', function() {
//   return gulp.src(
//       ['tmp/vendor.js', 'tmp/ts.js', 'tmp/jade/**/*.js']
//     )
//     .pipe(concat('index.js'))
//     .pipe(gulp.dest('../public'));
// });

// gulp.task('compile_less', function () {
//   return gulp.src('../src/**/*.less')
//     .pipe(less())
//     .pipe(gulp.dest('tmp/css'));
// });

// gulp.task('concat_css', function() {
//   return gulp.src('tmp/css/**/*.css')
//     .pipe(concat('index.css'))
//     .pipe(gulp.dest('../public'));
// });

// gulp.task('move_img', function(){
//   return gulp.src(
//       ['../src/**/*.png', '../src/**/*.jpg']
//     )
//     .pipe(
//       rename({dirname: ''})
//     )
//     .pipe(gulp.dest('../public/img'));
// });

// gulp.task('clean_tmp', function() {
//   return gulp.src(
//       'tmp',
//       {read: false}
//     )
//     .pipe(clean());
// });


// gulp.task('build', function() {
//   run_sequence(
//     'clean_public',
//     ['concat_vendor', 'compile_typescript', 'compile_jade', 'compile_less', 'move_img'],
//     ['concat_js', 'concat_css'],
//     'clean_tmp'
//   );
// });


// gulp.task('run', function () {
//   run_sequence('build', function() {
//     gulp.watch(['../src/**/*.ts', '../src/**/*.jade', '../src/**/*.less'], ['build']);
//   });
// });