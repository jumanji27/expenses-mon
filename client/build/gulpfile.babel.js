import gulp from 'gulp';
import runRequence from 'run-sequence';
import babel from 'gulp-babel';
import browserify from 'browserify';
import babelify from 'babelify';
import source from 'vinyl-source-stream';
import jade from 'gulp-jade';
import stylus from 'gulp-stylus';
import concat from 'gulp-concat';
import del from 'del';
import rename from 'gulp-rename';
import through from 'through2';
import path from 'path';

let modifyJade = () => {
  return through.obj((file, enc, callback) => { // Was crashed without return
    let context = gulp.src('../src/**/*.jade');

    if (!file.isBuffer()) {
      context.push(file);
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
    context.push(file);

    callback();
  });
};


gulp.task('clean_public', () => {
  return del(
    '../public/**/*',
    {
      force: true
    }
  )
});

gulp.task('concat_vendor', () => {
  return gulp.src(
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
  return browserify({
    entries: ['../src/router.js', '../src/main.js'],
    debug: true
  })
  .transform(babelify)
  .bundle()
  .pipe(source('babel.js'))
  .pipe(gulp.dest('tmp'))
});

gulp.task('compile_jade', () => {
  return gulp.src('../src/**/*.jade')
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
  return gulp.src(
      ['tmp/vendor.js', 'tmp/jade.js', 'tmp/babel.js']
    )
    .pipe(concat('main.js'))
    .pipe(gulp.dest('../public'))
});

gulp.task('compile_stylus', function () {
  return gulp.src('../src/**/*.styl')
    .pipe(stylus())
    .pipe(concat('main.css'))
    .pipe(gulp.dest('../public'))
});

gulp.task('move_img', () => {
  return gulp.src(
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
  return del('tmp')
});


gulp.task('build', () => {
  return runRequence(
    'clean_public',
    ['concat_vendor', 'compile_babel', 'compile_jade', 'compile_stylus', 'move_img'],
    'concat_js',
    'clean_tmp'
  )
});


gulp.task('run', function () {
  return runRequence('build', () => {
    gulp.watch(
      ['../src/**/*.js', '../src/**/*.jade', '../src/**/*.stylus'],
      ['build']
    )
  })
});