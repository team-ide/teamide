<template>
  <div class="toolbox-redis-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout width="500px" class="">
          <Keys
            ref="Keys"
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </Keys>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <Tabs
            :source="source"
            :toolbox="toolbox"
            :toolboxType="toolboxType"
            :wrap="wrap"
          >
          </Tabs>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
import Keys from "./Keys";
import Tabs from "./Tabs";

export default {
  components: { Keys, Tabs },
  props: ["source", "toolboxType", "toolbox", "option", "wrap"],
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
      this.$nextTick(() => {
        this.$refs.Keys && this.$refs.Keys.refresh();
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
.toolbox-redis-editor {
  width: 100%;
  height: 100%;
  user-select: text;
}
</style>
