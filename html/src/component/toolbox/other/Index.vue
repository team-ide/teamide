<template>
  <div class="toolbox-other-editor">
    <template v-if="ready">
      <template v-if="extend.type == 'format'">
        <Format
          :source="source"
          :toolbox="toolbox"
          :extend="extend"
          :wrap="wrap"
        ></Format>
      </template>
      <template v-else-if="extend.type == 'generateDataFile'">
        <GenerateDataFile
          :source="source"
          :toolbox="toolbox"
          :extend="extend"
          :wrap="wrap"
        ></GenerateDataFile>
      </template>
    </template>
  </div>
</template>


<script>
import Format from "./Format";
import GenerateDataFile from "./GenerateDataFile";

export default {
  components: { Format, GenerateDataFile },
  props: ["source", "toolboxType", "toolbox", "extend", "wrap"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
    },
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
