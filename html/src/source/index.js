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
}

source.frame = {
    show: true,
    remove: false,
}

source.login = {
    show: false,
    remove: true,
    user: null,
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

source.init = (data) => {
    if (data != null) {
        source.url = data.url;
        source.api = data.api;
        source.status = "connected";
        source.ready = true;
    } else {
        source.status = "error";
        source.ready = false;
    }
}

export default source;