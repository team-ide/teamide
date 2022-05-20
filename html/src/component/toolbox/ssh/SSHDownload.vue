<template>
  <div></div>
</template>

<script>
import Zmodem from "zmodem.js";
export default {
  components: {},
  props: ["source", "wrap", "token"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    show(zsession, term, callback) {
      this.zsession = zsession;
      this.term = term;
      this.callback = callback;
      this.downloadFile();
    },
    downloadFile() {
      let that = this;
      that.zsession.on("offer", function (xfer) {
        function on_form_submit() {
          if (xfer.get_details().size > 2048 * 1024 * 1024) {
            xfer.skip();
            that.tool.warn(`${xfer.get_details().name} 超过 2048 MB, 无法下载`);
            return;
          }
          let FILE_BUFFER = [];
          xfer.on("input", (payload) => {
            that.updateProgress(xfer);
            FILE_BUFFER.push(new Uint8Array(payload));
          });

          xfer.accept().then(() => {
            that.saveFile(xfer, FILE_BUFFER);
            that.term.write("\r\n");
          }, console.error.bind(console));
        }
        on_form_submit();
      });
      let promise = new Promise((res) => {
        that.zsession.on("session_end", () => {
          res();
        });
      });
      that.zsession.start();
      return promise;
    },
    bytesHuman(bytes, precision) {
      if (!/^([-+])?|(\.\d+)(\d+(\.\d+)?|(\d+\.)|Infinity)$/.test(bytes)) {
        return "-";
      }
      if (bytes === 0) return "0";
      if (typeof precision === "undefined") precision = 1;
      const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "BB"];
      const num = Math.floor(Math.log(bytes) / Math.log(1024));
      const value = (bytes / Math.pow(1024, Math.floor(num))).toFixed(
        precision
      );
      return `${value} ${units[num]}`;
    },
    updateProgress(xfer) {
      let detail = xfer.get_details();
      let name = detail.name;
      let total = detail.size;
      let offset = xfer.get_offset();
      let percent;
      if (total === 0 || total === offset) {
        percent = 100;
      } else {
        percent = Math.round((offset / total) * 100);
      }
      this.term.write(
        "\r" +
          "下载文件" +
          name +
          " " +
          this.bytesHuman(offset) +
          " " +
          this.bytesHuman(total) +
          " " +
          percent +
          "% "
      );
    },
    saveFile(xfer, buffer) {
      let name = xfer.get_details().name;

      let blob = new Blob(buffer, { type: "application/octet-stream" });
      let downloadElement = document.createElement("a");
      let href = window.URL.createObjectURL(blob); //创建下载的链接
      downloadElement.href = href;
      downloadElement.download = name; //下载后文件名
      document.body.appendChild(downloadElement);
      downloadElement.click(); //点击下载
      document.body.removeChild(downloadElement); //下载完成移除元素
      window.URL.revokeObjectURL(href); //释放blob对象
      this.callback && this.callback();
    },
  },
  created() {},
  mounted() {
    this.wrap.showSSHDownload = this.show;
    this.wrap.hideSSHDownload = this.show;
  },
};
</script>

<style>
</style>
