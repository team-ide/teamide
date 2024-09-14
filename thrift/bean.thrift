namespace go bean

struct Setting {
    1: optional map<string, string> option;
}

struct ToolboxGroup {
    1: optional i64 groupId;
    2: optional string name;
    3: optional string comment;
    4: optional string option;
    5: optional i32 sequence;
}

struct Toolbox {
    1: optional i64 toolboxId;
    2: optional string toolboxType;
    3: optional i64 groupId;
    4: optional string name;
    5: optional string comment;
    6: optional string option;
    7: optional i32 visibility;
    8: optional i32 sequence;
}
