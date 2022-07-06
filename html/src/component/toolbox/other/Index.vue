<template>
  <div class="toolbox-other-editor">
    <template v-if="extend != null">
      <template v-if="extend.type == 'format'">
        <Format
          :source="source"
          :toolboxWorker="toolboxWorker"
          :extend="extend"
        ></Format>
      </template>
      <template v-else-if="extend.type == 'generateDataFile'">
        <GenerateDataFile
          :source="source"
          :toolboxWorker="toolboxWorker"
          :extend="extend"
        ></GenerateDataFile>
      </template>
      <template v-else-if="extend.type == 'fontImage'">
        <FontImage
          :source="source"
          :toolboxWorker="toolboxWorker"
          :extend="extend"
        ></FontImage>
      </template>
    </template>
  </div>
</template>


<script>
import Format from "./Format";
import GenerateDataFile from "./GenerateDataFile";
import FontImage from "./FontImage";

export default {
  components: { Format, GenerateDataFile, FontImage },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {},
    refresh() {
      this.$nextTick(() => {});
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-other-editor {
  width: 100%;
  height: 100%;
}
</style>
