import gulp from 'gulp';
import runRequence from 'run-sequence';
import concat from 'gulp-concat';
import clean from 'gulp-clean';
import jade from 'gulp-jade';
import rename from 'gulp-rename';
import through from 'through2';
import path from 'path';



gulp.task('default', () =>
  console.log('xxx')
);

var modifyJade = (file, enc, callback) => {
  through.obj((file, enc, callback) => {
    if (!file.isBuffer()) {
      this.push(file);
      callback();
      return;
    }

    var file_name =
      file.path
        .substring(
          file.path.indexOf('views/')
        )
        .replace('views/', '')
        .replace('.js', '')
        .replace(/\//g, '_');

    var contents =
      file.contents
        .toString()
        .replace('function template(locals) {', 'function tmpl_' + file_name + ' (locals) {');
    file.contents = new Buffer(contents);
    this.push(file);

    callback();
  });
}

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