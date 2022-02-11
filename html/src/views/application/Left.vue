<template>
  <div>
    <div class="context-group-box">
      <template v-for="(modelType, index) in source.modelTypes">
        <div :key="'group-' + index" class="context-group">
          <div
            class="context-group-title"
            :class="{
              'context-group-body-hide':
                application.modelTypeOpens.indexOf(modelType.name) < 0,
            }"
          >
            <template
              v-if="application.modelTypeOpens.indexOf(modelType.name) < 0"
            >
              <span
                class="tm-pointer ft-12 color-grey pdlr-5"
                @click="modelTypeOpen(modelType)"
              >
                <b-icon icon="caret-right"></b-icon>
              </span>
            </template>
            <template v-else>
              <span
                class="tm-pointer ft-12 color-grey pdlr-5"
                @click="modelTypeClose(modelType)"
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
              @dblclick="modelTypeOpenOrClose(modelType)"
            >
              {{ modelType.text }}
            </div>
            <div class="context-btn-group ft-12">
              <span
                class="tm-pointer color-green mgr-5"
                @click="toInsert(modelType)"
              >
                <b-icon icon="plus-square"></b-icon>
              </span>
            </div>
          </div>
          <div
            class="context-group-body"
            :class="{
              'context-group-body-hide':
                application.modelTypeOpens.indexOf(modelType.name) < 0,
            }"
            :style="modelTypeStyleObject(modelType)"
          >
            <div class="context-model-box">
              <template
                v-if="
                  context[modelType.name] == null ||
                  context[modelType.name].length == 0
                "
              >
                <div class="text-center pdtb-10">
                  <div class="ft-12 color-grey-7">暂无模型</div>
                </div>
              </template>
              <template v-else>
                <template v-for="(model, index) in context[modelType.name]">
                  <div :key="'model-' + index" class="context-model">
                    <div
                      class="context-model-title"
                      @dblclick="modelOpen(modelType, model)"
                    >
                      <div
                        style="
                          max-width: calc(100% - 40px);
                          text-overflow: ellipsis;
                          overflow: hidden;
                          white-space: nowrap;
                        "
                      >
                        {{ model.name }}
                        <template v-if="tool.isNotEmpty(model.comment)">
                          <span class="mgl-3 color-grey-6">
                            {{ model.comment }}
                          </span>
                        </template>
                      </div>
                      <div class="context-btn-group ft-12">
                        <span
                          class="tm-pointer color-blue mgl-5"
                          @click="toUpdate(modelType, model)"
                        >
                          <b-icon icon="pencil-square" class="ft-13"></b-icon>
                        </span>
                        <span
                          class="tm-pointer color-orange mgl-5"
                          @click="toDelete(modelType, model)"
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
  props: ["source", "application", "app", "context"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    modelTypeStyleObject: function (modelType, models) {
      var opened = this.application.modelTypeOpens.indexOf(modelType.name) >= 0;
      var height = 0;
      if (this.context[modelType.name] != null) {
        height = this.context[modelType.name].length * 25;
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
    modelTypeOpenOrClose(modelType) {
      if (this.application.modelTypeOpens.indexOf(modelType.name) < 0) {
        this.modelTypeOpen(modelType);
      } else {
        this.modelTypeClose(modelType);
      }
    },
    modelTypeOpen(modelType) {
      if (this.application.modelTypeOpens.indexOf(modelType.name) < 0) {
        this.application.modelTypeOpens.push(modelType.name);
        modelType.marginTop = 0;
      }
    },
    modelTypeClose(modelType) {
      if (this.application.modelTypeOpens.indexOf(modelType.name) >= 0) {
        this.application.modelTypeOpens.splice(
          this.application.modelTypeOpens.indexOf(modelType.name),
          1
        );
        modelType.marginTop = 100;
      }
    },
    modelOpen(modelType, model) {
      let tab = this.application.createTabByModel(modelType, model);
      this.application.addTab(tab);
      this.application.doActiveTab(tab);
    },
    toInsert(modelType) {
      let data = {};
      this.application.showModelForm(modelType, data, (g, m) => {
        return this.doInsert(g, m);
      });
    },
    toUpdate(modelType, model) {
      let data = {
        name: model.name,
        comment: model.comment,
      };
      this.updateData = model;
      this.application.showModelForm(modelType, data, (g, m) => {
        return this.doUpdate(g, m);
      });
    },
    toDelete(modelType, model) {
      this.tool
        .confirm(
          "删除[" +
            modelType.text +
            "]模型[" +
            model.name +
            "]将无法回复，确定删除？"
        )
        .then(() => {
          this.doDelete(modelType, model);
        })
        .catch((e) => {});
    },
    doInsert(modelType, model) {
      let context = Object.assign({}, this.context);
      context[modelType.name] = context[modelType.name] || [];
      let find;
      context[modelType.name].forEach((one) => {
        if (one.name == model.name) {
          find = one;
        }
      });
      if (find != null) {
        this.tool.error("[" + modelType.text + "]模型[" + model.name + "]已存在");
        return false;
      }
      context[modelType.name].push(model);

      let flag = this.doSave(context);

      if (flag) {
        if (this.application.modelTypeOpens.indexOf(modelType.name) < 0) {
          this.modelTypeOpen(modelType);
        }
      }
      return flag;
    },
    doUpdate(modelType, model) {
      let context = Object.assign({}, this.context);
      context[modelType.name] = context[modelType.name] || [];

      let find;
      context[modelType.name].forEach((one) => {
        if (one.name == model.name) {
          find = one;
        }
      });
      if (find != null && find != this.updateData) {
        this.tool.error("[" + modelType.text + "]模型[" + model.name + "]已存在");
        return false;
      }
      if (context[modelType.name].indexOf(this.updateData) < 0) {
        this.tool.error(
          "[" + modelType.text + "]模型[" + this.updateData.name + "]不存在"
        );
        return false;
      }
      Object.assign(this.updateData, model);

      return this.doSave(context);
    },
    doDelete(modelType, model) {
      let context = Object.assign({}, this.context);
      context[modelType.name] = context[modelType.name] || [];
      if (context[modelType.name].indexOf(model) < 0) {
        this.tool.error("[" + modelType.text + "]模型[" + model.name + "]不存在");
        return false;
      }
      context[modelType.name].splice(context[modelType.name].indexOf(model), 1);

      return this.doSave(context);
    },
    async saveModel(modelType, model, isInsert) {
      let context = Object.assign({}, this.context);
      context[modelType.name] = context[modelType.name] || [];

      let find;
      context[modelType.name].forEach((one) => {
        if (one.name == model.name) {
          find = one;
        }
      });
      if (isInsert) {
        if (find != null) {
          this.tool.error("[" + modelType.text + "]模型[" + model.name + "]已存在");
          return false;
        }
        context[modelType.name].push(model);
      } else {
        if (find == null) {
          this.tool.error("[" + modelType.text + "]模型[" + model.name + "]不存在");
          return false;
        }
        context[modelType.name].splice(context[modelType.name].indexOf(find), 1, model);
      }

      return await this.doSave(context);
    },
    async doSave(context) {
      let res = await this.server.application.context.save({
        appName: this.app.name,
        content: JSON.stringify(context),
      });
      if (res.code == 0) {
        this.application.context = res.data;
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.application.saveModel = this.saveModel;
  },
};
</script>

<style>
.context-group-box {
  width: 100%;
  position: relative;
  user-select: none;
}
.context-group {
  line-height: 20px;
  font-size: 12px;
  font-weight: 400;
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
.context-group-title {
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
.context-group-title.context-group-body-hide {
  /* height: 0px !important; */
  border-bottom-color: transparent !important;
  /* margin-top: -50px; */
}
.context-group-body {
  margin: 0px 5px;
  height: 50px;
  transition: all 0.3s;
  position: relative;
  z-index: 0;
  overflow: hidden;
}
.context-group-body.context-group-body-hide {
  height: 0px !important;
  /* border-top-color: transparent !important; */
  /* margin-top: -50px; */
}
.context-btn-group {
  text-align: right;
  flex: 1;
  display: none;
}
.context-group-title:hover .context-btn-group {
  display: block;
}

.context-model-title {
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
.context-model-title:hover .context-btn-group {
  display: block;
}
</style>
