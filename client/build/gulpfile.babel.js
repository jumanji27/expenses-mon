import gulp from 'gulp';
import runRequence from 'run-sequence';
import babel from 'gulp-babel';
import jade from 'gulp-jade';
import stylus from 'gulp-stylus';
import concat from 'gulp-concat';
import clean from 'gulp-clean';
import rename from 'gulp-rename';
import through from 'through2';
import path from 'path';


let modifyJade = function() {
  return through.obj(function (file, enc, callback) {
    if (!file.isBuffer()) {
      this.push(file);
      callback();
      return;
    }

    let file_name =
      file.path
        .substring(
          file.path.indexOf('views/')
        )
        .replace('views/', '')
        .replace('.js', '')
        .replace(/\//g, '_');

    let contents =
      file.contents
        .toString()
        .replace('function template(locals) {', 'function tmpl_' + file_name + ' (locals) {');
    file.contents = new Buffer(contents);
    this.push(file);

    callback();
  });
}


gulp.task('clean_public', () => {
  gulp.src(
      '../public/**/*',
      {
        read: false
      }
    )
    .pipe(
      clean({
        force: true
      })
    )
});

gulp.task('concat_vendor', () => {
  gulp.src(
      [
        'bower_components/jquery/dist/jquery.js',
        'bower_components/underscore/underscore.js',
        'bower_components/backbone/backbone.js',
        'bower_components/jade/jade.js',
        'bower_components/jade/runtime.js'
      ]
    )
    .pipe(concat('vendor.js'))
    .pipe(gulp.dest('tmp'))
});

gulp.task('compile_babel', () => {
  gulp.src('../src/**/*.js')
    .pipe(babel())
    .pipe(concat('babel.js'))
    .pipe(gulp.dest('tmp'))
});

gulp.task('compile_jade', () => {
  gulp.src('../src/**/*.jade')
    .pipe(
      jade({
        client: true
      })
    )
    .pipe(modifyJade())
    .pipe(concat('jade.js'))
    .pipe(gulp.dest('tmp'))
});

gulp.task('concat_js', () => {
  gulp.src(
      ['tmp/vendor.js', 'tmp/jade.js', 'tmp/babel.js']
    )
    .pipe(concat('main.js'))
    .pipe(gulp.dest('../public'))
});

gulp.task('compile_stylus', function () {
  gulp.src('../src/**/*.styl')
    .pipe(stylus())
    .pipe(concat('main.css'))
    .pipe(gulp.dest('tmp/css'))
});

gulp.task('move_img', () => {
  gulp.src(
      ['../src/**/*.png', '../src/**/*.jpg']
    )
    .pipe(
      rename({
        dirname: ''
      })
    )
    .pipe(gulp.dest('../public/img'))
});

gulp.task('clean_tmp', () => {
  gulp.src(
      'tmp',
      {
        read: false
      }
    )
    .pipe(clean())
});


gulp.task('build', () => {
  runRequence(
    'clean_public',
    ['concat_vendor', 'compile_babel', 'compile_jade', 'compile_stylus', 'move_img'],
    'concat_js',
    'clean_tmp'
  )
});


gulp.task('run', function () {
  runRequence('build', () => {
    gulp.watch(
      ['../src/**/*.js', '../src/**/*.jade', '../src/**/*.stylus'],
      ['build']
    )
  })
});