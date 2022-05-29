
import co from '../tool.js';
import Color from './Color.js';
import config from './config.js';
import css_text from './css.js';
let style = new Object();
style.config = config;

style.getSizesCSS = function () {
	var css = '\n';
	config.sizes.forEach((one, index) => {
		let key = one.name;
		let tmKey = '';
		let sizeKey = '';
		if (co.isEmpty(key)) {
			key = '';
			tmKey = '';
			sizeKey = '';
		} else {
			tmKey = 'tm-window-' + key + ' ';
			sizeKey = '.tm-size-' + key + '';
			key = '-' + key;
		}
		var text = css_text.sizeCSS;
		text = co.replaceAll(text, 'tm-key ', tmKey);

		text = co.replaceAll(text, '.tm-size-key', sizeKey);
		text = co.replaceAll(text, '-key', key);
		for (let n in one) {
			text = co.replaceAll(text, '"' + n + '"', one[n]);
		}
		css += text;
	});
	return css;
};
style.getDistancesCSS = function () {
	var css = '\n';
	config.distances.forEach((distance, index) => {
		let key = distance;
		var text = css_text.distanceCSS;
		text = co.replaceAll(text, '-key', '-' + key + '');
		text = co.replaceAll(text, '"value"', distance + 'px');
		css += text;
		css += '\n';

	});
	return css;
};
style.getColsCSS = function () {
	var css = '\n';
	config.cols.forEach((col, index) => {
		var key = col.value;
		key = key.replace(".", "-");

		var text = css_text.colCSS;
		text = co.replaceAll(text, '-key', '-' + key + '');
		text = co.replaceAll(text, '"value"', col.width);
		css += text;
		css += '\n';

	});
	return css;
};

style.getColorsCSS = function () {
	var css = '';
	config.colors.forEach((one, index) => {
		css += style.getColorCSS(one);
	});
	return css;
};
style.getColorCSS = function (one) {
	var css = '';
	var name = one.value;
	var text = one.text;
	var colors = one.colors;
	colors.forEach((color, index) => {
		var key = name;
		if (index > 0) {
			key = name + '-' + index;
		}
		var whiteColor = "";
		if (name != 'white') {
			whiteColor = "#FFFFFF";
		}
		let colorText = css_text.colorCSS;
		colorText = co.replaceAll(colorText, '-key', '-' + key + '');
		colorText = co.replaceAll(colorText, '"whiteColor"', whiteColor);
		colorText = co.replaceAll(colorText, '"color"', color);
		css += colorText;
	});
	return css;
};

