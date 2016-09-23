var browserify   = require('browserify');
var changed      = require('gulp-changed');
var coffeeify    = require('coffeeify');
var connect      = require('gulp-connect');
var del          = require('del');
var fs           = require('fs');
var gulp         = require('gulp');
var gutil        = require('gulp-util');
var jade         = require('gulp-jade');
var jadeify      = require('jadeify');
var notify       = require('gulp-notify');
var path         = require('path');
var plumber      = require('gulp-plumber');
var prettyHrtime = require('pretty-hrtime');
var ractiveify   = require('ractiveify');
var source       = require('vinyl-source-stream');
var stylus       = require('gulp-stylus');
var watchify     = require('watchify');

const srcJade = ['src/**/*.jade', '!src/component/**/*.jade'];

ractiveify.extensions.push('ractive');

gulp.task('default', ['build']);

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

gulp.task('browserify', () => {
  var plugins = [];
  if (global.isWatching) {
    plugins.push(watchify);
  }
  var bundler = browserify({
    debug: true,
    entries: ['src/index.coffee'],
    extensions: ['.coffee', '.js'],
    cache: {},
    packageCache: {},
    plugins: plugins,
  });
  bundler.transform(jadeify, {compileDebug: true, pretty: true});
  bundler.transform(coffeeify);
  bundler.transform(ractiveify);

  var bundle = () => {
    bundleLogger.start();

    return bundler
      .bundle()
      .on('error', handleErrors)
      .pipe(source('bundle.js'))
      .pipe(gulp.dest('build/'))
      .on('end', bundleLogger.end);
  };

  if (global.isWatching) {
    bundler.on('update', bundle);
  }

  return bundle();
});

gulp.task('build', ['clean', 'browserify', 'copy', 'jade', 'stylus']);

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
    .pipe(jade())
    .pipe(gulp.dest('build/'));
});

gulp.task('stylus', () => {
  return gulp.src(['src/**/*.styl'])
    .pipe(stylus())
    .pipe(gulp.dest('build/'));
});

gulp.task('setWatch', () => {
  global.isWatching = true;
});

gulp.task('watch', ['clean', 'setWatch', 'build'], () => {
  gulp.watch(['src/**/*.coffee', 'src/**/*.js', 'src/component/**/*'], ['browserify']);
  gulp.watch(['src/**/*.html', 'src/**/*.png'], ['copy']);
  gulp.watch(['src/**/*.styl'], ['stylus']);
  gulp.watch(srcJade, ['jade']);
});
