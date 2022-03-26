<template>
  <div v-if="ready">
    <template v-for="(one, index) in bean.steps">
      <ModelEditorStep
        :key="'step-' + index"
        :source="source"
        :context="context"
        :wrap="wrap"
        :bean="one"
        :parent="bean"
      >
      </ModelEditorStep>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "context", "wrap", "bean"],
  data() {
    return {
      ready: false,
      webApiOpen: false,
    };
  },
  computed: {},
  watch: {
    bean() {
      this.init();
    },
  },
  methods: {
    init() {
      if (this.bean.steps == null) {
        this.bean.steps = [];
        this.wrap.refresh();
        return;
      }
      this.ready = true;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>
