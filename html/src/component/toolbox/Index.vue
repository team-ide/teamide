<template>
  <div class="toolbox-editor" tabindex="-1" v-if="toolboxType != null">
    {{ data }}
    <template v-if="option != null">
      {{ option }}
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox"],
  data() {
    return {
      key: this.tool.getNumber(),
      option: null,
      ready: false,
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
      this.initOption();
    },
    initOption() {
      let option = null;
      if (this.tool.isNotEmpty(this.data.option)) {
        option = JSON.parse(this.data.option);
      }
      this.set(option);
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
};
</script>

<style>
.toolbox-editor {
  width: 100%;
  height: 100%;
  overflow: auto;
}
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
}
.toolbox-editor .text {
  display: inline-block;
  min-width: 80px;
}
.toolbox-editor .text,
.toolbox-editor .input,
.toolbox-editor .comment {
  padding: 0px 5px;
}
</style>
