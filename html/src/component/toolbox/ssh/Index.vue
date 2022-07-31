<template>
  <div class="toolbox-ssh-editor-box">
    <template v-if="ready">
      <template v-if="extend && extend.isFTP">
        <FTP
          ref="ftp"
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </FTP>
      </template>
      <template v-else>
        <SSH
          ref="ssh"
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
          :initToken="initToken"
          :initSocket="initSocket"
        >
        </SSH>
      </template>
    </template>
  </div>
</template>


<script>
import FTP from "./FTP";
import SSH from "./SSH";

export default {
  components: { FTP, SSH },
  props: ["source", "extend", "toolboxWorker"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.ready = true;
    },
    onFocus() {
      this.$children.forEach((one) => {
        one.onFocus && one.onFocus();
      });
    },
    refresh() {
      this.$children.forEach((one) => {
        one.refresh && one.refresh();
      });
    },
    async initToken(obj) {
      if (this.tool.isEmpty(obj.token)) {
        let param = {};
        let res = await this.toolboxWorker.work("createToken", param);
        res.data = res.data || {};
        obj.token = res.data.token;
      }
      obj.tokenWork = async (work, param) => {
        param = param || {};
        param.token = obj.token;
        let res = await this.toolboxWorker.work(work, param);
        return res;
      };
      return obj.token;
    },
    initSocket(obj) {
      if (obj.socket != null) {
        obj.socket.close();
      }

      obj.writeData = (data) => {
        obj.socket.send(data);
      };
      obj.writeMessage = (message) => {
        obj.socket.send(this.source.sshTeamIDEMessage + message);
      };
      obj.writeEvent = (event) => {
        obj.socket.send(this.source.sshTeamIDEEvent + event);
      };
      obj.writeError = (error) => {
        obj.socket.send(this.source.sshTeamIDEError + error);
      };

      let url = this.source.api;
      url = url.substring(url.indexOf(":"));
      if (obj.isFTP) {
        url = "ws" + url + "api/toolbox/ssh/ftp";
      } else {
        url = "ws" + url + "api/toolbox/ssh/shell";
      }
      url += "?token=" + encodeURIComponent(obj.token);
      url += "&jwt=" + encodeURIComponent(obj.tool.getJWT());
      obj.socket = new WebSocket(url);
      obj.socket.binaryType = "arraybuffer";
      obj.socket.onopen = () => {
        obj.onEvent("socket open");
      };
      obj.socket.onmessage = (event) => {
        let message = event.data;
        // if (message instanceof ArrayBuffer) {
        //   try {
        //     let data = new Uint8Array(message);
        //     if (this.tool.isUtf8(data)) {
        //       message = this.tool.Utf8ArrayToStr(data);
        //     }
        //   } catch (e) {
        //     message = event.data;
        //   }
        // }
        if (typeof message == "string") {
          if (message.indexOf(this.source.sshTeamIDEEvent) == 0) {
            obj.onEvent &&
              obj.onEvent(
                message.substring(this.source.sshTeamIDEEvent.length)
              );
          } else if (message.indexOf(this.source.sshTeamIDEError) == 0) {
            obj.onError &&
              obj.onError(
                message.substring(this.source.sshTeamIDEError.length)
              );
          } else if (message.indexOf(this.source.sshTeamIDEMessage) == 0) {
            obj.onMessage &&
              obj.onMessage(
                message.substring(this.source.sshTeamIDEMessage.length)
              );
          } else if (message.indexOf(this.source.sshTeamIDEAlert) == 0) {
            obj.onAlert &&
              obj.onAlert(
                message.substring(this.source.sshTeamIDEAlert.length)
              );
          } else if (message.indexOf(this.source.sshTeamIDEConsole) == 0) {
            obj.onConsole &&
              obj.onConsole(
                message.substring(this.source.sshTeamIDEConsole.length)
              );
          } else if (message.indexOf(this.source.sshTeamIDEStdout) == 0) {
            obj.onStdout &&
              obj.onStdout(
                message.substring(this.source.sshTeamIDEStdout.length)
              );
          } else {
            obj.onData && obj.onData(message);
          }
        } else {
          obj.onData && obj.onData(message);
        }
      };
      obj.socket.onclose = () => {
        obj.onEvent("socket close");
        obj.socket = null;
      };
      obj.socket.onerror = () => {
        console.log("socket error");
      };
      return obj.socket;
    },
    destroy() {
      if (this.socket != null) {
        this.socket.close();
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeUpdate() {},
  beforeDestroy() {
    this.destroy();
  },
};
</script>

<style>
.toolbox-ssh-editor-box {
  width: 100%;
  height: 100%;
}
</style>
