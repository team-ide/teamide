let source = {};

source.ready = false;
source.url = location.origin + location.pathname;
if (!source.url.endsWith('/')) {
    source.url = source.url + '/';
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
    show: true,
    remove: false,
}

source.console = {
    show: true,
    remove: false,
}

source.enum = {
};

source.log = {
};

export default source;