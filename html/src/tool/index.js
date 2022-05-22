import server from "@/server/index.js";
import source from "@/source/index.js";
import jQuery from 'jquery'

import tm from 'teamide-ui'

import md5 from 'js-md5';
import cryptoJS from 'crypto-js';
import keyCode from '@/tool/keyCode.js';

let tool = {};
Object.assign(tool, tm);
Object.assign(tool, keyCode.keyEvent);
tool.md5 = md5;
tool.jQuery = jQuery;

tool.init = function () {
    source.status = 'connecting';
    server.data().then(res => {
        if (res.code == 0) {
            let data = res.data;
            source.init(data);

            tool.initSession();
        } else {
            tool.error(res.msg);
            source.init();
        }
    }).catch(() => {
        source.init();
    })
};
var sessionLoadding = false;
var refreshSessionStart = false;
function refreshSession() {
    function nextContinue() {
        setTimeout(() => {
            refreshSession();
        }, 1000 * 60 * 10);
    }
    if (sessionLoadding) {
        nextContinue();
    } else {
        server.session().then(res => {
            if (res.code == 0) {
                let data = res.data;
                source.initSession(data)
            } else {
                source.initSession();
            }
            nextContinue();
        }).catch(() => {
            nextContinue();
        })
    }
}
tool.initSession = function () {
    sessionLoadding = true;
    server.session().then(res => {
        if (res.code == 0) {
            let data = res.data;
            source.initSession(data)
        } else {
            tool.error(res.msg);
            source.initSession();
        }
        sessionLoadding = false;
    }).catch(() => {
        source.initSession();
        sessionLoadding = false;
    })
    if (!refreshSessionStart) {
        refreshSessionStart = true;
        setTimeout(() => {
            refreshSession();
        }, 1000 * 60 * 10);
    }
};
tool.isManagePage = function (path) {
    if (path == '/manage' || path.indexOf('/manage/') == 0) {
        return true;
    }
    return false;
};
tool.isUserPage = function (path) {
    if (path == '/user' || path.indexOf('/user/') == 0) {
        return true;
    }
    return false;
};
tool.isToolboxPage = function (path) {
    if (path == '/toolbox' || path.indexOf('/toolbox/') == 0) {
        return true;
    }
    return false;
};
tool.isApplicationPage = function (path) {
    if (path == '/application' || path.indexOf('/application/') == 0) {
        return true;
    }
    return false;
};
tool.isWorkspacePage = function (path) {
    if (path == '/workspace' || path.indexOf('/workspace/') == 0) {
        return true;
    }
    return false;
};
tool.toLogin = function () {
    tool.hideRegister();
    source.login.remove = false;
    source.login.show = true;
};

tool.hideLogin = function () {
    source.login.remove = true;
    source.login.show = false;
};

tool.toRegister = function () {
    tool.hideLogin();
    source.register.remove = false;
    source.register.show = true;
};

tool.hideRegister = function () {
    source.register.remove = true;
    source.register.show = false;
};

tool.toLogout = function () {
    server.logout().then(res => {
        if (res.code == 0) {
            tool.setJWT("");
            tool.initSession();
        }
    }).catch(() => {
    })
};

tool.stopEvent = function (event) {
    event = event || window.event;
    if (event) {
        event.stopPropagation && event.stopPropagation();
        event.preventDefault && event.preventDefault();
    }
};
tool.setJWT = function (jwt) {
    if (tool.isNotEmpty(jwt)) {
        tool.setCookie("team-ide-jwt", jwt, 60);
    } else {
        tool.setCookie("team-ide-jwt", jwt, 0);
    }
}
tool.getJWT = function () {
    return tool.getCookie("team-ide-jwt");
}
tool.setCookie = function (cname, cvalue, exms) {
    var d = new Date();
    d.setTime(d.getTime() + (exms * 60 * 1000));
    var expires = "expires=" + d.toGMTString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
}
tool.getCookie = function (cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i].trim();
        if (c.indexOf(name) == 0) { return c.substring(name.length, c.length); }
    }
    return "";
}
tool.stringToJSON = function (value) {
    let data = null;
    if (tool.isNotEmpty(value)) {
        try {
            data = JSON.parse(value);
        } catch (error) {
            try {
                data = eval("(" + value + ")");
            } catch (error2) {
                throw error;
            }
        }
    }
    return data;
};
tool.up = function (bean, name, value) {
    if (bean != null && bean[name] != null) {
        let index = bean[name].indexOf(value);
        if (index > 0) {
            bean[name].splice(index, 1);
            bean[name].splice(index - 1, 0, value);
        }
    }
};
tool.down = function (bean, name, value) {
    if (bean != null && bean[name] != null) {
        let index = bean[name].indexOf(value);
        if (index >= 0 && index < bean[name].length - 1) {
            bean[name].splice(index, 1);
            bean[name].splice(index + 1, 0, value);
        }
    }
};
tool.copyText = async function (text) {
    let result = await navigator.clipboard.writeText(text)
    return result;
};
tool.readClipboardText = async function () {
    let result = await navigator.permissions.query({ name: 'clipboard-read' })
    if (result.state == 'granted' || result.state == 'prompt') { //读取剪贴板 
        let text = await navigator.clipboard.readText()
        return text;
    } else {
        tool.warn('请允许读取剪贴板！')
        return null;
    }
};
tool.byteToString = function (arr) {
    if (typeof arr === 'string') {
        return arr;
    }
    var str = '',
        _arr = arr;
    for (var i = 0; i < _arr.length; i++) {
        var one = _arr[i].toString(2),
            v = one.match(/^1+?(?=0)/);
        if (v && one.length == 8) {
            var bytesLength = v[0].length;
            var store = _arr[i].toString(2).slice(7 - bytesLength);
            for (var st = 1; st < bytesLength; st++) {
                store += _arr[st + i].toString(2).slice(2);
            }
            str += String.fromCharCode(parseInt(store, 2));
            i += bytesLength - 1;
        } else {
            str += String.fromCharCode(_arr[i]);
        }
    }
    return str;
}

