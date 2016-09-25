var changed      = require('gulp-changed');
var connect      = require('gulp-connect');
var del          = require('del');
var fs           = require('fs');
var gulp         = require('gulp');
var gutil        = require('gulp-util');
var jade         = require('gulp-jade');
var notify       = require('gulp-notify');
var path         = require('path');
var plumber      = require('gulp-plumber');
var prettyHrtime = require('pretty-hrtime');
var source       = require('vinyl-source-stream');
var stylus       = require('gulp-stylus');
require('./gulp-webpack.js');

const srcJade = ['src/**/*.jade', '!src/component/**/*.jade'];

gulp.task('default', ['build:dev']);

var handleErrors = (err) => {
  const msg = err.toString();
  gutil.log(gutil.colors.red(msg));

  notify.onError({
    title: 'Compile Error',
    message: "<%=error.message%>",
  })(err);

  gulp.emit('end');
};

var _gulpsrc = gulp.src;
gulp.src = (globs) => {
  return _gulpsrc(globs)
    .pipe(plumber({errorHandler: handleErrors}));
};

// bundleLogger
var bundleStartTime;
var bundleLogger = {
  start: () => {
    bundleStartTime = process.hrtime();
    gutil.log('Running', gutil.colors.green("'bundle'") + '...');
  },

  end: () => {
    var taskTime = process.hrtime(bundleStartTime);
    var prettyTime = prettyHrtime(taskTime);
    gutil.log('Finished', gutil.colors.green("'bundle'"), 'in', gutil.colors.magenta(prettyTime));
  }
};

gulp.task('build:dev', ['clean', 'copy', 'jade', 'stylus', 'webpack:dev']);
gulp.task('build:prod', ['clean', 'copy', 'jade', 'stylus', 'webpack:prod']);

gulp.task('clean', () => {
  return del(['build/**/*']);
});

gulp.task('connect', () => {
  return connect.server({
    port: 8002,
    root: 'build/',
  });
});

gulp.task('copy', () => {
  return gulp.src(['src/**/*.png', './vendor/**/*'])
  .pipe(gulp.dest('build/'));
});

gulp.task('jade', () => {
  return gulp.src(srcJade)
    .pipe(jade({
      locals: {
        version: 'unknown-build-error',
        title: 'Fluffy Acorn',
        environment: 'dev',
      },
    }))
    .pipe(gulp.dest('build/'));
});

gulp.task('stylus', () => {
  return gulp.src(['src/**/*.styl'])
    .pipe(stylus())
    .pipe(gulp.dest('build/'));
});

gulp.task('watch', ['clean', 'build:dev'], () => {
  gulp.watch(['src/**/*.coffee', 'src/**/*.js', 'src/component/**/*'], ['build:dev']);
  gulp.watch(['src/**/*.html', 'src/**/*.png'], ['copy']);
  gulp.watch(['src/**/*.styl'], ['stylus']);
  gulp.watch(srcJade, ['jade']);
});
