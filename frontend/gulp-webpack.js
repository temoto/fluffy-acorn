var gulp = require('gulp');
var webpack = require('webpack');
var gutil = require("gulp-util")


const logger = (err, stats) => {
  //var statColor = stats.compilation.warnings.length < 1 ? 'green' : 'yellow';

  var stringStats = stats.toString({});
  if(stringStats.trim().length > 0) {
    gutil.log(gutil.colors.cyan("Webpack Stats: ") + "\n", stringStats);
  }

  if(err) throw new gutil.PluginError("webpack", err);

  if(stats.compilation.errors.length > 0) {

    // stats.compilation.errors.forEach(function(error){
    //  handleErrors(error)
    //  statColor = 'red'
    // })

  } else {
    var compileTime = (stats.endTime - stats.startTime)
    //gutil.log(gutil.colors[statColor](stats))
    gutil.log('Compiled with', gutil.colors.cyan('webpack:dev'), 'in', gutil.colors.magenta(compileTime))
  }
}

gulp.task('webpack:dev', (callback) => {
  const config = require('./webpack.config')('dev');
  webpack(config, (err, stats) => {
    logger(err, stats);
    if(stats.compilation.errors.length > 0) {
      callback(err);
      return;
    }
    callback();
  })
});

gulp.task('webpack:prod', (callback) => {
  const config = require('./webpack.config')('production');
  webpack(config, (err, stats) => {
    logger(err, stats);
    if(stats.compilation.errors.length > 0) {
      callback(err);
      return;
    }
    callback();
  });
});
