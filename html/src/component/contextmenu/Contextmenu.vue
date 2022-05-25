<template>
  <div
    class="contextmenu-box"
    :class="{
      show: showContextmenu,
      showbottom: showbottom,
      showleft: showleft,
    }"
    :style="{
      top: top,
      left: left,
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
      showleft: false,
      showContextmenu: false,
      top: null,
      left: null,
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

      this.$nextTick(() => {
        let left = clientX + "px";
        let top = clientY + "px";
        let bottomHeight = window.innerHeight - clientY;
        let menuHeight = this.$el.offsetHeight;
        let showbottom = false;
        let showleft = false;
        let hasSub = false;
        if (this.contextmenu.menus) {
          this.contextmenu.menus.forEach((one) => {
            if (one.menus && one.menus.length > 0) {
              hasSub = true;
            }
          });
        }
        if (bottomHeight < menuHeight + 30) {
          showbottom = true;
          top = clientY - menuHeight + "px";
        }

        let rightWidth = window.innerWidth - clientX;
        let menuWidth = this.$el.offsetWidth;
        let offsetLeft = 30;
        if (hasSub) {
          offsetLeft = 200;
        }
        if (rightWidth < menuWidth + offsetLeft) {
          showleft = true;
          left = clientX - menuWidth + "px";
        }

        this.left = left;
        this.top = top;
        this.showbottom = showbottom;
        this.showleft = showleft;
        this.showContextmenu = true;
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
  float: left;
  min-width: 160px;
  padding: 5px 0;
  margin: 2px 0 0;
  list-style: none;
  background-color: #fff;
  border: 1px solid #ccc;
  font-size: 13px;
  border-radius: 1px;
  text-align: left;
  transform: scale(0);
}
.contextmenu-box.show {
  transform: scale(1);
}
</style>