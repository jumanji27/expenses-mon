'use strict';

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }

var _gulp = require('gulp');

var _gulp2 = _interopRequireDefault(_gulp);

var _gulpBabel = require('gulp-babel');

var _gulpBabel2 = _interopRequireDefault(_gulpBabel);

var _runSequence = require('run-sequence');

var _runSequence2 = _interopRequireDefault(_runSequence);

var _gulpConcat = require('gulp-concat');

var _gulpConcat2 = _interopRequireDefault(_gulpConcat);

var _gulpClean = require('gulp-clean');

var _gulpClean2 = _interopRequireDefault(_gulpClean);

var _gulpJade = require('gulp-jade');

var _gulpJade2 = _interopRequireDefault(_gulpJade);

var _gulpRename = require('gulp-rename');

var _gulpRename2 = _interopRequireDefault(_gulpRename);

var _through2 = require('through2');

var _through22 = _interopRequireDefault(_through2);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var modifyJade = function modifyJade(file, enc, callback) {
  _through22['default'].obj(function (file, enc, callback) {
    if (!file.isBuffer()) {
      undefined.push(file);
      callback();
      return;
    }

    var file_name = file.path.substring(file.path.indexOf('views/')).replace('views/', '').replace('.js', '').replace(/\//g, '_');

    var contents = file.contents.toString().replace('function template(locals) {', 'function tmpl_' + file_name + ' (locals) {');
    file.contents = new Buffer(contents);
    undefined.push(file);

    callback();
  });
};

_gulp2['default'].task('compile_babel', function () {
  return _gulp2['default'].src('../**/*.js').pipe((0, _gulpBabel2['default'])()).pipe(_gulp2['default'].dest('public'));
});

// gulp.task('clean_public', () => {
//   return gulp.src(
//       ['../public/**/*', '!../public/storage', '!../public/storage/**/*.jpg'],
//       {read: false}
//     )
//     .pipe(
//       clean({force: true})
//     );
// });

// gulp.task('concat_vendor', () => {
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

// gulp.task('compile_jade', () => {
//   return gulp.src('../src/**/*.jade')
//     .pipe(
//       jade({client: true})
//     )
//     .pipe(modifyJade())
//     .pipe(gulp.dest('tmp/jade'))
// });

// gulp.task('concat_js', () => {
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

// gulp.task('concat_css', () => {
//   return gulp.src('tmp/css/**/*.css')
//     .pipe(concat('index.css'))
//     .pipe(gulp.dest('../public'));
// });

// gulp.task('move_img', () =>{
//   return gulp.src(
//       ['../src/**/*.png', '../src/**/*.jpg']
//     )
//     .pipe(
//       rename({dirname: ''})
//     )
//     .pipe(gulp.dest('../public/img'));
// });

// gulp.task('clean_tmp', () => {
//   return gulp.src(
//       'tmp',
//       {read: false}
//     )
//     .pipe(clean());
// });

// gulp.task('build', () => {
//   runRequence(
//     'clean_public',
//     ['concat_vendor', 'compile_typescript', 'compile_jade', 'compile_less', 'move_img'],
//     ['concat_js', 'concat_css'],
//     'clean_tmp'
//   );
// });

// gulp.task('run', function () {
//   runRequence('build', () => {
//     gulp.watch(['../src/**/*.ts', '../src/**/*.jade', '../src/**/*.less'], ['build']);
//   });
// });