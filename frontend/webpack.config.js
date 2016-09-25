var path = require('path');
var webpack = require('webpack');

var loggingOptions = {
  //Available options: https://github.com/webpack/webpack/blob/master/lib/Stats.js#L26-L40
  assets: true,
  version: true,
  timings: true,
  hash: true,
  chunks: false,
  chunkModules: false,
  errorDetails: true,
  reasons: false,
  colors: true,
};

module.exports = function(env) {
  var webpackConfig = {
    context: path.resolve('src/'),
    resolve: {
      extensions: ['', '.coffee', '.js'],
      root: [path.resolve('.')],
    },
    plugins: [],
    stats: loggingOptions,
    module: {
      loaders: [
        {test: /\.es6\.js$/, loaders: ['babel'], exclude: /node_modules/},
        {test: /\.coffee$/, loaders: ['coffee'], exclude: /node_modules/},
        {test: /\.ract\.jade$/, loaders: ['ractive', 'jade-html'], exclude: /node_modules/},
        // {test: /\.jade$/, loaders: ['jade-html'], exclude: /node_modules/},
        {test: /\.ract$/, loaders: ['ractive'], exclude: /node_modules/},
      ]
    },
  }

  if(env !== 'test') {
    // Karma doesn't need entry points or output settings
    webpackConfig.entry = {
      index: ['src/index.coffee'],
      // component: ['src/component/messages.coffee'],
    }

    webpackConfig.output = {
      path: 'build/',
      filename: env === 'production' ? '[name]-[hash].js' : '[name].js',
      publicPath: './',
    }
  }

  if(env === 'dev') {
    webpackConfig.devtool = 'source-map-inline'
    webpack.debug = true
  }

  if(env === 'production') {
    webpackConfig.plugins.push(
      new webpackManifest('src/', 'build/'),
      new webpack.DefinePlugin({'process.env': {'NODE_ENV': JSON.stringify('production')}}),
      new webpack.optimize.DedupePlugin(),
      new webpack.optimize.UglifyJsPlugin(),
      new webpack.NoErrorsPlugin()
    )
  }

  return webpackConfig
}
