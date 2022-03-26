<template>
  <li v-if="ready">
    <template v-if="tool.isNotEmpty(field.text)">
      <span class="text">{{ field.text }}<span class="pdlr-5">:</span></span>
    </template>

    <template v-if="field.fields == null">
      <ModelEditorFieldInput
        :source="source"
        :context="context"
        :bean="bean"
        :field="field"
        :wrap="wrap"
      ></ModelEditorFieldInput>
    </template>

    <template v-if="tool.isNotEmpty(field.comment)">
      <span class="comment">{{ field.comment }}</span>
    </template>
  </li>
</template>


<script>
export default {
  components: {},
  props: ["source", "context", "wrap", "bean", "field"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
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
