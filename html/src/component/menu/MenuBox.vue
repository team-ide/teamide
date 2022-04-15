<template>
  <div class="menu-box" :class="{ 'menu-sub-left': subLeft }">
    <slot></slot>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["subLeft"],
  data() {
    return {};
  },
  beforeMount() {},
  watch: {},
  methods: {},
  mounted() {},
};
</script>

<style >
.menu-box,
.menu-box ul,
.menu-box li {
  list-style: none;
  padding: 0px;
  margin: 0px;
  user-select: none;
}

.menu-box {
  min-width: 160px;
  padding: 5px 0;
  margin: 2px 0 0;
  list-style: none;
  background-color: #fff;
  border: 1px solid #ccc;
  border: 1px solid rgba(0, 0, 0, 0.2);
  font-family: helvetica neue, Helvetica, Arial, sans-serif;
  font-size: 14px;
  *border-right-width: 2px;
  *border-bottom-width: 2px;
  -webkit-border-radius: 6px;
  -moz-border-radius: 6px;
  border-radius: 6px;
  -webkit-box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
  -moz-box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
  -webkit-background-clip: padding-box;
  -moz-background-clip: padding;
  background-clip: padding-box;
  text-align: left;
  width: 200px;
}
</style>