style.getComponentsCSS = function () {
	var css = '';
	config.colors.forEach((one, index) => {
		css += style.getComponentCSS(one);
	});
	return css;
};
style.getComponentCSS = function (one) {
	var css = '';
	var name = one.value;
	var text = one.text;
	var colors = one.colors;
	colors.forEach((color, index) => {
		var key = name;
		if (index > 0) {
			key = name + '-' + index;
			// return;
		}
		if (index >= colors.length - 1) {
			return;
		}
		var whiteColor = "";
		if (name != 'white') {
			whiteColor = "#FFFFFF";
		}
		var activeColor = color;
		if (colors.length >= 10) {
			if (index == 0 || index == 5) {
				activeColor = colors[8];
			} else {
				activeColor = colors[index + 1];
			}
		}

		var text = css_text.componentCSS;
		text = co.replaceAll(text, '-key', '-' + key + '');
		text = co.replaceAll(text, '"whiteColor"', whiteColor);
		text = co.replaceAll(text, '"activeColor"', activeColor);
		text = co.replaceAll(text, '"color"', color);
		css += text;
		css += '\n';
	});
	css += style.getComponentActiveCSS(one);
	return css;
};
style.getComponentActiveCSS = function (one) {
	var css = '';
	var name = one.value;
	var text = one.text;
	var colors = one.colors;
	colors.forEach((color, index) => {
		var key = name;
		if (index > 0) {
			key = name + '-' + index;
			// return;
		}
		if (index >= colors.length - 1) {
			return;
		}
		var whiteColor = "";
		if (name != 'white') {
			whiteColor = "#FFFFFF";
		}
		var activeColor = color;
		if (colors.length >= 10) {
			if (index == 0 || index == 5) {
				activeColor = colors[8];
			} else {
				activeColor = colors[index + 1];
			}
		}

		var text = css_text.componentActiveCSS;
		text = co.replaceAll(text, '-key', '-' + key + '');
		text = co.replaceAll(text, '"whiteColor"', whiteColor);
		text = co.replaceAll(text, '"activeColor"', activeColor);
		text = co.replaceAll(text, '"color"', color);
		css += text;
		css += '\n';
	});
	return css;
};
style.initStyle = function (arg) {

	if (window._COOS_STYLE_DOCUMENT == null) {
		let heads = document.getElementsByTagName('head');
		if (heads != null && heads.length > 0) {
			let head = heads[0];
			let dom = document.createElement('style');
			dom.setAttribute('id', 'COOS-STYLE');
			dom.setAttribute('type', 'text/css');
			if (head.firstChild) {
				head.insertBefore(dom, head.firstChild);
			} else {
				head.appendChild(dom);
			}
			window._COOS_STYLE_DOCUMENT = dom;
		}

		if (window._COOS_STYLE_DOCUMENT != null) {
			window._COOS_STYLE_BASE_NODE = document.createTextNode("");
			window._COOS_STYLE_DOCUMENT.appendChild(window._COOS_STYLE_BASE_NODE);

			window._COOS_STYLE_COMPONENT_NODE = document.createTextNode("");
			window._COOS_STYLE_DOCUMENT.appendChild(window._COOS_STYLE_COMPONENT_NODE);
		}
	}
};
style.appendBase = function (css) {
	style.initStyle();
	if (window._COOS_STYLE_BASE_NODE != null && css) {
		window._COOS_STYLE_BASE_NODE.appendData(css);
	}
};
style.appendComponent = function (css) {
	style.initStyle();
	if (window._COOS_STYLE_COMPONENT_NODE != null && css) {
		window._COOS_STYLE_COMPONENT_NODE.appendData(css);
	}
};

style.init = function () {
	var css = '';
	css += style.getSizesCSS();
	css += style.getDistancesCSS();
	css += style.getColsCSS();
	css += style.getColorsCSS();
	style.appendBase(css);
	css = style.getComponentsCSS();
	style.appendComponent(css);
};

style.writeColorCSS = function (one) {
	var css = style.getColorCSS(one);
	style.appendBase(css);
	css = style.getComponentCSS(one);
	style.appendComponent(css);
};
var map = {};
style.addColor = function (main) {
	if (map[main]) {
		return map[main];
	}
	if (co.isEmpty(main) || main.indexOf('#') != 0) {
		return;
	}
	var one = {};
	var color = new Color({
		main: main
	});
	one.colors = color.getColors();
	one.text = main;
	one.value = main.substring(1);
	style.writeColorCSS(one);
	map[main] = one;
	return one;
};

style.getColorClass = function (color) {
	if (co.isEmpty(color) || color.indexOf('{') >= 0) {
		return "";
	}
	if (color.indexOf('#') == 0) {
		style.addColor(color);
		color = color.substring('1');
	}
	return 'color-' + color;
};
style.getBGColorClass = function (color) {
	if (co.isEmpty(color) || color.indexOf('{') >= 0) {
		return "";
	}
	if (color.indexOf('#') == 0) {
		style.addColor(color);
		color = color.substring('1');
	}
	return 'bg-' + color;
};

style.getActiveColorClass = function (color) {
	if (co.isEmpty(color) || color.indexOf('{') >= 0) {
		return "";
	}
	if (color.indexOf('#') == 0) {
		style.addColor(color);
		color = color.substring('1');
	}
	return 'active-color-' + color;
};
style.getActiveBGColorClass = function (color) {
	if (co.isEmpty(color) || color.indexOf('{') >= 0) {
		return "";
	}
	if (color.indexOf('#') == 0) {
		style.addColor(color);
		color = color.substring('1');
	}
	return 'active-bg-' + color;
};

export default style;