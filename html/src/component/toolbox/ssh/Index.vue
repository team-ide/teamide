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
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-ssh-editor-box {
  width: 100%;
  height: 100%;
}
</style>
