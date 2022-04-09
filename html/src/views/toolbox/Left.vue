<template>
  <div>
    <div class="toolbox-type-box">
      <template v-for="(toolboxType, index) in source.toolboxTypes">
        <div :key="'type-' + index" class="toolbox-type">
          <div
            class="toolbox-type-title"
            :class="{
              'toolbox-type-body-hide':
                toolbox.toolboxTypeOpens.indexOf(toolboxType.name) < 0,
            }"
          >
            <template
              v-if="toolbox.toolboxTypeOpens.indexOf(toolboxType.name) < 0"
            >
              <span
                class="tm-pointer ft-12 color-grey pdlr-5"
                @click="toolboxTypeOpen(toolboxType)"
              >
                <b-icon icon="caret-right"></b-icon>
              </span>
            </template>
            <template v-else>
              <span
                class="tm-pointer ft-12 color-grey pdlr-5"
                @click="toolboxTypeClose(toolboxType)"
              >
                <b-icon icon="caret-down-fill"></b-icon>
              </span>
            </template>
            <div
              style="
                max-width: calc(100% - 40px);
                text-overflow: ellipsis;
                overflow: hidden;
                white-space: nowrap;
                cursor: pointer;
              "
              @dblclick="toolboxTypeOpenOrClose(toolboxType)"
            >
              {{ toolboxType.text }}
            </div>
            <div class="toolbox-btn-type ft-12">
              <span
                class="tm-pointer color-green mgr-5"
                @click="toInsert(toolboxType)"
              >
                <b-icon icon="plus-square"></b-icon>
              </span>
            </div>
          </div>
          <div
            class="toolbox-type-body"
            :class="{
              'toolbox-type-body-hide':
                toolbox.toolboxTypeOpens.indexOf(toolboxType.name) < 0,
            }"
            :style="toolboxTypeStyleObject(toolboxType)"
          >
            <div class="toolbox-one-box">
              <template
                v-if="
                  context[toolboxType.name] == null ||
                  context[toolboxType.name].length == 0
                "
              >
                <div class="text-center pdtb-10">
                  <div class="ft-12 color-grey-7">暂无模型</div>
                </div>
              </template>
              <template v-else>
                <template v-for="(data, index) in context[toolboxType.name]">
                  <div :key="'data-' + index" class="toolbox-one">
                    <div
                      class="toolbox-one-title"
                      @dblclick="dataOpen(toolboxType, data)"
                    >
                      <div
                        style="
                          max-width: calc(100% - 40px);
                          text-overflow: ellipsis;
                          overflow: hidden;
                          white-space: nowrap;
                        "
                      >
                        {{ data.name }}
                        <template v-if="tool.isNotEmpty(data.comment)">
                          <span class="mgl-3 color-grey-6">
                            {{ data.comment }}
                          </span>
                        </template>
                      </div>
                      <div class="toolbox-btn-type ft-12">
                        <span
                          v-if="toolboxType.name == 'ssh'"
                          class="tm-pointer color-blue mgl-5"
                          @click="dataOpenSfpt(toolboxType, data)"
                        >
                          <i class="mdi mdi-folder color-orange ft-13"></i>
                        </span>
                        <span
                          class="tm-pointer color-blue mgl-5"
                          @click="toUpdate(toolboxType, data)"
                        >
                          <b-icon icon="pencil-square" class="ft-13"></b-icon>
                        </span>
                        <span
                          class="tm-pointer color-orange mgl-5"
                          @click="toDelete(toolboxType, data)"
                        >
                          <b-icon icon="x-square"></b-icon>
                        </span>
                      </div>
                    </div>
                  </div>
                </template>
              </template>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "context"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    toolboxTypeStyleObject: function (toolboxType) {
      var opened = this.toolbox.toolboxTypeOpens.indexOf(toolboxType.name) >= 0;
      var height = 0;
      if (this.context[toolboxType.name] != null) {
        height = this.context[toolboxType.name].length * 25;
      }
      if (height == 0) {
        height = 40;
      } else {
        height += 15;
      }
      var marginTop = -height;
      if (opened) {
        marginTop = 0;
      }
      return {
        height: height + "px",
        // marginTop: marginTop + "px",
      };
    },
    toolboxTypeOpenOrClose(toolboxType) {
      if (this.toolbox.toolboxTypeOpens.indexOf(toolboxType.name) < 0) {
        this.toolboxTypeOpen(toolboxType);
      } else {
        this.toolboxTypeClose(toolboxType);
      }
    },
    toolboxTypeOpen(toolboxType) {
      if (this.toolbox.toolboxTypeOpens.indexOf(toolboxType.name) < 0) {
        this.toolbox.toolboxTypeOpens.push(toolboxType.name);
        toolboxType.marginTop = 0;
      }
    },
    toolboxTypeClose(toolboxType) {
      if (this.toolbox.toolboxTypeOpens.indexOf(toolboxType.name) >= 0) {
        this.toolbox.toolboxTypeOpens.splice(
          this.toolbox.toolboxTypeOpens.indexOf(toolboxType.name),
          1
        );
        toolboxType.marginTop = 100;
      }
    },
    dataOpen(toolboxType, data) {
      this.open(data);
    },
    dataOpenSfpt(toolboxType, data) {
      this.open(data, {
        isFTP: true,
      });
    },
    async open(data, extend) {
      let extendStr = null;
      if (extend != null) {
        extendStr = JSON.stringify(extend);
      }
      let param = {
        toolboxId: data.toolboxId,
        extend: extendStr,
      };
      let res = await this.server.toolbox.open(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let openData = res.data.open;
      let tab = await this.openByOpenData(openData);
      if (tab != null) {
        this.toolbox.doActiveTab(tab);
      }
    },
    async openByOpenData(openData) {
      let data = this.getToolboxData(openData.toolboxId);
      if (data == null) {
        await this.closeOpen(openData.openId);
        return;
      }
      let toolboxType = this.getToolboxType(data.toolboxType);
      if (toolboxType == null) {
        await this.closeOpen(openData.openId);
      }
      openData.data = data;
      openData.toolboxType = toolboxType;
      if (this.tool.isNotEmpty(openData.extend)) {
        openData.extend = JSON.parse(openData.extend);
      } else {
        openData.extend = null;
      }
      let tab = this.toolbox.createTabByData(openData);
      this.toolbox.addTab(tab);
      return tab;
    },
    async activeOpen(openId) {
      let param = {
        openId: openId,
      };
      let res = await this.server.toolbox.open(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    async closeOpen(openId) {
      let param = {
        openId: openId,
      };
      let res = await this.server.toolbox.close(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    getToolboxType(type) {
      let res = null;
      this.source.toolboxTypes.forEach((one) => {
        if (one == type || one.name == type || one.name == type.name) {
          res = one;
        }
      });
      return res;
    },
    getToolboxData(data) {
      let res = null;
      for (let type in this.context) {
        if (this.context[type] == null) {
          continue;
        }
        this.context[type].forEach((one) => {
          if (
            one == data ||
            one.toolboxId == data ||
            one.toolboxId == data.toolboxId
          ) {
            res = one;
          }
        });
      }
      return res;
    },
    toInsert(toolboxType) {
      let data = {};
      this.toolbox.showToolboxForm(toolboxType, data, (g, m) => {
        let flag = this.doInsert(g, m);
        return flag;
      });
    },
    toUpdate(toolboxType, data) {
      this.updateData = data;
      this.toolbox.showToolboxForm(toolboxType, data, (g, m) => {
        let flag = this.doUpdate(g, m);
        return flag;
      });
    },
    toDelete(toolboxType, data) {
      this.tool
        .confirm(
          "删除[" +
            toolboxType.text +
            "]工具[" +
            data.name +
            "]将无法回复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(toolboxType, data);
        })
        .catch((e) => {});
    },
    async doDelete(toolboxType, data) {
      let res = await this.server.toolbox.delete(data);
      if (res.code == 0) {
        this.toolbox.initContext();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdate(toolboxType, data) {
      data.toolboxType = toolboxType.name;
      data.toolboxId = this.updateData.toolboxId;
      let res = await this.server.toolbox.update(data);
      if (res.code == 0) {
        this.toolbox.initContext();
        let tab = this.toolbox.getTabByData(data);
        if (tab != null) {
          Object.assign(tab.data, data);
        }
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doInsert(toolboxType, data) {
      data.toolboxType = toolboxType.name;
      let res = await this.server.toolbox.insert(data);
      if (res.code == 0) {
        this.toolbox.initContext();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async initOpens() {
      let opens = await this.loadOpens();

      await opens.forEach(async (openData) => {
        await this.openByOpenData(openData);
      });

      // 激活最后
      let activeOpenData = null;
      opens.forEach(async (openData) => {
        if (activeOpenData == null) {
          activeOpenData = openData;
        } else {
          if (
            new Date(openData.openTime).getTime() >
            new Date(activeOpenData.openTime).getTime()
          ) {
            activeOpenData = openData;
          }
        }
      });
      if (activeOpenData != null) {
        this.toolbox.doActiveTab(activeOpenData.openId);
      }
    },
    async loadOpens() {
      let param = {};
      let res = await this.server.toolbox.queryOpens(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let opens = res.data.opens || [];
      return opens;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.initOpens();
    this.toolbox.closeOpen = this.closeOpen;
    this.toolbox.activeOpen = this.activeOpen;
    this.toolbox.saveToolbox = this.saveToolbox;
  },
};
</script>

<style>
.toolbox-type-box {
  width: 100%;
  position: relative;
  user-select: none;
}
.toolbox-type {
  line-height: 20px;
  font-size: 12px;
  font-weight: 400;
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
.toolbox-type-title {
  line-height: 28px;
  font-size: 14px;
  font-weight: 500;
  margin: 0px 0px;
  color: #8c8c8c;
  position: relative;
  z-index: 1;
  border-bottom: 1px dotted #4e4e4e;
  display: flex;
  height: 30px;
}
.toolbox-type-title.toolbox-type-body-hide {
  /* height: 0px !important; */
  border-bottom-color: transparent !important;
  /* margin-top: -50px; */
}
.toolbox-type-body {
  margin: 0px 5px;
  height: 50px;
  transition: all 0.3s;
  position: relative;
  z-index: 0;
  overflow: hidden;
}
.toolbox-type-body.toolbox-type-body-hide {
  height: 0px !important;
  /* border-top-color: transparent !important; */
  /* margin-top: -50px; */
}
.toolbox-btn-type {
  text-align: right;
  flex: 1;
  display: none;
}
.toolbox-type-title:hover .toolbox-btn-type {
  display: block;
}

.toolbox-one-title {
  line-height: 23px;
  font-size: 12px;
  font-weight: 400;
  margin: 0px 0px 0px 15px;
  color: #cdcdcd;
  position: relative;
  display: flex;
  border-bottom: 1px dotted #4e4e4e;
  height: 25px;
  cursor: pointer;
}
.toolbox-one-title:hover .toolbox-btn-type {
  display: block;
}
</style>