tool.stringToByte = function (str) {
    var bytes = new Array();
    var len, c;
    len = str.length;
    for (var i = 0; i < len; i++) {
        c = str.charCodeAt(i);
        if (c >= 0x010000 && c <= 0x10FFFF) {
            bytes.push(((c >> 18) & 0x07) | 0xF0);
            bytes.push(((c >> 12) & 0x3F) | 0x80);
            bytes.push(((c >> 6) & 0x3F) | 0x80);
            bytes.push((c & 0x3F) | 0x80);
        } else if (c >= 0x000800 && c <= 0x00FFFF) {
            bytes.push(((c >> 12) & 0x0F) | 0xE0);
            bytes.push(((c >> 6) & 0x3F) | 0x80);
            bytes.push((c & 0x3F) | 0x80);
        } else if (c >= 0x000080 && c <= 0x0007FF) {
            bytes.push(((c >> 6) & 0x1F) | 0xC0);
            bytes.push((c & 0x3F) | 0x80);
        } else {
            bytes.push(c & 0xFF);
        }
    }
    return bytes;
}
let k = tool.byteToString([81, 53, 54, 104, 70, 65, 97, 117, 87, 107, 49, 56, 71, 121, 50, 105]);

//加密
tool.encrypt = function (word, key) {
    const keyParse = cryptoJS.enc.Utf8.parse(key)
    const srcs = cryptoJS.enc.Utf8.parse(word)
    const encrypted = cryptoJS.AES.encrypt(srcs, keyParse, {
        iv: keyParse,
        mode: cryptoJS.mode.CBC,
        padding: cryptoJS.pad.Pkcs7
    })
    return encrypted.toString()
};
//解密
tool.decrypt = function (word, key) {
    const keyParse = cryptoJS.enc.Utf8.parse(key)
    const decrypted = cryptoJS.AES.decrypt(word, keyParse, {
        iv: keyParse,
        mode: cryptoJS.mode.CBC,
        padding: cryptoJS.pad.Pkcs7
    })
    return cryptoJS.enc.Utf8.stringify(decrypted).toString()
};
//随机生成指定数量的16进制key
tool.generatekey = function (num) {
    let library = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let key = "";
    for (var i = 0; i < num; i++) {
        let randomPoz = Math.floor(Math.random() * library.length);
        key += library.substring(randomPoz, randomPoz + 1);
    }
    return key;
};

// 加密
tool.aesEncrypt = function (str) {
    return tool.encrypt(str, k);
};
// 解密
tool.aesDecrypt = function (str) {
    return tool.decrypt(str, k);
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

tool.computeFontSize = function (str, size, family) {
    let spanDom = document.createElement("span");
    spanDom.style.fontSize = size;
    spanDom.style.opacity = "0";
    if (family) {
        spanDom.style.fontFamily = family;
    }
    if (tool.isNotEmpty(str)) {
        str = tool.replaceAll(str, " ", "&nbsp;")
        str = tool.replaceAll(str, "<", "&lt;")
        str = tool.replaceAll(str, ">", "&gt;")
    }
    spanDom.innerHTML = str;
    document.body.append(spanDom);
    let sizeD = {};
    sizeD.width = spanDom.offsetWidth;
    sizeD.height = spanDom.offsetHeight;
    spanDom.remove();
    return sizeD;
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

tool.initTreeWidth = function (tree, treeBox) {

    if (tree.initWidthIng) {
        return;
    }
    tree.initWidthIng = true;

    tree.$nextTick(() => {
        let hasActive = false;
        if (
            tree.$el.getElementsByClassName("v-enter-active").length > 0 ||
            tree.$el.getElementsByClassName("v-leave-active").length > 0
        ) {
            hasActive = true;
        }
        if (hasActive) {
            window.setTimeout(() => {
                tree.initWidthIng = false;
                tool.initTreeWidth(tree, treeBox);
            }, 100);
        } else {
            tree.$el.style.minWidth = "0px";
            let minWidth = 0;
            if (treeBox) {
                if (minWidth < treeBox.scrollWidth) {
                    minWidth = treeBox.scrollWidth;
                }
            }
            tree.minWidth = minWidth;
            tree.$el.style.minWidth = minWidth + "px";
            tree.initWidthIng = false;
        }
    });
};

tool.initInputWidth = (event, options) => {
    let target = tool.jQuery(event.target || event);
    let value = target.val();
    let str = value;
    let w = 10;
    if (target.find('option').length > 0) {
        w = 30;
        target.find('option').each((index, one) => {
            one = tool.jQuery(one);
            if (one.attr('value') == value) {
                str = one.attr('text') || one.text();
            }
        });
    }
    let compute = tool.computeFontSize(str, "12px");
    tool.jQuery(target).css('width', (compute.width + w) + 'px')
};
export default tool;