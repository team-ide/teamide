<template>
  <div class="toolbox-fontImage-editor">
    <template v-if="toOpenUp">
      <tm-layout height="100%">
        <tm-layout width="auto" class="">
          <div class="fontImage-box">
            <canvas
              ref="canvas"
              class="fontImage-canvas"
              :style="{
                width: `${form.width}px`,
                height: `${form.height}px`,
                left: `calc(50% - ${form.width / 2}px)`,
                top: `calc(50% - ${form.height / 2}px)`,
              }"
            >
            </canvas>
          </div>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="300px" class=""> </tm-layout>
      </tm-layout>
    </template>
    <template v-else>
      <div class="text-center color-orange pdtb-20">功能开发中，敬请期待！</div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      ready: false,
      toOpenUp: false,
      form: {
        width: 100,
        height: 100,
        fontSize: 30,
        text: `Team IDE`,
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
      // this.buildView();
    },
    refresh() {
      this.$nextTick(() => {});
    },
    buildView() {
      let canvas = this.$refs.canvas;
      var context = canvas.getContext("2d");
      context.clearRect(0, 0, this.form.width, this.form.height);

      context = canvas.getContext("2d");
      context.clearRect(0, 0, this.form.width, this.form.height);
      context.font = `${this.form.fontSize}px Font_HVD_Comic_Serif_Pro`;
      context.fillStyle = "red";
      context.textAlign = "center";
      context.textBaseline = "middle";
      let textMetrics = context.measureText(this.form.text);
      console.log(textMetrics);
      // 左上角显示文字，所以我用了x=0, y=0，但只显示部分文字。
      context.strokeStyle = "cornflowerblue";
      context.fillText(this.form.text, textMetrics.width, 0);
      context.strokeText(
        this.form.text,
        textMetrics.width,
        this.form.height / 2
      );
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-fontImage-editor {
  width: 100%;
  height: 100%;
}
.fontImage-box {
  width: 100%;
  height: 100%;
}
.fontImage-canvas {
  margin: 0px auto;
  border: 1px solid #4f4f4f;
  position: absolute;
}
</style>
