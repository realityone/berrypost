const path = require('path');

module.exports = {
    mode: 'production',
    entry: {
        index: {
            import: './src/index/index.js',
            dependOn: 'vendor'
        },
        invoke: {
            import: './src/invoke/invoke.js',
            dependOn: 'vendor'
        },
        config: {
            import: './src/config/config.js',
            dependOn: 'vendor'
        },
        sidebar: {
            import: './src/sidebar/sidebar.js',
            dependOn: 'vendor'
        },
        dashboard: {
            import: './src/dashboard/dashboard.js',
            dependOn: 'vendor'
        },
        blueprint: {
            import: './src/blueprint/blueprint.js',
            dependOn: 'vendor'
        },
        login: {
            import: './src/login/login.js',
            dependOn: 'vendor'
        },
        vendor: './src/vendor.js',
    },
    output: {
        filename: '[name].bundle.js',
        path: path.resolve(__dirname, 'dist'),
        clean: true,
    },
    module: {
        rules: [{
            test: /\.css$/i,
            use: ['style-loader', 'css-loader'],
        },
        {
            test: /\.(png|svg|jpg|jpeg|gif)$/i,
            type: 'asset/resource',
        },
        {
            test: /\.(woff|woff2|eot|ttf|otf)$/i,
            type: 'asset/resource',
        }],
    },
};
