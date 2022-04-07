<template>
  <div class="toolbox-ssh-editor-box">
    <template v-if="ready">
      <template v-if="extend && extend.isFTP">
        <ToolboxFTP
          ref="ftp"
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :wrap="wrap"
          :extend="extend"
          :token="token"
          :socket="socket"
        >
        </ToolboxFTP>
      </template>
      <template v-else>
        <ToolboxSSH
          ref="ssh"
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :wrap="wrap"
          :extend="extend"
          :token="token"
          :socket="socket"
        >
        </ToolboxSSH>
      </template>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: [
    "source",
    "data",
    "toolboxType",
    "toolbox",
    "option",
    "wrap",
    "extend",
  ],
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
      this.wrap.writeError = this.writeError;
      this.wrap.writeEvent = this.writeEvent;
      await this.initToken();
      this.initSocket();
      this.ready = true;
    },
    async initToken() {
      if (this.tool.isEmpty(this.token)) {
        let param = {};
        let res = await this.wrap.work("createToken", param);
        res.data = res.data || {};
        this.token = res.data.token;
      }
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
      url = "ws" + url + "ws/toolbox/ssh/connection?token=" + this.token;
      this.socket = new WebSocket(url);

      this.socket.onopen = () => {
        this.onEvent("socket open");
      };
      this.socket.onmessage = (event) => {
        // 接收推送的消息
        let data = event.data.toString();
        if (data.indexOf(this.TeamIDEEvent) == 0) {
          this.onEvent(data.substring(this.TeamIDEEvent.length));
        } else if (data.indexOf(this.TeamIDEError) == 0) {
          this.onError(data.substring(this.TeamIDEError.length));
        } else if (data.indexOf(this.TeamIDEMessage) == 0) {
          this.onMessage(data.substring(this.TeamIDEMessage.length));
        } else {
          this.onData(data);
        }
      };
      this.socket.onclose = () => {
        this.onEvent("socket close");
      };
      this.socket.onerror = () => {
        console.log("socket error");
      };
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    if (this.socket != null) {
      this.socket.close();
    }
  },
};
</script>

<style>
.toolbox-ssh-editor-box {
  width: 100%;
  height: 100%;
}
</style>
