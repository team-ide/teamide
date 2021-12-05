module.exports = {
  assetsDir: "static",
  publicPath: "./",
  devServer: {
    port: 21081,
    proxy: {
      "/server": {
        target: "http://127.0.0.1:21080/server",
        changeOrigin: true
      },
    }
  },
  productionSourceMap: false,
  transpileDependencies: [
    "vuetify"
  ],
}