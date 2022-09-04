<template>
  <div class="toolbox-file-manager-editor">
    <tm-layout height="100%">
      <tm-layout height="auto">
        <template v-for="one in list">
          <tm-layout ref="layout" :key="`layout-${one.id}`" :width="one.width">
            <FileManager
              :source="source"
              :toolboxWorker="toolboxWorker"
              :place="one.place"
              :placeId="one.placeId"
              :changeDir="changeDir"
            ></FileManager>
          </tm-layout>
          <template v-if="one.hasBar">
            <tm-layout-bar :key="`layout-bar-${one.id}`" right></tm-layout-bar>
          </template>
        </template>
      </tm-layout>
      <tm-layout-bar top></tm-layout-bar>
      <tm-layout height="200px">
        <Progress :source="source" :toolboxWorker="toolboxWorker"></Progress>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
import FileManager from "./FileManager.vue";
import Progress from "./Progress.vue";

export default {
  components: { FileManager, Progress },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      list: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.list.splice(0, this.list.length);
      this.$nextTick(() => {
        if (this.extend) {
          this.addOne(this.extend.place, this.extend.placeId);
        }
      });
    },
    addOne(place, placeId) {
      let list = [];
      if (this.list) {
        this.list.forEach((one) => {
          let data = Object.assign({}, one);
          list.push(data);
        });
      }
      list.push({
        id: "ID-" + this.tool.getNumber(),
        place: place,
        placeId: placeId,
        width: "100%",
        hasBar: false,
      });

      if (list.length > 1) {
        list.forEach((one, index) => {
          if (index == list.length - 1) {
            one.width = "auto";
          } else {
            one.hasBar = true;
            one.width = `${100 / list.length}%`;
          }
        });
      }
      this.list = list;
      this.$nextTick(() => {
        if (this.$refs["layout"] && this.$refs["layout"].forEach) {
          this.$refs["layout"].forEach((one) => {
            one.initSize && one.initSize();
          });
        }
      });
    },
    changeDir() {},
    refresh() {
      console.log(this.list);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-file-manager-editor {
  width: 100%;
  height: 100%;
}
</style>
