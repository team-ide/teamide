<template>
  <div class="toolbox-ssh-editor">
    <tm-layout height="100%">
      <tm-layout height="auto">
        <FileManager
          :source="source"
          :toolboxWorker="toolboxWorker"
          place="ssh"
          :placeId="`${toolboxWorker.toolboxId}`"
          :openDir="extend == null ? '' : extend.openDir"
          :onChangeOpenDir="onChangeOpenDir"
        ></FileManager>
      </tm-layout>
      <tm-layout-bar top></tm-layout-bar>
      <tm-layout height="200px">
        <Progress :source="source" :toolboxWorker="toolboxWorker"></Progress>
      </tm-layout>
    </tm-layout>
  </div>
</template>

<script>
import FileManager from "../file-manager/FileManager.vue";
import Progress from "../file-manager/Progress.vue";
export default {
  components: { FileManager, Progress },
  props: ["source", "extend", "toolboxWorker"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    async init() {},
    onChangeOpenDir(openDir) {
      let data = this.extend || {};
      if (this.tool.isEmpty(openDir)) {
        openDir = "";
      }
      if (this.tool.isEmpty(openDir) && this.tool.isEmpty(data.openDir)) {
        return;
      }
      if (openDir == data.openDir) {
        return;
      }
      data.openDir = openDir;
      var keyValueMap = {};
      keyValueMap["openDir"] = openDir;
      this.toolboxWorker.updateExtend(keyValueMap);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.isDestroyed = true;
  },
};
</script>

<style>
.toolbox-ftp-editor {
  width: 100%;
  height: 100%;
}
</style>
