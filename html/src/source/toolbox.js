import tool from "../tool";

let toolbox = {};

toolbox.toolboxTypeOpens = [];

toolbox.tabs = [];
toolbox.activeTab = null;
toolbox.context = null;

toolbox.formatExtend = (toolboxType, data, extend) => {
    if (toolboxType.name == "ssh") {
        toolbox.formatSSHExtend(toolboxType, data, extend);
    }
};
toolbox.formatSSHExtend = (toolboxType, data, extend) => {

};

export default toolbox;