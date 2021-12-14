import server from "@/server/index.js";
import source from "@/source/index.js";

import tm from 'teamide-ui'

import md5 from 'js-md5';
let tool = {};
Object.assign(tool, tm);
tool.md5 = md5;
tool.init = function () {
    source.status = 'connecting';
    server.data().then(res => {
        if (res.code == 0) {
            let data = res.data;
            source.init(data)
        } else {
            tool.error(res.msg);
            source.init();
        }
    }).catch(() => {
        source.init();
    })
};

tool.toLogin = function () {

    source.login.remove = false;
    source.login.show = true;
    // source.login.user = {
    //     name: "张三",
    //     avatar: "static/logo.png",
    //     avatarUrl: source.url + "static/logo.png",
    // }
};

tool.toLogout = function () {
    source.login.user = null;
};

tool.formatDateByTime = function (time, format) {
    if (time == null || time <= 0) {
        return "";
    }
    return tool.formatDate(new Date(time), format);
};

tool.replaceAll = function (str, s1, s2) {
    if (str == null) {
        return str;
    }
    str = '' + str;
    return str.replace(new RegExp(s1, "gm"), s2);
};

/**
 * 获取顶部div的距离
 */
tool.getElementTop = (e) => {
    var offset = e.offsetTop
    if (e.offsetParent != null) {
        offset += tool.getElementTop(e.offsetParent)
    }
    return offset
};
/**
 * 获取左侧div的距离
 */
tool.getElementLeft = (e) => {
    var offset = e.offsetLeft
    if (e.offsetParent != null) {
        offset += tool.getElementLeft(e.offsetParent)
    }
    return offset
};
tool.getEnum = (type, value) => {
    let result = null;
    let options = source.ENUM_MAP[type];
    if (options) {
        options.forEach(one => {
            if (one.value == value) {
                result = one;
            }
        });
        if (result == null) {
            options.forEach(one => {
                if (one.text == value) {
                    result = one;
                }
            });
        }
    }
    return result;
};
tool.getEnumName = (type, value) => {
    let one = tool.getEnum(type, value);
    if (one == null) {
        return;
    }
    return one.name;
};
tool.getEnumValue = (type, value) => {
    let one = tool.getEnum(type, value);
    if (one == null) {
        return;
    }
    return one.value;
};
tool.getEnumText = (type, value) => {
    let one = tool.getEnum(type, value);
    if (one == null) {
        return;
    }
    return one.text;
};
tool.getEnumTextSpan = (type, value) => {
    let text = '';
    let color = '';
    let one = tool.getEnum(type, value);
    if (one) {
        text = one.text;
        color = one.color;
    }
    let html = '';
    html += '<span '
    if (tool.isNotEmpty(color)) {
        if (color.startsWith('#')) {
            html += ' style="color:' + color + ';" ';
        } else {
            html += ' class="color-' + color + '" ';
        }
    }
    html += ' >';
    html += text;
    html += '</span>'
    return html;
};
let chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789".split('');
tool.getRandom = function (length) {
    let res = '';
    let max = chars.length;
    let min = 0;
    for (let i = 0; i < length; i++) {
        let index = Math.floor(Math.random() * (max - min) + min);
        res += chars[index];
    }
    return res;
};

tool.getCacheKey = function (key) {
    return source.ROOT_URL + '-' + key;
};

tool.getCache = function (key) {
    key = tool.getCacheKey(key);
    return localStorage.getItem(key);
};
tool.setCache = function (key, value) {
    key = tool.getCacheKey(key);
    localStorage.setItem(key, value);
};
tool.removeCache = function (key) {
    key = tool.getCacheKey(key);
    return localStorage.removeItem(key);
};
tool.isTrimEmpty = function (arg) {
    if (arg == null) {
        return true;
    }
    return tool.isEmpty(('' + arg).trim());
};

tool.isNotTrimEmpty = function (arg) {
    return !tool.isTrimEmpty(arg);
};

export default tool;