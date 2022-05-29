import tool from './tool.js';
import style from './style/index.js';


tool.style = style;

tool.addColor = function (color) {
    if (tool.isNotEmpty(color)) {
        if (color.startsWith("#")) {
            style.addColor(color);
            color = color.substring(1);
        }
    }
    return color;
}

export default tool;