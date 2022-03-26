<template>
  <div
    class="menu-item"
    :class="{
      divider: divider,
      disabled: disabled,
      header: header,
      'has-sub-box': hasSubBox,
    }"
    @click="onClick"
  >
    <a class :class="aClass" :target="target" :href="href">
      <slot ></slot>
    </a>
    <slot name="MenuSubBox"></slot>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["divider", "disabled", "header", "target", "href", "aClass"],
  data() {
    return {
      hasSubBox: false,
    };
  },
  beforeMount() {},
  watch: {},
  methods: {
    onClick() {
      if (this.disabled) {
        return;
      }
      this.$emit("click");
    },
  },
  mounted() {
    this.$children.forEach((one) => {
      if (one.isMenuSubBox) {
        this.hasSubBox = true;
      }
    });
  },
};
</script>

<style  >
.menu-item.divider {
  *width: 100%;
  height: 2px;
  margin: 9px 1px;
  *margin: -5px 0 5px;
  overflow: hidden;
  background-color: #e5e5e5;
  border-bottom: 1px solid #fff;
}

.menu-item.header {
  display: block;
  padding: 3px 30px 3px 15px;
  font-size: 12px;
  font-weight: 700;
  line-height: 24px;
  color: #999;
  text-shadow: 0 1px 0 rgba(255, 255, 255, 0.5);
  white-space: nowrap;
  text-decoration: none;
  cursor: default;
}
.menu-item {
  position: relative;
}

.menu-item > a {
  display: block;
  padding: 3px 20px;
  clear: both;
  font-weight: 400;
  line-height: 24px;
  color: #333 !important;
  white-space: nowrap;
  text-decoration: none;
  cursor: pointer;
}

.menu-item > a:hover,
.menu-item > a:focus {
  color: #fff !important;
  text-decoration: none;
  background-color: #009688;
  background-image: -moz-linear-gradient(top, #00897b, #00796b);
  background-image: -webkit-gradient(
    linear,
    0 0,
    0 100%,
    from(#00897b),
    to(#00796b)
  );
  background-image: -webkit-linear-gradient(top, #00897b, #00796b);
  background-image: -o-linear-gradient(top, #00897b, #00796b);
  background-image: linear-gradient(to bottom, #00897b, #00796b);
  background-repeat: repeat-x;
  filter: progid:dximagetransform.microsoft.gradient(startColorstr='#ff0088cc',
		endColorstr='#ff0077b3', GradientType=0);
}
.menu-item.disabled {
  cursor: no-drop;
  pointer-events: none;
  opacity: 0.65;
  filter: alpha(opacity = 65);
}
.menu-item.disabled > a,
.menu-item.disabled > a:hover {
  color: #999;
}

.menu-item.disabled > a:hover {
  text-decoration: none;
  cursor: default;
  background-color: transparent;
}

.menu-item:hover > .menu-sub-box {
  display: block;
}

.menu-item.has-sub-box > a {
  padding-right: 30px;
}
.menu-item.has-sub-box > a:after {
  display: block;
  float: right;
  width: 0;
  height: 0;
  margin-top: 5px;
  margin-right: -20px;
  border-color: transparent;
  border-left-color: #ccc;
  border-style: solid;
  border-width: 5px 0 5px 5px;
  content: " ";
}

.menu-item.has-sub-box > a:hover > a:after {
  border-left-color: #fff;
}
</style>