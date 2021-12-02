module.exports = {
  assetsDir: process.env.assetsDir,
  publicPath: process.env.publicPath,
  devServer: {
    port: 8081,
    proxy: {
      "/server": {
        target: "http://127.0.0.1:8081/server",
        changeOrigin: true
      },
    }
  },
  productionSourceMap: false,
  transpileDependencies: [
    "vuetify"
  ],
}