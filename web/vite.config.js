export default {
  server: {
    port: 3000,
    strictPort: true,
    proxy: {
      "/workflow": {
        target: "http://127.0.0.1:2323",
        changeOrigin: true
      }
    }
  }
}
