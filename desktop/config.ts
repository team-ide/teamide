export default {
    tray: {
        toolTip: "Team · IDE",
    },
    window: {
        title: "Team · IDE",
        width: 1440,
        height: 900,
        // 首页地址 如：./assets/html/index.html 或 http://127.0.0.1:8080/xxx
        index: "",
        // 使用 服务输出的地址
        // 服务 控制台 输出 字符串 如：event:serverUrl:http://127.0.0.1:8080/xxx
        useServerUrl: true,
        // 开启 关闭窗口 最小化
        hideWhenClose: true,
        // 开启 启动最小化
        hideWhenStart: false,
    },
    // 服务配置
    server: {
        // 服务根目录
        dir: "./assets/server",
        // so、dll库等
        libDir: "./assets/server/lib",
        // darwin 系统服务配置
        darwin: {
            exec: "./server",
            args: ["--isElectron"],
        },
        // linux 系统服务配置
        linux: {
            exec: "./server",
            args: ["--isElectron"],
        },
        // win 系统服务配置
        win: {
            exec: "./server.exe",
            args: ["--isElectron"],
        }
    },
}