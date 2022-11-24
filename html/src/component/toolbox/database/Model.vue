<template>
  <div class="toolbox-database-ddl" v-loading="loading">
    <template v-if="ready">
      <el-form class="pdt-10 pdlr-10" size="mini" inline>
        <el-form-item label="类型">
          <el-select v-model="formData.modelType" style="width: 100px">
            <el-option
              v-for="(one, index) in modelTypes"
              :key="index"
              :value="one.value"
              :label="one.text"
              :disabled="one.disabled"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div style="height: calc(100% - 140px) !important">
        <Editor
          ref="Editor"
          :source="source"
          :value="showDDL"
          language="go"
        ></Editor>
      </div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "actived", "ownerName", "tableName"],
  data() {
    return {
      ready: false,
      showDDL: null,
      loading: false,
      modelTypes: [{ text: "Table", value: "table" }],
      formData: {
        modelType: "table",
      },
    };
  },
  computed: {},
  watch: {
    "formData.modelType"() {
      this.toLoad();
    },
  },
  methods: {
    onFocus() {
      if (this.inited) {
        return;
      }
      this.$nextTick(async () => {
        this.init();
      });
    },
    init() {
      this.inited = true;
      this.ready = true;
      this.toLoad();
    },
    async toLoad() {
      this.loading = true;
      await this.load();
      this.loading = false;
    },
    async load() {
      let data = await this.loadModel(this.ownerName, this.tableName);
      this.$refs.Editor.setValue(data.content);
    },
    async loadModel(ownerName, tableName) {
      let param = Object.assign({}, this.formData);
      param.ownerName = ownerName;
      param.tableName = tableName;
      let res = await this.toolboxWorker.work("model", param);
      res.data = res.data || {};
      return res.data;
    },
  },
  created() {},
  mounted() {
    if (this.actived) {
      this.init();
    }
  },
};
</script>

<style>
.toolbox-database-ddl {
  width: 100%;
  height: 100%;
}
</style>
