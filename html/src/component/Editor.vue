<template>
  <div class="teamide-editor" ref="editor" @keydown="keydown"></div>
</template>

<script>
export default {
  props: ["source", "language", "readonly", "value", "change", "onContextMenu"],
  components: {},
  data() {
    return {
      hints: [],
    };
  },
  watch: {},
  methods: {
    setValue(value) {
      if (this.monacoInstance) {
        if (!this.getValue() != value) {
          this.isSetValue = true;
          this.monacoInstance.setValue(value);
        }
      }
    },
    getValue() {
      return this.monacoInstance.getValue();
    },
    keydown(event) {
      event = event || window.event;
      if (event.ctrlKey) {
        // Ctrl + S
        if (event.keyCode == 83) {
          event.stopPropagation && event.stopPropagation();
          event.preventDefault && event.preventDefault();
          this.$emit("save", this.getValue());
        }
      }
    },
    getSelection() {
      return this.monacoInstance
        .getModel()
        .getValueInRange(this.monacoInstance.getSelection());
      // return this.monacoInstance.getSelection().toString();
    },
    init() {
      let monaco = window.monaco;
      let languageList = monaco.languages.getLanguages();
      let yamlLanguage = null;
      let htmlLanguage = null;
      let jsLanguage = null;
      let language = null;
      languageList.forEach((one) => {
        if (one.id == "javascript") {
          jsLanguage = one;
        } else if (one.id == "html") {
          htmlLanguage = one;
        } else if (one.id == "yaml") {
          yamlLanguage = one;
        }
        if (one.id == this.language) {
          language = one;
        }
      });
      if (language == null) {
        language = htmlLanguage;
      }
      this.monacoInstance = monaco.editor.create(this.$refs.editor, {
        theme: "vs-dark", //官方自带三种主题vs, hc-black, or vs-dark
        minimap: { enabled: false }, // 缩略导航
        value: this.value || "", //编辑器初始显示文字
        language: language.id,
        selectOnLineNumbers: true, //显示行号
        roundedSelection: false,
        readOnly: this.readonly, // 只读
        cursorStyle: "line", //光标样式
        automaticLayout: true, //自动布局
        glyphMargin: true, //字形边缘
        useTabStops: false,
        fontSize: 16, //字体大小
        autoIndent: true, //自动布局
        // quickSuggestionsDelay: 500, //代码提示延时
        contextmenu: false,
      });
      this.monacoInstance.onDidChangeModelContent((e) => {
        if (this.isSetValue) {
          delete this.isSetValue;
        } else {
          this.change && this.change(this.getValue());
        }
      });
      this.monacoInstance.onContextMenu((e) => {
        this.tool.stopEvent();
        this.onContextMenu && this.onContextMenu();
      });
      //提示项设值
      //       monaco.languages.registerCompletionItemProvider("java", {
      //         provideCompletionItems: (model, position) => {
      //           console.log(model);
      //           console.log(position);
      //           let suggestions = [];

      //           suggestions.push({
      //             label: "main", // 显示的提示内容
      //             kind: monaco.languages.CompletionItemKind["Function"], // 用来显示提示内容后的不同的图标
      //             insertText: `public static void main(string[] args){
      // \t
      // }`, // 选择后粘贴到编辑器中的文字
      //             detail: "public static void main", // 提示内容后的说明
      //           });
      //           return { suggestions: suggestions };
      //         },
      //         // 光标选中当前自动补全item时触发动作，一般情况下无需处理
      //         resolveCompletionItem(item, token) {
      //           return null;
      //         },
      //       });
    },
  },
  mounted() {
    this.init();
  },
  destroyed() {
    if (this.monacoInstance != null) {
      this.monacoInstance.dispose(); //使用完成销毁实例
    }
  },
  beforeCreate() {},
};
</script>

<style >
.teamide-editor {
  width: 100%;
  height: 100%;
}
</style>
