import { pattern, rule } from '@/form/base.js';

let login = {
    fields: [
        {
            label: "账号",
            name: "account",
            rules: rule.account,
        },
        {
            label: "密码",
            name: "password",
            type: "password",
            rules: [
                {
                    required: true,
                    message: `密码不能为空!`
                }
            ],
        }
    ],
};

let register = {
    fields: [
        {
            label: "名称",
            name: "name",
            rules: [
                {
                    required: true,
                    message: `名称不能为空!`
                },
            ],
        },
        {
            label: "账号",
            name: "account",
            rules: rule.account,
        },
        {
            label: "邮箱",
            name: "email",
            type: "email",
            rules: rule.email,
        },
        {
            label: "密码",
            name: "password",
            type: "password",
            rules: rule.password,
        }
    ],
};

let app = {
    fields: [
        {
            label: "名称",
            name: "name",
            rules: [
                {
                    required: true,
                    message: `应用名称不能为空!`
                },
                {
                    pattern: /^[a-zA-Z0-9_]+$/,
                    message: `应用名称只能包含数字、字母、下划线!`
                },

            ],
        },
    ],
};


let model = {
    fields: [
        {
            label: "名称",
            name: "name",
            rules: [
                {
                    required: true,
                    message: `模型名称不能为空!`
                },
                {
                    pattern: /^[a-zA-Z0-9_/]+$/,
                    message: `模型名称只能包含数字、字母、下划线!`
                },

            ],
        },
        {
            label: "说明",
            name: "comment",
            rules: [
            ],
        },
    ],
};


let toolbox = {
    fields: [
        {
            label: "名称",
            name: "name",
            rules: [
                {
                    required: true,
                    message: `名称不能为空!`
                },

            ],
        },
        {
            label: "说明",
            name: "comment",
            rules: [
            ],
        },
    ],
    share: {
        fields: [
            {
                label: "分享对象",
                name: "targetUserId",
                type: "select",
                options: [
                    { text: "任意用户", value: "0" },
                ],
                defaultValue: "0",
            },
            {
                label: "权限",
                name: "power",
                type: "select",
                options: [
                    { text: "只读", value: "1" },
                    { text: "所有操作", value: "2" },
                ],
                defaultValue: "1",
                rules: [
                    {
                        required: true,
                        message: `权限不能为空!`
                    },

                ],
            },
            {
                label: "说明",
                name: "comment",
                rules: [
                ],
            },
            {
                label: "开始时间（可不填写）",
                name: "endTime",
                type: "datetime",
                rules: [
                ],
            },
            {
                label: "结束时间（可不填写）",
                name: "endTime",
                type: "datetime",
                rules: [
                ],
            },
        ],
    },
    group: {
        fields: [
            {
                label: "名称",
                name: "name",
                rules: [
                    {
                        required: true,
                        message: `名称不能为空!`
                    },

                ],
            },
            {
                label: "说明",
                name: "comment",
                rules: [
                ],
            },
        ],
        option: {
            fields: [
            ],
        },
    },
};

var nodeOptions = []
var netProxyNetTypes = [{ text: "TCP", value: "tcp" }];
let node = {
    nodeOptions: nodeOptions,
    local: {

        fields: [
            {
                label: "名称",
                name: "name",
                rules: [
                    {
                        required: true,
                        message: `名称不能为空!`
                    },

                ],
            },
            {
                label: "说明",
                name: "comment",
                rules: [
                ],
            },
            {
                label: "绑定地址(ip:port)",
                name: "bindAddress",
                rules: [
                    {
                        required: true,
                        message: `绑定地址不能为空!`
                    },

                ],
            },
            {
                label: "Token(用于节点连接验证)",
                name: "bindToken",
                rules: [
                    {
                        required: true,
                        message: `Token不能为空!`
                    },

                ],
            },
        ],
    },
    toNode: {
        fields: [
            {
                label: "名称",
                name: "name",
                rules: [
                    {
                        required: true,
                        message: `名称不能为空!`
                    },

                ],
            },
            {
                label: "说明",
                name: "comment",
                rules: [
                ],
            },
            {
                label: "目标节点ID",
                name: "serverId",
                rules: [
                    {
                        required: true,
                        message: `目标节点ID不能为空!`
                    },

                ],
            },
            {
                label: "目标节点地址(ip:port)",
                name: "connAddress",
                rules: [
                    {
                        required: true,
                        message: `目标节点地址不能为空!`
                    },

                ],
            },
            {
                label: "目标节点Token(用于节点连接验证)",
                name: "connToken",
                rules: [
                    {
                        required: true,
                        message: `目标节点Token不能为空!`
                    },

                ],
            },
        ],
    },
    fromNode: {
        fields: [
            {
                label: "名称",
                name: "name",
                rules: [
                    {
                        required: true,
                        message: `名称不能为空!`
                    },

                ],
            },
            {
                label: "说明",
                name: "comment",
                rules: [
                ],
            },
            {
                label: "节点ID",
                name: "serverId",
                rules: [
                    {
                        required: true,
                        message: `节点ID不能为空!`
                    },

                ],
            },
        ],
    },
    netProxy: {
        labelWidth: "140px",
        fields: [
            {
                label: "名称",
                name: "name",
                columnSize: 6,
                rules: [
                    {
                        required: true,
                        message: `名称不能为空!`
                    },

                ],
            },
            {
                label: "说明",
                name: "comment",
                columnSize: 6,
                rules: [
                ],
            },
            {
                label: "输入端节点",
                name: "innerServerId",
                columnSize: 4,
                type: "select",
                options: nodeOptions,
                rules: [
                    {
                        required: true,
                        message: `输入端节点不能为空!`
                    },

                ],
            },
            {
                label: "输入端类型",
                name: "innerType",
                columnSize: 4,
                type: "select",
                options: netProxyNetTypes,
                defaultValue: "tcp",
                rules: [
                    {
                        required: true,
                        message: `输入端类型不能为空!`
                    },

                ],
            },
            {
                label: "输入端绑定地址",
                name: "innerAddress",
                columnSize: 4,
                rules: [
                    {
                        required: true,
                        message: `输入端绑定地址不能为空!`
                    },

                ],
            },
            {
                label: "输出端节点",
                name: "outerServerId",
                columnSize: 4,
                type: "select",
                options: nodeOptions,
                rules: [
                    {
                        required: true,
                        message: `输出端节点不能为空!`
                    },

                ],
            },
            {
                label: "输出端类型",
                name: "outerType",
                columnSize: 4,
                type: "select",
                options: netProxyNetTypes,
                defaultValue: "tcp",
                rules: [
                    {
                        required: true,
                        message: `输出端类型不能为空!`
                    },

                ],
            },
            {
                label: "输出端绑定地址",
                name: "outerAddress",
                columnSize: 4,
                rules: [
                    {
                        required: true,
                        message: `输出端绑定地址不能为空!`
                    },

                ],
            },
        ],
    },
};
export default {
    login,
    register,
    app,
    model,
    toolbox,
    node,
};