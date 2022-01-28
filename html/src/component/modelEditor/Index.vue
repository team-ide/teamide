<template>
  <div
    class="model-editor"
    tabindex="-1"
    @keydown="keydown"
    v-if="group != null && group.fields != null"
  >
    <template v-if="data != null">
      <div class="model-editor-toolbar">
        <b-button-group size="sm">
          <b-button title="保存（Ctrl+s）" @click="save()">
            <b-icon icon="sticky-fill"></b-icon>
          </b-button>
          <b-button title="复制" @click="copy()">
            <b-icon icon="stickies-fill"></b-icon>
          </b-button>
          <b-button title="刷新重置" @click="reload()">
            <b-icon icon="arrow-clockwise"></b-icon>
          </b-button>
          <b-button title="上一步" :disabled="true" @click="previousStep()">
            <b-icon icon="arrow-left"></b-icon>
          </b-button>
          <b-button title="下一步" :disabled="true" @click="nextStep()">
            <b-icon icon="arrow-right"></b-icon>
          </b-button>
          <b-button
            v-if="group.isAction"
            title="测试"
            :disabled="true"
            @click="toTest()"
          >
            <b-icon icon="skip-end-fill"></b-icon>
          </b-button>
          <b-button title="帮助说明" :disabled="true" @click="help()">
            <b-icon icon="exclamation-circle-fill"></b-icon>
          </b-button>
          <b-button title="删除" :disabled="true" @click="toDelete()">
            <b-icon icon="x-circle"></b-icon>
          </b-button>
        </b-button-group>
      </div>
      <template v-if="group.isAction">
        <ModelEditorAction
          class="mgt-40"
          :source="source"
          :context="context"
          :bean="data"
          :wrap="wrap"
        >
        </ModelEditorAction>
      </template>
      <template v-else>
        <ul v-if="!group.isAction" class="mgt-40">
          <template v-for="(one, index) in group.fields">
            <ModelEditorField
              :key="'field-' + index"
              :source="source"
              :context="context"
              :field="one"
              :bean="data"
              :wrap="wrap"
            >
            </ModelEditorField>
            <template v-if="one.fields != null && !one.isList">
              <ModelEditorFieldBean
                :key="'field-bean-' + index"
                :source="source"
                :context="context"
                :field="one"
                :bean="data"
                :wrap="wrap"
              >
              </ModelEditorFieldBean>
            </template>
            <template v-if="one.fields != null && one.isList">
              <ModelEditorFieldList
                :key="'field-list-' + index"
                :source="source"
                :context="context"
                :field="one"
                :bean="data"
                :wrap="wrap"
              >
              </ModelEditorFieldList>
            </template>
          </template>
        </ul>
      </template>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "context", "model", "group", "application"],
  data() {
    return {
      key: this.tool.getNumber(),
      wrap: {},
      data: null,
      ready: false,
    };
  },
  computed: {},
  watch: {
    model(newModel, oldModel) {
      this.initData();
    },
  },
  methods: {
    init() {
      this.set(JSON.parse(JSON.stringify(this.model || {})));
    },
    initData() {
      this.set(JSON.parse(JSON.stringify(this.data || {})));
    },
    get() {
      return this.data;
    },
    set(data) {
      this.data = data;
    },
    refresh() {
      this.initData();
    },
    onChange(bean, name, value) {
      bean = bean || {};
      bean[name] = value;
      this.toCallChange();
    },
    push(bean, name, value) {
      bean = bean || {};
      bean[name] = bean[name] || [];
      value = value || {};
      bean[name].push(value);
      this.toCallChange();
    },
    del(bean, name, value) {
      if (bean != null && bean[name] != null) {
        let index = bean[name].indexOf(value);
        if (index >= 0) {
          bean[name].splice(index, 1);
          this.toCallChange();
        }
      }
    },
    up(bean, name, value) {
      if (bean != null && bean[name] != null) {
        let index = bean[name].indexOf(value);
        if (index > 0) {
          bean[name].splice(index, 1);
          bean[name].splice(index - 1, 0, value);
          this.toCallChange();
        }
      }
    },
    down(bean, name, value) {
      if (bean != null && bean[name] != null) {
        let index = bean[name].indexOf(value);
        if (index >= 0 && index < bean[name].length - 1) {
          bean[name].splice(index, 1);
          bean[name].splice(index + 1, 0, value);
          this.toCallChange();
        }
      }
    },
    toCallChange() {
      let change_model = JSON.parse(JSON.stringify(this.data));
      // console.log(this.data);
      this.last_change_model = change_model;
      this.$emit("change", this.group, change_model);
    },
    keydown(event) {
      //ctrl+s
      if (event.keyCode === 83 && event.ctrlKey) {
        event.preventDefault && event.preventDefault();
        event.stopPropagation && event.stopPropagation();
        this.save();
        return false;
      }
    },
    save() {
      let model = JSON.parse(JSON.stringify(this.data));
      this.$emit("save", this.group, model);
    },
    copy() {
      let model = JSON.parse(JSON.stringify(this.data));
      this.application.showModelForm(this.group, model, (group, param) => {
        Object.assign(model, param);
        return this.application.saveModel(group, model, true);
      });
    },
    reload() {
      let model = JSON.parse(JSON.stringify(this.model));
      this.set(model);
      this.$emit("change", this.group, model);
    },
    toDelete() {},
    previousStep() {},
    nextStep() {},
    help() {},
    toTest() {},
  },
  created() {
    this.wrap.push = this.push;
    this.wrap.del = this.del;
    this.wrap.up = this.up;
    this.wrap.down = this.down;
    this.wrap.onChange = this.onChange;
    this.wrap.refresh = this.refresh;
  },
  mounted() {
    this.init();
  },
};
</script>

<style>
.model-editor {
  width: 100%;
  height: 100%;
  overflow: auto;
}
.model-editor ul {
  margin-top: 10px;
}
.model-editor ul,
.model-editor li {
  list-style: none;
  padding: 0px;
  font-size: 12px;
}
.model-editor li {
  display: block;
  line-height: 22px;
  margin-bottom: 3px;
}
.model-editor .text,
.model-editor .input,
.model-editor .comment {
  padding: 0px 5px;
}
.model-editor-toolbar {
  width: 100%;
  text-align: center;
  padding: 5px 0px;
  position: absolute;
}

.model-editor-toolbar .btn {
  border-radius: 0px;
}
</style>
