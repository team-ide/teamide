import tool from "@/tool/index.js";

let source = {};

source.status = null;
source.ready = false;
source.url = null;
source.api = null;

source.header = {
    title: "Team IDE",
    toggleable: "lg",
    type: "dark",
    variant: "dark",
    navs: [
        {
            name: "首页",
            icon: "home",
            link: "/",
        },
        {
            name: "工作台",
            icon: "home",
            link: "/workspace",
        },
        {
            name: "工具箱",
            icon: "home",
            link: "/toolbox",
        },
        {
            name: "系统管理",
            icon: "home",
            link: "/manage",
            match(path) {
                if (path == '/manage' || path.indexOf('/manage/') == 0) {
                    return true;
                }
            },
        },
    ],
}

source.frame = {
    show: true,
    remove: false,
}
source.colors = ["primary", "secondary", "success", "warning", "danger", "info", "dark",]

source.login = {
    show: false,
    remove: false,
    user: null,
    headerNavs: [
        {
            name: "个人主页",
            icon: "person-circle",
            link: "/user/home",
            color: "deep-purple",
        },
        {
            name: "个人中心",
            icon: "person",
            link: "/user/home",
            color: "purple",
        },
        {
            name: "修改密码",
            icon: "shield",
            link: "/user/home",
            color: "indigo",
        },
        {
            name: "任务计划",
            icon: "calendar-minus",
            link: "/user/home",
            color: "teal",
        },
        {
            name: "消息通知",
            icon: "chat-text",
            link: "/user/home",
            color: "green",
        },
        {
            name: "帮助中心",
            icon: "exclamation",
            link: "/user/home",
            color: "lime",
        },
        {
            name: "问题建议",
            icon: "question",
            link: "/user/home",
            color: "orange",
        },
        {
            name: "设置",
            icon: "gear",
            link: "/user/home",
            color: "brown",
        },
    ],
    manageNavs: [
        {
            name: "用户管理",
            icon: "person-circle",
            navs: [
                {
                    name: "用户列表",
                    icon: "person-circle",
                    link: "/manage/user",
                },
                {
                    name: "元数据列表",
                    icon: "person-circle",
                    link: "/manage/user/metadata",
                },
                {
                    name: "授权列表",
                    icon: "person-circle",
                    link: "/manage/user/auth",
                },
            ],
        },
        {
            name: "权限管理",
            icon: "person-circle",
            navs: [
                {
                    name: "角色列表",
                    icon: "person-circle",
                    link: "/manage/role",
                },
            ],
        },
        {
            name: "企业管理",
            icon: "person-circle",
            navs: [
                {
                    name: "企业列表",
                    icon: "person-circle",
                    link: "/manage/enterprise",
                },
            ],
        },
        {
            name: "群组管理",
            icon: "person-circle",
            navs: [
                {
                    name: "群组列表",
                    icon: "person-circle",
                    link: "/manage/group",
                },
            ],
        },
        {
            name: "任务管理",
            icon: "person-circle",
            navs: [
                {
                    name: "任务列表",
                    icon: "person-circle",
                    link: "/manage/job",
                },
            ],
        },
        {
            name: "安全管理",
            icon: "person-circle",
            navs: [
                {
                    name: "日志管理",
                    icon: "person-circle",
                    link: "/manage/log",
                },
                {
                    name: "登录记录",
                    icon: "person-circle",
                    link: "/manage/login",
                },
            ],
        },
        {
            name: "系统管理",
            icon: "person-circle",
            navs: [
                {
                    name: "系统设置",
                    icon: "person-circle",
                    link: "/manage/system/setting",
                },
                {
                    name: "系统日志",
                    icon: "person-circle",
                    link: "/manage/system/log",
                },
            ],
        },
    ]
}

source.register = {
    show: false,
    remove: true,
}

source.workspace = {
    show: false,
    remove: false,
}

source.console = {
    show: false,
    remove: false,
}

source.enum = {
};

source.log = {
};

source.powers = [];
source.init = (data) => {
    if (data != null) {
        source.url = data.url;
        source.api = data.api;
    } else {
        source.status = "error";
        source.ready = false;
    }
}
source.initSession = (data) => {
    if (data != null) {
        source.login.user = data.user;
        source.powers = data.powers || [];
        tool.setJWT(data.JWT);
    } else {
        source.login.user = null;
        source.powers = [];
    }
    source.status = "connected";
    source.ready = true;
}

source.hasPower = function (action) {
    action = '' + action;
    source.powers = source.powers || [];
    let find = false;
    for (let i = 0; i < source.powers.length; i++) {
        if (source.powers[i].toLowerCase() == action.toLowerCase()) {
            find = true;
            break;
        }
    }
    return find;
}
export default source;