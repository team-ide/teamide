<template>
  <div class="toolbox-zookeeper-editor">
    <template v-if="ready">
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
              @node-expand="nodeExpand"
              @node-collapse="nodeCollapse"
            >
              <span
                class="toolbox-editor-tree-span"
                slot-scope="{ node, data }"
              >
                <template v-if="data.path == '/'">
                  <span>/</span>
                </template>
                <template v-else>
                  <span>{{ node.label }}</span>
                </template>
                <div class="toolbox-editor-tree-btn-group">
                  <div
                    class="tm-link color-grey ft-14 mgr-2"
                    @click="toReloadChildren(data)"
                  >
                    <i class="mdi mdi-reload"></i>
                  </div>
                  <div
                    class="tm-link color-blue ft-16 mgr-2"
                    @click="toInsert(data)"
                  >
                    <i class="mdi mdi-plus"></i>
                  </div>
                  <div
                    class="tm-link color-orange ft-15 mgr-2"
                    @click="toDelete(data)"
                  >
                    <i class="mdi mdi-delete-outline"></i>
                  </div>
                </div>
              </span>
            </el-tree>
          </div>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <div class="pd-10">
            <el-form ref="form" size="mini" @submit.native.prevent>
              <el-form-item label="目录">
                <el-input v-model="form.dir" @change="dirChange"> </el-input>
              </el-form-item>
              <el-form-item label="名称">
                <el-input v-model="form.name" @change="nameChange"> </el-input>
              </el-form-item>
              <template v-if="form.nameFormat != null">
                <el-form-item label="名称解码">
                  <el-input
                    type="textarea"
                    v-model="form.nameFormat"
                    :autosize="{ minRows: 5, maxRows: 10 }"
                  >
                  </el-input>
                </el-form-item>
              </template>
              <el-form-item label="Path">
                <el-input v-model="form.path" @change="pathChange"> </el-input>
              </el-form-item>
              <el-form-item label="Value">
                <el-input
                  type="textarea"
                  v-model="form.value"
                  :autosize="{ minRows: 5, maxRows: 10 }"
                >
                </el-input>
              </el-form-item>
              <template v-if="form.valueJson != null">
                <el-form-item label="值JSON预览">
                  <el-input
                    type="textarea"
                    v-model="form.valueJson"
                    :autosize="{ minRows: 5, maxRows: 10 }"
                  >
                  </el-input>
                </el-form-item>
              </template>
            </el-form>
            <div class="pdtb-20">
              <div class="tm-btn bg-teal-8" @click="toSave">保存</div>
            </div>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      style: {
        left: {
          width: "600px",
        },
        main: {},
      },
      form: {
        dir: null,
        name: null,
        nameFormat: null,
        path: null,
        value: null,
        valueJson: null,
      },
      ready: false,
      expands: [],
      defaultProps: {
        children: "children",
        label: "name",
        isLeaf: "leaf",
      },
    };
  },
  computed: {},
  watch: {
    "form.path"() {
      this.pathChange();
    },
    "form.name"(value) {
      this.form.nameFormat = null;
      if (this.tool.isNotEmpty(value)) {
        try {
          let nameFormat = decodeURIComponent(value);
          if (nameFormat != value) {
            this.form.nameFormat = nameFormat;
            value = nameFormat;
          }
        } catch (e) {}
        try {
          if (
            (value.startsWith("{") && value.endsWith("}")) ||
            (value.startsWith("[") && value.endsWith("]"))
          ) {
            let data = JSON.parse(value);
            this.form.nameFormat = JSON.stringify(data, null, "    ");
          }
        } catch (e) {}
      }
    },
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
    dirChange() {
      this.form.dir = this.form.dir || "";
      if (!this.form.dir.endsWith("/")) {
        this.form.dir = this.form.dir + "/";
      }
      this.isFromWatchDir = true;
      this.form.path = this.form.dir + this.form.name;
    },
    nameChange() {
      this.isFromWatchName = true;
      this.form.path = this.form.dir + this.form.name;
    },
    pathChange() {
      if (this.isFromWatchDir || this.isFromWatchName) {
        delete this.isFromWatchDir;
        delete this.isFromWatchName;
        return;
      }
      this.form.dir = "/";
      this.form.name = "";
      if (this.tool.isEmpty(this.form.path)) {
        return;
      }
      let index = this.form.path.lastIndexOf("/");
      if (index >= 0) {
        let dir = this.form.path.substring(0, index);
        if (!dir.endsWith("/")) {
          dir = dir + "/";
        }
        let name = "";
        if (index < this.form.path.length - 1) {
          name = this.form.path.substring(index + 1);
        }
        this.form.dir = dir;
        this.form.name = name;
      } else {
        this.form.name = this.form.path;
      }
    },
    init() {
      if (this.extend && this.extend.expands) {
        this.expands = this.extend.expands;
      }
      this.ready = true;
    },
    nodeExpand(data) {
      let index = this.expands.indexOf(data.key);
      if (index < 0) {
        this.expands.push(data.key);
        this.toolboxWorker.updateExtend({
          expands: this.expands,
        });
      }
    },
    nodeCollapse(data) {
      console.log(data);
      let needDeletes = [];
      this.expands.forEach((one) => {
        if (data.key == "/") {
          needDeletes.push(one);
        } else if (one == data.key || one.startsWith(data.key + "/")) {
          needDeletes.push(one);
        }
      });
      if (needDeletes.length > 0) {
        needDeletes.forEach((one) => {
          let index = this.expands.indexOf(one);
          if (index >= 0) {
            this.expands.splice(index, 1);
          }
        });
        this.toolboxWorker.updateExtend({
          expands: this.expands,
        });
      }
    },
    refresh() {
      this.reloadChildren(this.$refs.tree.root);
    },
    toReloadChildren(data) {
      this.tool.stopEvent();
      this.reloadChildren(data);
    },
    nodeClick(data, node) {
      this.rowClickTimeCache = this.rowClickTimeCache || {};
      let nowTime = new Date().getTime();
      let clickTime = this.rowClickTimeCache[node];
      this.rowClickTimeCache[node] = nowTime;
      if (clickTime) {
        let timeout = nowTime - clickTime;
        if (timeout < 300) {
          delete this.rowClickTimeCache[node];
          this.nodeDbClick(node);
        }
      }
    },
    nodeDbClick(node) {
      if (node.expanded) {
        node.expanded = false;
      } else {
        node.loaded = false;
        node.expand();
      }
    },
    reloadChildren(key) {
      this.tool.stopEvent();
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
        let dataList = [];
        list.forEach((one) => {
          let name = one.name;
          let oneData = { name: name };
          dataList.push(oneData);
        });
        this.formatDataList(parent, dataList);
        resolve(dataList);
        this.initTreeWidth();
      }
    },
    initTreeWidth() {},
    formatDataList(parent, dataList) {
      dataList = dataList || [];
      dataList.forEach((data) => {
        this.formatData(parent, data);
      });
    },
    formatData(parent, data) {
      data.path = "/" + data.name;
      if (parent != null && parent.path != "/") {
        data.path = parent.path + "/" + data.name;
      }
      data.key = data.path;
      data.leaf = false;

      this.loadHasChildren(data.path).then((res) => {
        if (res.code == 0) {
          if (!res.data.hasChildren) {
            data.leaf = true;
            let node = this.$refs.tree.getNode(data);
            if (node != null) {
              node.isLeaf = true;
            }
          }
        }
      });
    },
    currentChange(data) {
      this.toUpdate(data);
    },
    toSave() {
      this.doSave();
    },
    toInsert(one) {
      this.tool.stopEvent();
      if (one.path.endsWith("/")) {
        this.form.path = one.path + "xxx";
      } else {
        this.form.path = one.path + "/xxx";
      }
      this.form.value = null;
    },
    async toUpdate(one) {
      this.tool.stopEvent();
      let data = await this.get(one.path);
      if (data == null) {
        data = {};
      }
      this.form.path = one.path;
      this.form.value = data.data;
    },
    toDelete(data) {
      this.tool.stopEvent();
      this.tool
        .confirm(
          "将删除[" + data.path + "]节点和该节点下所有子节点，确认删除？"
        )
        .then(async () => {
          this.doDelete(data.path);
        })
        .catch((e) => {});
    },
    async loadHasChildren(path) {
      let param = {
        path: path,
      };
      let res = await this.toolboxWorker.work("hasChildren", param);
      return res;
    },
    async loadChildren(path) {
      let param = {
        path: path,
      };
      let res = await this.toolboxWorker.work("getChildren", param);
      return res;
    },
    async get(path) {
      let param = {
        path: path,
      };
      let res = await this.toolboxWorker.work("get", param);
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

      let res = await this.toolboxWorker.work("save", param);
      if (res.code == 0) {
        this.tool.success("保存成功!");

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
      let res = await this.toolboxWorker.work("delete", param);
      if (res.code == 0) {
        this.tool.success("删除成功!");

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
</style>
