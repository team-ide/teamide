<template>
  <div
    class="contextmenu-box"
    :class="{ show: showContextmenu, showbottom: showbottom }"
    :style="{
      top: contextmenu.top,
      left: contextmenu.left,
      'z-index': contextmenu.zIndex,
    }"
  >
    <template v-if="contextmenu.menus != null">
      <ContextmenuMenus
        :menus="contextmenu.menus"
        :contextmenu="contextmenu"
      ></ContextmenuMenus>
    </template>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["contextmenu"],
  data() {
    return {
      showbottom: false,
      showContextmenu: false,
    };
  },
  beforeMount() {},
  watch: {},
  methods: {
    hide() {
      this.showContextmenu = false;
    },
    show(event) {
      this.showing = true;
      event = event || window.event;
      let clientX = event.clientX;
      let clientY = event.clientY;
      this.contextmenu.left = clientX + "px";
      this.contextmenu.top = clientY + "px";
      this.showContextmenu = true;
      this.showbottom = false;

      this.$nextTick(() => {
        let bottomHeight = window.innerHeight - clientY;
        let menuHeight = this.$el.offsetHeight;
        if (bottomHeight < menuHeight + 30) {
          this.showbottom = true;
          this.contextmenu.top = clientY - menuHeight + "px";
        }

        delete this.showing;
      });
    },
  },
  mounted() {
    if (this.contextmenu) {
      this.contextmenu.show = this.show;
      this.contextmenu.hide = this.hide;
    }
    document.addEventListener("mouseup", (e) => {
      if (this.$el.contains(e.target)) {
        return;
      }
      if (!this.showing || e.button != 2) {
        this.hide();
      }
    });
  },
};
</script>

<style >
.contextmenu-box,
.contextmenu-box ul,
.contextmenu-box li {
  list-style: none;
  padding: 0px;
  margin: 0px;
  user-select: none;
}
</style>
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style  scoped>
.contextmenu-box {
  position: fixed;
  top: 0px;
  left: 0;
  z-index: 1000;
  display: none;
  float: left;
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
}
.contextmenu-box.show {
  display: block;
}
</style>