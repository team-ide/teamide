<template>
  <div class="toolbox-ssh-editor-box">
    <template v-if="ready">
      <template v-if="extend && extend.isFTP">
        <FTP
          ref="ftp"
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
          :token="token"
          :socket="socket"
        >
        </FTP>
      </template>
      <template v-else>
        <SSH
          ref="ssh"
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
          :token="token"
          :socket="socket"
        >
        </SSH>
      </template>
    </template>
    <FileEdit :source="source" :toolbox="toolbox" :wrap="wrap"> </FileEdit>
  </div>
</template>


<script>
import FTP from "./FTP";
import SSH from "./SSH";
import FileEdit from "./FileEdit";

export default {
  components: { FTP, SSH, FileEdit },
  props: ["source", "toolboxType", "toolbox", "option", "extend", "wrap"],
  data() {
    return {
      ready: false,
      token: null,
      socket: null,
      TeamIDEEvent: "TeamIDE:event:",
      TeamIDEMessage: "TeamIDE:message:",
      TeamIDEError: "TeamIDE:error:",
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.wrap.writeData = this.writeData;
      this.wrap.writeMessage = this.writeMessage;
      this.wrap.writeEvent = this.writeEvent;
      this.wrap.writeError = this.writeError;
      this.wrap.tokenWork = this.tokenWork;

      await this.initToken();
      this.initSocket();
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
    async initToken() {
      if (this.tool.isEmpty(this.token)) {
        let param = {};
        let res = await this.wrap.work("createToken", param);
        res.data = res.data || {};
        this.token = res.data.token;
      }
    },
    async tokenWork(work, param) {
      param = param || {};
      param.token = this.token;
      let res = await this.wrap.work(work, param);
      return res;
    },
    onEvent(arg) {
      this.$refs.ftp && this.$refs.ftp.onEvent && this.$refs.ftp.onEvent(arg);
      this.$refs.ssh && this.$refs.ssh.onEvent && this.$refs.ssh.onEvent(arg);
    },
    onError(arg) {
      this.$refs.ftp && this.$refs.ftp.onError && this.$refs.ftp.onError(arg);
      this.$refs.ssh && this.$refs.ssh.onError && this.$refs.ssh.onError(arg);
    },
    onMessage(arg) {
      this.$refs.ftp &&
        this.$refs.ftp.onMessage &&
        this.$refs.ftp.onMessage(arg);
      this.$refs.ssh &&
        this.$refs.ssh.onMessage &&
        this.$refs.ssh.onMessage(arg);
    },
    onData(arg) {
      this.$refs.ftp && this.$refs.ftp.onData && this.$refs.ftp.onData(arg);
      this.$refs.ssh && this.$refs.ssh.onData && this.$refs.ssh.onData(arg);
    },
    writeData(data) {
      this.socket.send(data);
    },
    writeMessage(message) {
      this.socket.send(this.TeamIDEMessage + message);
    },
    writeEvent(event) {
      this.socket.send(this.TeamIDEEvent + event);
    },
    writeError(error) {
      this.socket.send(this.TeamIDEError + error);
    },
    initSocket() {
      if (this.socket != null) {
        this.socket.close();
      }

      let url = this.source.api;
      url = url.substring(url.indexOf(":"));
      if (this.extend && this.extend.isFTP) {
        url = "ws" + url + "api/toolbox/ssh/ftp";
      } else {
        url = "ws" + url + "api/toolbox/ssh/shell";
      }
      url += "?token=" + encodeURIComponent(this.token);
      url += "&jwt=" + encodeURIComponent(this.tool.getJWT());
      this.socket = new WebSocket(url);
      this.socket.binaryType = "arraybuffer";
      this.socket.onopen = () => {
        this.onEvent("socket open");
      };
      this.socket.onmessage = (event) => {
        if (typeof event.data == "string") {
          let data = event.data;
          if (data.indexOf(this.TeamIDEEvent) == 0) {
            this.onEvent(data.substring(this.TeamIDEEvent.length));
          } else if (data.indexOf(this.TeamIDEError) == 0) {
            this.onError(data.substring(this.TeamIDEError.length));
          } else if (data.indexOf(this.TeamIDEMessage) == 0) {
            this.onMessage(data.substring(this.TeamIDEMessage.length));
          } else {
            this.onData(data);
          }
        } else {
          this.onData(event.data);
        }
      };
      this.socket.onclose = () => {
        this.onEvent("socket close");
        this.socket = null;
      };
      this.socket.onerror = () => {
        console.log("socket error");
      };
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
