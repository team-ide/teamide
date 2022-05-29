
var common = {
	black: '#00000',
	white: '#FFFFFF'
};
var Color = function (options) {
	options = options || {};
	this.options = options;
	this.initOptions();
	this.init();
};

Color.prototype.initOptions = function () {
	this.main = this.options.main;
	if (this.main == null || this.main.trim().length == 0) {
		throw new Error('main不能为空，请设置主色');
	}
	this.contrastThreshold = this.options.contrastThreshold || 3;
	this.tonalOffset = this.options.tonalOffset || 0.2;

};
Color.prototype.init = function () {
	this.initColors();
};
Color.prototype.initColors = function () {
	this.colors = [];
	this.colors.push(this.main);

	this.colors.push(lighten(this.main, .4));
	this.colors.push(lighten(this.main, .3));
	this.colors.push(lighten(this.main, .2));
	this.colors.push(lighten(this.main, .1));
	this.colors.push(this.main);
	this.colors.push(darken(this.main, .1));
	this.colors.push(darken(this.main, .2));
	this.colors.push(darken(this.main, .3));
	this.colors.push(darken(this.main, .4));
};
Color.prototype.getColors = function () {
	return this.colors;
};
Color.prototype.getContrastText = function (background) {
	return getContrastRatio(background, common.white) >= this.contrastThreshold ? common.white : common.black;
};
/**
 * 返回其值仅限于给定范围的数字。
 */
function clamp(value) {
	var min = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : 0;
	var max = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : 1;

	if (value < min) {
		return min;
	}

	if (value > max) {
		return max;
	}

	return value;
}
/**
 * 将颜色从CSS十六进制格式转换为CSS rgb格式。
 */
function convertHexToRGB(color) {
	color = color.substr(1);
	var re = new RegExp(".{1,".concat(color.length / 3, "}"), 'g');
	var colors = color.match(re);

	if (colors && colors[0].length === 1) {
		colors = colors.map(function (n) {
			return n + n;
		});
	}

	return colors ? "rgb(".concat(colors.map(function (n) {
		return parseInt(n, 16);
	}).join(', '), ")") : '';
}
/**
 * 返回具有颜色类型和值的对象. 注意：不支持rgb％值。
 */
function decomposeColor(color) {
	if (color.charAt(0) === '#') {
		return decomposeColor(convertHexToRGB(color));
	}

	var marker = color.indexOf('(');
	var type = color.substring(0, marker);
	var values = color.substring(marker + 1, color.length - 1).split(',');
	values = values.map(function (value) {
		return parseFloat(value);
	});
	{
		if (['rgb', 'rgba', 'hsl', 'hsla'].indexOf(type) === -1) {
			throw new Error('紧持以下格式：#nnn，＃nnnnnn，rgb（），rgba（），hsl（），hsla（）');
		}
	}

	return {
		type: type,
		values: values
	};
}
/**
 * 将具有类型和值的颜色对象转换为字符串。
 */
function recomposeColor(color) {
	var type = color.type;
	var values = color.values;
	if (type.indexOf('rgb') !== -1) {
		// Only convert the first 3 values to int (i.e. not alpha)
		values = values.map(function (n, i) {
			return i < 3 ? parseInt(n, 10) : n;
		});
	}
	if (type.indexOf('hsl') !== -1) {
		values[1] = "".concat(values[1], "%");
		values[2] = "".concat(values[2], "%");
	}

	return "".concat(color.type, "(").concat(values.join(', '), ")");
}
/**
 * 计算两种颜色之间的对比度。
 */
function getContrastRatio(foreground, background) {
	var lumA = getLuminance(foreground);
	var lumB = getLuminance(background);
	return (Math.max(lumA, lumB) + 0.05) / (Math.min(lumA, lumB) + 0.05);
}
/**
 * 颜色空间中任何点的相对亮度，对于最暗的黑色和黑色标准化为0,1为最轻的白色。
 */
function getLuminance(color) {
	var decomposedColor = decomposeColor(color);

	if (decomposedColor.type.indexOf('rgb') !== -1) {
		var rgb = decomposedColor.values.map(function (val) {
			val /= 255; // normalized

			return val <= 0.03928 ? val / 12.92 : Math.pow((val + 0.055) / 1.055, 2.4);
		}); // Truncate at 3 digits

		return Number((0.2126 * rgb[0] + 0.7152 * rgb[1] + 0.0722 * rgb[2]).toFixed(3));
	} // else if (decomposedColor.type.indexOf('hsl') !== -1)

	return decomposedColor.values[2] / 100;
}
/**
 * 根据亮度，使颜色变暗或变浅。 浅色变暗，暗色变浅。
 */
function emphasize(color) {
	var coefficient = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : 0.15;
	return getLuminance(color) > 0.5 ? darken(color, coefficient) : lighten(color, coefficient);
}
/**
 * 设置颜色的绝对透明度。 任何现有的alpha值都会被覆盖。
 */
function fade(color, value) {
	if (!color)
		return color;
	color = decomposeColor(color);
	value = clamp(value);

	if (color.type === 'rgb' || color.type === 'hsl') {
		color.type += 'a';
	}

	color.values[3] = value;
	return recomposeColor(color);
}
/**
 * 使颜色变暗。
 */
function darken(color, coefficient) {
	if (!color)
		return color;
	color = decomposeColor(color);
	coefficient = clamp(coefficient);

	if (color.type.indexOf('hsl') !== -1) {
		color.values[2] *= 1 - coefficient;
	} else if (color.type.indexOf('rgb') !== -1) {
		for (var i = 0; i < 3; i += 1) {
			color.values[i] *= 1 - coefficient;
		}
	}

	return recomposeColor(color);
}
/**
 * 淡化颜色。
 */
function lighten(color, coefficient) {
	if (!color)
		return color;
	color = decomposeColor(color);
	coefficient = clamp(coefficient);

	if (color.type.indexOf('hsl') !== -1) {
		color.values[2] += (100 - color.values[2]) * coefficient;
	} else if (color.type.indexOf('rgb') !== -1) {
		for (var i = 0; i < 3; i += 1) {
			color.values[i] += (255 - color.values[i]) * coefficient;
		}
	}

	return recomposeColor(color);
}

Color.getContrastText = function (background) {
	return getContrastRatio(background, common.white) >= 3 ? common.white : common.black;
};

export default Color;