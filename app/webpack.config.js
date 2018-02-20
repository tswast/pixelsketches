const path = require('path')

module.exports = {
  entry: './src/index.js',
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist')
  },
  // Don't use in production.
  // https://webpack.js.org/guides/development/
  devtool: 'inline-source-map',
  devServer: {
    contentBase: './dist'
  }
}
