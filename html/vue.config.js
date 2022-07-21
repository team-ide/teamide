module.exports = {
  assetsDir: process.env.assetsDir,
  publicPath: process.env.publicPath,
  devServer: {
    port: 21081,
    proxy: {
      "/api": {
        ws: true,
        target: "http://127.0.0.1:21080/teamide/",
        changeOrigin: true
      },
    }
  },
  productionSourceMap: false,
  transpileDependencies: [
    "vuetify"
  ],
  runtimeCompiler: true,
}