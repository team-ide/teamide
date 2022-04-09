<template>
  <div class="toolbox-editor" v-if="toolboxType != null">
    <template v-if="ready">
      <template v-if="toolboxType.name == 'redis'">
        <ToolboxRedisEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxRedisEditor>
      </template>
      <template v-else-if="toolboxType.name == 'database'">
        <ToolboxDatabaseEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxDatabaseEditor>
      </template>
      <template v-else-if="toolboxType.name == 'zookeeper'">
        <ToolboxZookeeperEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxZookeeperEditor>
      </template>
      <template v-else-if="toolboxType.name == 'elasticsearch'">
        <ToolboxElasticsearchEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxElasticsearchEditor>
      </template>
      <template v-else-if="toolboxType.name == 'kafka'">
        <ToolboxKafkaEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxKafkaEditor>
      </template>
      <template v-else-if="toolboxType.name == 'ssh'">
        <ToolboxSSHEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :data="data"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxSSHEditor>
      </template>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "extend", "toolboxType", "toolbox"],
  data() {
    return {
      key: "toolbox-" + this.tool.getNumber(),
      option: null,
      ready: false,
      wrap: {},
    };
  },
  computed: {},
  watch: {
    data() {
      this.initOption();
    },
  },
  methods: {
    init() {
      this.wrap.work = this.work;
      this.initOption();
      this.ready = true;
    },
    initOption() {
      let option = null;
      if (this.tool.isNotEmpty(this.data.option)) {
        option = JSON.parse(this.data.option);
      }
      this.set(option);
    },
    async work(work, data) {
      let param = {
        toolboxId: this.data.toolboxId,
        work: work,
        data: data,
      };
      let res = await this.server.toolbox.work(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      return res;
    },
    get() {
      return this.option;
    },
    set(option) {
      this.option = option;
    },
    refresh() {
      this.initData();
    },
    reload() {},
  },
  created() {},
  mounted() {
    this.init();
  },
  updated() {},
  beforeDestroy() {
    if (this.wrap.destroy != null) {
      this.wrap.destroy();
    }
  },
};
</script>

<style>
.toolbox-editor {
  width: 100%;
  height: 100%;
  overflow: auto;
}
/* 
.toolbox-editor ul {
  margin-top: 10px;
}
.toolbox-editor ul,
.toolbox-editor li {
  list-style: none;
  padding: 0px;
  font-size: 12px;
}
.toolbox-editor li {
  display: block;
  line-height: 22px;
  margin-bottom: 3px;
} */
.toolbox-editor .text {
  display: inline-block;
  min-width: 80px;
}
.toolbox-editor .text,
.toolbox-editor .input,
.toolbox-editor .comment {
  padding: 0px 5px;
}

.toolbox-editor table {
  padding: 0px 0px;
  width: 100%;
}
.toolbox-editor table thead {
  border: 1px solid #4e4e4e;
}
.toolbox-editor table th {
  text-align: center;
  line-height: 30px;
}
.toolbox-editor table td {
  border-right: 1px solid #4e4e4e;
  border-bottom: 1px solid #4e4e4e;
  padding: 3px 5px;
}
.toolbox-editor table tbody {
  border-left: 1px solid #4e4e4e;
}
.toolbox-editor table td .input {
  padding: 0px 0px;
}
.toolbox-editor table td .model-input {
  min-width: 80px;
}

.part-box {
  line-height: 20px;
  font-size: 12px;
  overflow: auto;
  width: 100%;
  height: 100%;
}
.part-box,
.part-box li {
  padding: 0px;
  margin: 0px;
  list-style: none;
}
.part-box li {
  text-overflow: ellipsis;
  white-space: nowrap;
  word-break: keep-all;
}

.part-box input,
.part-box select {
  color: #ffffff;
  width: 40px;
  min-width: 40px;
  border: 1px dashed transparent;
  background-color: transparent;
  height: 20px;
  max-width: 100%;
  padding: 0px;
  padding-left: 2px;
  padding-right: 2px;
  box-sizing: border-box;
  outline: none;
  font-size: 12px;
}

.part-box input {
  border-bottom: 1px dashed #636363;
}
.part-box select {
  -moz-appearance: auto;
  -webkit-appearance: auto;
}
.part-box option {
  background-color: #ffffff;
  color: #3e3e3e;
}
.part-box input[type="checkbox"] {
  width: 10px;
  min-width: 10px;
  height: 13px;
  vertical-align: -3px;
  margin-left: 6px;
}

.part-box textarea {
  color: #ffffff;
  height: 70px;
  border: 1px dashed #636363;
  text-align: left;
  padding: 5px;
  min-width: 500px;
  background-color: transparent;
  font-size: 12px;
  vertical-align: text-top;
}
</style>
