<template>
  <div class="toolbox-zookeeper-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout :width="style.left.width">
          <tm-layout height="50px">
            <div class="pdlr-10 pdt-5">
              <div class="tm-btn tm-btn-xs bg-grey-6" @click="refresh">
                刷新
              </div>
              <div
                class="tm-btn tm-btn-xs bg-teal-8"
                @click="toInsert({ path: '/' })"
              >
                新建
              </div>
              <div
                class="tm-btn tm-btn-xs bg-grey"
                @click="toolboxWorker.showInfo()"
              >
                信息
              </div>
              <div class="color-orange ft-12 pdt-5">
                双击展开目录，单击查看节点数据
              </div>
            </div>
          </tm-layout>
          <tm-layout height="auto" class="app-scroll-bar">
            <div class="pd-10">
              <el-input placeholder="输入关键字进行过滤" v-model="filterText">
              </el-input>
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
                :filter-node-method="filterNode"
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
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <div class="pd-10">
            <el-form ref="form" size="mini" @submit.native.prevent>
              <div class="ft-12 color-grey" v-if="stat != null">
                节点信息：
                <div class="tm-link color-green" @click="reloadForm()">
                  刷新
                </div>
                <div>
                  <div class="pdtb-2">
                    创建时间:
                    <span class="mglr-5">
                      {{ tool.formatDateByTime(stat.ctime) }}
                    </span>
                    创建事务ID:
                    <span class="mglr-5">
                      {{ stat.czxid }}
                    </span>
                  </div>
                  <div class="pdtb-2">
                    修改时间:
                    <span class="mglr-5">
                      {{ tool.formatDateByTime(stat.mtime) }}
                    </span>
                    修改事务ID:
                    <span class="mglr-5">
                      {{ stat.mzxid }}
                    </span>
                  </div>
                  <div class="pdtb-2" v-if="stat.ephemeralOwner != null">
                    临时节点:
                    <span class="mglr-5 color-green"> 是 </span>
                    所有者SessionID:
                    <span class="mglr-5">
                      {{ stat.ephemeralOwner }}
                    </span>
                  </div>
                  <div class="pdtb-2" v-if="stat.dataVersion != null">
                    更改次数:
                    <span class="mglr-5">
                      {{ stat.dataVersion }}
                    </span>
                  </div>
                  <div class="pdtb-2" v-if="stat.numChildren != null">
                    子节点数量:
                    <span class="mglr-5">
                      {{ stat.numChildren }}
                    </span>
                  </div>
                </div>
              </div>
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
            </el-form>
            <div class="pdtb-20">
              <div
                class="tm-btn bg-grey-7"
                @click="toolboxWorker.showJSONData(form.value)"
              >
                预览值
              </div>
              <div class="tm-btn bg-teal-8" @click="toSave">保存</div>
            </div>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
    <ShowInfo :source="source" :toolboxWorker="toolboxWorker"></ShowInfo>
  </div>
</template>


<script>
import ShowInfo from "./ShowInfo.vue";

var JSONbig = require("json-bigint");
var JSONbigString = JSONbig({});

export default {
  components: { ShowInfo },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      style: {
        left: {
          width: "600px",
        },
        main: {},
      },
      searchForm: {
        pattern: null,
      },
      form: {
        dir: null,
        name: null,
        nameFormat: null,
        path: null,
        value: null,
      },
      stat: null,
      ready: false,
      expands: [],
      defaultProps: {
        children: "children",
        label: "name",
        isLeaf: "leaf",
      },
      filterText: "",
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
        if (this.tool.isJSONString(value)) {
          try {
            let data = JSONbigString.parse(value);
            this.form.nameFormat = JSON.stringify(data, null, "    ");
          } catch (e) {}
        }
      }
    },
    filterText(val) {
      this.$refs.tree.filter(val);
    },
  },
  methods: {
    filterNode(value, data) {
      if (!value) return true;
      return (
        data.name && data.name.toLowerCase().indexOf(value.toLowerCase()) !== -1
      );
    },
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
        let list = res.data || [];
        let dataList = [];
        list.forEach((one) => {
          let oneData = Object.assign({}, one);
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
      if (!data.hasChildren) {
        data.leaf = true;
      }
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
      this.stat = null;
      this.form.value = null;
    },
    reloadForm() {
      if (this.stat != null) {
        this.toUpdate({ path: this.form.path });
      } else {
        this.toInsert({ path: this.form.path });
      }
    },
    async toUpdate(one) {
      this.tool.stopEvent();
      let data = await this.get(one.path);
      if (data == null) {
        data = {};
      }
      this.form.path = one.path;
      this.stat = data.stat;
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
    async loadChildren(path) {
      let param = this.toolboxWorker.getWorkParam({
        path: path,
      });
      let res = await this.server.zookeeper.getChildren(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      return res;
    },
    async get(path) {
      let param = this.toolboxWorker.getWorkParam({
        path: path,
      });
      let res = await this.server.zookeeper.get(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      return res.data;
    },
    async doSave() {
      let param = this.toolboxWorker.getWorkParam({
        path: this.form.path,
        data: this.form.value,
      });

      if (this.tool.isEmpty(param.path)) {
        this.tool.error("路径不能为空！");
        return;
      }
      if (param.path.indexOf("//") >= 0 || param.path.endsWith("/")) {
        this.tool.error("路径格式有误，请检查路径格式！");
        return;
      }

      let res = await this.server.zookeeper.save(param);
      if (res.code == 0) {
        this.tool.success("保存成功!");

        let path = param.path;
        if (path.lastIndexOf("/") < path.length - 1) {
          let key = path.substring(0, path.lastIndexOf("/"));
          if (!key.startsWith("/")) {
            key = "/" + key;
          }
          this.reloadChildren(key);
          this.toUpdate({
            path: path,
          });
        }
      } else {
        this.tool.error(res.msg);
      }
    },
    async doDelete(path) {
      let param = this.toolboxWorker.getWorkParam({
        path: path,
      });
      let res = await this.server.zookeeper.delete(param);
      if (res.code == 0) {
        this.tool.success("删除成功!");

        if (path.lastIndexOf("/") < path.length - 1) {
          let key = path.substring(0, path.lastIndexOf("/"));
          if (!key.startsWith("/")) {
            key = "/" + key;
          }
          this.reloadChildren(key);
        }
      } else {
        this.tool.error(res.msg);
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
  user-select: text;
}
.toolbox-zookeeper-editor .el-tree {
  user-select: text;
}
</style>
