require.config({ paths: { vs: "static/monaco-editor/min/vs" } });
// vs/editor/editor.main.xxx.js  有多个文件需要导入
require(
    [
        "vs/editor/editor.main",
    ]
    , () => {
        window.monaco = monaco;
    })