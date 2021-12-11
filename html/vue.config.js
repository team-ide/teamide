module.exports = {
  assetsDir: "static",
  publicPath: "",
  devServer: {
    port: 21081,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:21080/teamide/",
        changeOrigin: true
      },
    }
  },
  productionSourceMap: false,
  transpileDependencies: [
    "vuetify"
  ],
}