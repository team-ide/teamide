
let keyCodeMap = {
    BackSpace: 8,
    Tab: 9,
    Clear: 12,
    Enter: [13, 108],
    Shift: 16,
    Ctrl: 17,
    Alt: 18,
    CapeLock: 20,
    Esc: 27,
    Spacebar: 32,
    PageUp: 33,
    PageDown: 34,
    End: 35,
    Home: 36,
    Left: 37,
    Up: 38,
    Right: 39,
    Dw: 40,
    Insert: 45,
    Delete: 46,


    0: [48, 96],
    1: [49, 97],
    2: [50, 98],
    3: [51, 99],
    4: [52, 100],
    5: [53, 101],
    6: [54, 102],
    7: [55, 103],
    8: [56, 104],
    9: [57, 105],

    A: 65,
    B: 66,
    C: 67,
    D: 68,
    E: 69,
    F: 70,
    G: 71,
    H: 72,
    I: 73,
    J: 74,
    K: 75,
    L: 76,
    M: 77,
    N: 78,
    O: 79,
    P: 80,
    Q: 81,
    R: 82,
    S: 83,
    T: 84,
    U: 85,
    V: 86,
    W: 87,
    X: 88,
    Y: 89,
    Z: 90,

    "*": 106,
    "+": 107,
    // Enter: 108,
    "-": 109,
    ".": 110,
    "/": 111,

    F1: 112,
    F2: 113,
    F3: 114,
    F4: 115,
    F5: 116,
    F6: 117,
    F7: 118,
    F8: 119,
    F9: 120,
    F10: 121,
    F11: 122,
    F12: 123,
};

let keyEvent = {};
for (let key in keyCodeMap) {
    let codes = keyCodeMap[key];
    if (!codes.length) {
        codes = [codes];
    }

    keyEvent['keyIs' + key] = (event) => {
        event = event || window.event;
        return codes.indexOf(event.keyCode) >= 0;
    }
    keyEvent['keyIsCtrl' + key] = (event) => {
        event = event || window.event;
        return event.ctrlKey && codes.indexOf(event.keyCode) >= 0;
    }
    keyEvent['keyIsShift' + key] = (event) => {
        event = event || window.event;
        return event.shiftKey && codes.indexOf(event.keyCode) >= 0;
    }
    keyEvent['keyIsAlt' + key] = (event) => {
        event = event || window.event;
        return event.altKey && codes.indexOf(event.keyCode) >= 0;
    }
}
export default {
    keyEvent
}