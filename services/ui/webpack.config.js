const path = require('path');
const webpack = require('webpack');
const CopyWebpackPlugin = require("copy-webpack-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

const distFolder = "www"

module.exports = {
    // entryPoint.js is the main file of your application
    // from there all required parts of the application are imported
    // wepack will start to traverse imports starting from this file
    entry: {
        main: "./src/index.tsx",
    },
    resolve: {
        // Add `.ts` and `.tsx` as a resolvable extension.
        extensions: [".ts", ".tsx", ".js"]
    },
    module: {
        rules: [
            // all files with a `.ts` or `.tsx` extension will be handled by `ts-loader`
            { 
                test: /\.tsx?$/, 
                loader: "ts-loader",
                exclude: /node_modules/
            }
        ]
    },
    plugins: [
        // clean files in webpack_dist before doing anything
        new CleanWebpackPlugin({
            cleanOnceBeforeBuildPatterns: [ './' + distFolder ]
        }),
        new CopyWebpackPlugin({
            patterns: [
              { from: path.resolve(__dirname, "uiconfig.local.json"), to: "config.json"},
            ],
          }),
        // we provide es5 and es6 shim as plugins via webpack
        // you need to extend the entry point for this to work, see below
        // new webpack.ProvidePlugin({
        //     es5: 'es5-shim',
        //     es6: 'es6-shim'
        // }),
        new HtmlWebpackPlugin({
            title: 'Kubelab',
            template: "src/index.html",
            inject: false
        }),
        new webpack.HotModuleReplacementPlugin()
    ],
    devServer: {
        // contentBase: path.resolve(__dirname, distFolder),
        hot: true,
        serveIndex: true,
        historyApiFallback: {
          index: '/ui'
        },
        publicPath: "/ui"
    },
    devtool: 'inline-source-map',
    mode: "development",
    output: {
        libraryTarget: "umd",
        filename: 'kubelab.[name].bundle.js',
        path: __dirname + '/' + distFolder
    }
};