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
            rules: rule.password,
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
                    message: `工具名称不能为空!`
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

let toolboxOption = {
    database: {
        fields: [
            {
                label: "类型",
                name: "type",
                type: "select",
                defaultValue: "mysql",
                rules: [
                    {
                        required: true,
                        message: `数据库类型不能为空!`
                    },

                ],
                options: [
                    { text: "MySql", value: "mysql" },
                ],
            },
            {
                label: "Host（127.0.0.1）",
                name: "host",
                defaultValue: "127.0.0.1",
                rules: [
                    {
                        required: true,
                        message: `数据库连接地址不能为空!`
                    },
                ],
            },
            {
                label: "Port（3306）",
                name: "port",
                defaultValue: "3306",
                rules: [
                    {
                        required: true,
                        message: `数据库连接端口不能为空!`
                    },
                ],
            },
            {
                label: "Username",
                name: "username",
                rules: [
                ],
            },
            {
                label: "Password",
                name: "password",
                rules: [
                ],
            },
        ],
    },
    ssh: {
        fields: [
            {
                label: "类型",
                name: "type",
                type: "select",
                defaultValue: "tcp",
                rules: [
                    {
                        required: true,
                        message: `SSH类型不能为空!`
                    },

                ],
                options: [
                    { text: "TCP", value: "tcp" },
                ],
            },
            {
                label: "连接地址（127.0.0.1:22",
                name: "address",
                defaultValue: "127.0.0.1:22",
                rules: [
                    {
                        required: true,
                        message: `连接地址不能为空!`
                    },

                ],
            },
            {
                label: "User",
                name: "user",
                rules: [
                ],
            },
            {
                label: "Password",
                name: "password",
                rules: [
                ],
            },
        ],
    },
    redis: {
        fields: [
            {
                label: "连接地址（127.0.0.1:6379）",
                name: "address",
                defaultValue: "127.0.0.1:6379",
                rules: [
                    {
                        required: true,
                        message: `连接地址不能为空!`
                    },

                ],
            },
            {
                label: "密码",
                name: "auth",
                rules: [
                ],
            },
        ],
    },
    zookeeper: {
        fields: [
            {
                label: "连接地址（127.0.0.1:2181）",
                name: "address",
                defaultValue: "127.0.0.1:2181",
                rules: [
                    {
                        required: true,
                        message: `连接地址不能为空!`
                    },

                ],
            },
        ],
    },
    elasticsearch: {
        fields: [
            {
                label: "连接地址（http://127.0.0.1:9200）",
                name: "url",
                defaultValue: "http://127.0.0.1:9200",
                rules: [
                    {
                        required: true,
                        message: `连接地址不能为空!`
                    },

                ],
            },
        ],
        index: {
            fields: [
                {
                    label: "IndexName（索引）",
                    name: "indexName",
                    defaultValue: "index_xxx",
                    rules: [
                        {
                            required: true,
                            message: `索引不能为空!`
                        },
                    ],
                },
            ],
        },
    },
    kafka: {
        fields: [
            {
                label: "连接地址（127.0.0.1:9092）",
                name: "address",
                defaultValue: "127.0.0.1:9092",
                rules: [
                    {
                        required: true,
                        message: `连接地址不能为空!`
                    },

                ],
            },
        ],
        topic: {
            fields: [
                {
                    label: "Topic（主题）",
                    name: "topic",
                    defaultValue: "topic_xxx",
                    rules: [
                        {
                            required: true,
                            message: `主题不能为空!`
                        },
                    ],
                },
                {
                    label: "Partitions",
                    name: "numPartitions",
                    defaultValue: "1",
                    rules: [
                        {
                            required: true,
                            message: `分区不能为空!`
                        },
                    ],
                },
                {
                    label: "ReplicationFactor",
                    name: "replicationFactor",
                    defaultValue: "1",
                    rules: [
                        {
                            required: true,
                            message: `ReplicationFactor不能为空!`
                        },
                    ],
                },
            ],
        },
        push: {
            fields: [
                {
                    label: "Topic（主题）",
                    name: "topic",
                    rules: [
                        {
                            required: true,
                            message: `主题不能为空!`
                        },
                    ],
                },
                {
                    label: "KeyType",
                    name: "keyType",
                    defaultValue: "string",
                    rules: [
                    ],
                },
                {
                    label: "Key",
                    name: "key",
                    rules: [
                    ],
                },
                {
                    label: "ValueType",
                    name: "valueType",
                    defaultValue: "string",
                    rules: [
                    ],
                },
                {
                    label: "Value",
                    name: "value",
                    type: "textarea",
                    rules: [
                    ],
                },
            ],
        },
    },
};
export default {
    login,
    register,
    app,
    model,
    toolbox,
    toolboxOption,
};