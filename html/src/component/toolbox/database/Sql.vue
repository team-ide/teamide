<template>
  <div class="toolbox-database-sql">
    <template v-if="ready">
      <div>{{ sqlData }}</div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxType", "toolbox", "option", "wrap", "sqlData"],
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
.toolbox-database-sql {
  width: 100%;
  height: 100%;
}
</style>
