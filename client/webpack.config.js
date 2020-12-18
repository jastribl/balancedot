module.exports = {
    entry: "./src/index.js",
    output: {
        path: __dirname + "/public",
        publicPath: "/",
        filename: "bundle.js",
    },
    resolve: {
        extensions: [".jsx", ".js", "json", "css"],
    },
    mode: 'development',
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                loader: require.resolve("babel-loader"),
                exclude: /node_modules/,
                // Options for the plugin
                options: {
                    presets: [require.resolve("@babel/preset-react")],
                },
            },
        ],
    },
}
