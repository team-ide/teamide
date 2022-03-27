<template>
  <div class="toolbox-zookeeper-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="100%">
          <tm-layout :width="style.left.width" class="scrollbar">
            <div class="pd-10">
              <el-tree
                ref="tree"
                :load="loadNode"
                lazy
                :props="defaultProps"
                :default-expanded-keys="expands"
                node-key="key"
                @node-click="nodeClick"
                @current-change="currentChange"
                :expand-on-click-node="false"
              >
                <span class="toolbox-box-tree-span" slot-scope="{ node, data }">
                  <template v-if="data.path == '/'">
                    <span>/</span>
                  </template>
                  <template v-else>
                    <span>{{ node.label }}</span>
                  </template>
                  <div class="toolbox-box-tree-btn-group">
                    <a
                      class="tm-link color-grey ft-14 mgr-2"
                      @click="toReloadChildren(data)"
                    >
                      <i class="mdi mdi-reload"></i>
                    </a>
                    <a
                      class="tm-link color-blue ft-16 mgr-2"
                      @click="toInsert(data)"
                    >
                      <i class="mdi mdi-plus"></i>
                    </a>
                    <a
                      class="tm-link color-orange ft-15 mgr-2"
                      @click="toDelete(data)"
                    >
                      <i class="mdi mdi-delete-outline"></i>
                    </a>
                  </div>
                </span>
              </el-tree>
            </div>
          </tm-layout>
          <tm-layout-bar right></tm-layout-bar>
          <tm-layout width="auto">
            <b-form class="pd-10">
              <b-form-group label="Path" label-size="sm">
                <b-form-input size="sm" v-model="form.path"> </b-form-input>
              </b-form-group>
              <b-form-group label="Value" label-size="sm">
                <b-form-textarea
                  size="sm"
                  rows="5"
                  max-rows="10"
                  v-model="form.value"
                >
                </b-form-textarea>
              </b-form-group>
              <template v-if="form.valueJson != null">
                <b-form-group label="值JSON预览" label-size="sm">
                  <b-form-textarea
                    size="sm"
                    rows="5"
                    max-rows="10"
                    v-model="form.valueJson"
                  >
                  </b-form-textarea>
                </b-form-group>
              </template>
              <div class="pdtb-20">
                <div class="tm-btn bg-teal-8" @click="toSave">保存</div>
              </div>
            </b-form>
          </tm-layout>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      style: {
        left: {
          width: "800px",
        },
        main: {},
      },
      form: {
        path: null,
        value: null,
        valueJson: null,
      },
      ready: false,
      expands: [],
      opens: [],
      defaultProps: {
        children: "children",
        label: "name",
        isLeaf: "leaf",
      },
    };
  },
  computed: {},
  watch: {
    "form.value"(value) {
      this.form.valueJson = null;
      if (this.tool.isNotEmpty(value)) {
        try {
          if (
            (value.startsWith("{") && value.endsWith("}")) ||
            (value.startsWith("[") && value.endsWith("]"))
          ) {
            let data = JSON.parse(value);
            this.form.valueJson = JSON.stringify(data, null, "    ");
          }
        } catch (e) {
          this.form.valueJson = e;
        }
      }
    },
  },
  methods: {
    init() {
      this.ready = true;
    },
    toReloadChildren(data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      this.reloadChildren(data);
    },
    reloadChildren(key) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      let node = this.$refs.tree.getNode(key);
      if (node) {
        node.loaded = false;
        node.expand();
      }
    },
    async loadNode(node, resolve) {
      if (node.level === 0) {
        resolve([{ name: "/", path: "/", key: "/" }]);
        this.$nextTick(() => {
          let node = this.$refs.tree.getNode("/");
          node.expand();
        });
        return;
      }
      let parent = node.data;
      let path = parent.path;
      let res = await this.loadChildren(path);
      if (res.code != 0) {
        this.tool.error(res.msg);
        resolve([]);
        this.initTreeWidth();
      } else {
        let list = res.data.children || [];
        let datas = [];
        list.forEach((name) => {
          datas.push({ name: name });
        });
        this.formatDatas(parent, datas);
        resolve(datas);
        this.initTreeWidth();
      }
    },
    initTreeWidth() {},
    formatDatas(parent, datas) {
      datas = datas || [];
      datas.forEach((data) => {
        this.formatData(parent, data);
      });
    },
    formatData(parent, data) {
      data.path = "/" + data.name;
      if (parent != null && parent.path != "/") {
        data.path = parent.path + "/" + data.name;
      }
      data.key = data.path;
    },
    nodeClick() {},
    currentChange(data) {
      this.toUpdate(data);
    },
    toSave() {
      this.doSave();
    },
    toInsert(one) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      if (one.path.endsWith("/")) {
        this.form.path = one.path + "xxx";
      } else {
        this.form.path = one.path + "/xxx";
      }
      this.form.value = null;
    },
    async toUpdate(one) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      let data = await this.get(one.path);
      if (data == null) {
        data = {};
      }
      this.form.path = one.path;
      this.form.value = data.data;
    },
    toDelete(data) {
      if (window.event) {
        window.event.stopPropagation && window.event.stopPropagation();
      }
      this.tool
        .confirm(
          "将删除[" + data.path + "]节点和该节点下所有子节点，确认删除？"
        )
        .then(async () => {
          this.doDelete(data.path);
        })
        .catch((e) => {});
    },
    async loadChildren(path) {
      let param = {
        path: path,
      };
      let res = await this.wrap.work("getChildren", param);
      return res;
    },
    async get(path) {
      let param = {
        path: path,
      };
      let res = await this.wrap.work("get", param);
      return res.data;
    },
    async doSave() {
      let param = {
        path: this.form.path,
        data: this.form.value,
      };

      if (this.tool.isEmpty(param.path)) {
        this.tool.error("路径不能为空！");
        return;
      }
      if (param.path.indexOf("//") >= 0 || param.path.endsWith("/")) {
        this.tool.error("路径格式有误，请检查路径格式！");
        return;
      }

      let res = await this.wrap.work("save", param);
      if (res.code == 0) {
        this.tool.info("保存成功!");

        let path = param.path;
        if (path.lastIndexOf("/") < path.length - 1) {
          let key = path.substring(0, path.lastIndexOf("/"));
          if (!key.startsWith("/")) {
            key = "/" + key;
          }
          this.reloadChildren(key);
        }
      }
    },
    async doDelete(path) {
      let param = {
        path: path,
      };
      let res = await this.wrap.work("delete", param);
      if (res.code == 0) {
        this.tool.info("删除成功!");

        if (path.lastIndexOf("/") < path.length - 1) {
          let key = path.substring(0, path.lastIndexOf("/"));
          if (!key.startsWith("/")) {
            key = "/" + key;
          }
          this.reloadChildren(key);
        }
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-zookeeper-editor {
  width: 100%;
  height: 100%;
}
.toolbox-zookeeper-editor .el-tree {
  /* border: 1px solid #f3f3f3; */
  border-bottom: 0px;
}
.toolbox-zookeeper-editor .el-tree-node__content {
  border-bottom: 1px dotted #696969;
}

.toolbox-box .el-tree {
  background-color: transparent;
  color: unset;
  width: auto;
  font-size: 12px;
}
.toolbox-box .el-tree .mdi {
  vertical-align: middle;
}
.toolbox-box .el-tree-node__content {
  position: relative;
}

.toolbox-box .el-tree-node__content .toolbox-box-tree-btn-group {
  display: none;
}
.toolbox-box .el-tree-node__content:hover .toolbox-box-tree-btn-group {
  display: block;
}
.toolbox-box .toolbox-box-tree-span .toolbox-box-tree-btn-group {
  position: absolute;
  right: 5px;
  top: -3px;
}
.toolbox-box .toolbox-box-tree-span {
  /* flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between; */
  font-size: 14px;
}
.toolbox-box .el-tree__empty-block {
  display: none;
}
.toolbox-box .el-tree-node__children {
  overflow: visible !important;
}
.toolbox-box .el-tree-node__children.v-enter-active,
.toolbox-box .el-tree-node__children.v-leave-active {
  overflow: hidden !important;
}
.toolbox-box .el-tree .el-tree-node.is-current > .el-tree-node__content {
  background-color: #636363 !important;
}
.toolbox-box .el-tree .el-tree-node:focus > .el-tree-node__content {
  background-color: #545454;
}
.toolbox-box .el-tree .el-tree-node > .el-tree-node__content:hover {
  background-color: #545454;
}
</style>
