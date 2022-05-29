var co = new Object();

co.replaceAll = function (text, arg1, arg2) {
	if (co.isEmpty(text)) {
		return text;
	}
	return text.replace(new RegExp(arg1, "gm"), arg2);
};


co.createClass = function (arg1) {
	var Class = function () { }
	if (arg1) {
		Class = function (arg) {
			arg1.call(this, arg);
		}
		var Super = function () { };
		Super.prototype = arg1.prototype;
		Class.prototype = new Super();
	}
	return Class;
};
var numberindex = 0;
co.getNumber = function () {
	numberindex++;
	return new Date().getTime() - 1200000000000 + "" + Math.floor(Math.random() * 9 + 1) + "" + Math.floor(Math.random() * 9 + 1) + "" + numberindex;
};
co.trimObject = function (obj, trims) {
	trims = trims || [undefined, null];
	if (obj == null) {
		return;
	}
	if (co.isArray(obj)) {

		obj.forEach(o => {
			co.trimObject(o, trims);
		});
	} else if (co.isObject(obj)) {
		for (var key in obj) {
			var value = obj[key];
			if (co.isArray(value)) {
				co.trimObject(value, trims);
			} else if (co.isObject(value)) {
				co.trimObject(value, trims);
			} else {
				if (trims.indexOf(value) >= 0) {
					delete obj[key];
				}
			}
		}
	}
};
co.isArray = function (arg) {
	if (!co.isObject(arg)) {
		return false;
	}
	return arg.constructor == Array;
};

co.isString = function (arg) {
	if (arg == null) {
		return false;
	}
	return typeof (arg) === "string";
};
co.isObject = function (arg) {
	if (arg == null) {
		return false;
	}
	return typeof (arg) === "object";
};
co.isFunction = function (arg) {
	if (arg == null) {
		return false;
	}
	return typeof (arg) === "function";
};
co.isNumber = function (arg) {
	if (co.isEmpty(arg)) {
		return false;
	}
	return (/^[-]?[0-9]+\.?[0-9]*$/.test(arg));
};
co.isBoolean = function (arg) {
	if (co.isEmpty(arg)) {
		return false;
	}
	return typeof (arg) === "boolean";
};
co.isInteger = function (arg) {
	if (co.isEmpty(arg)) {
		return false;
	}
	return (/^-?[0-9]*$/.test(arg));
};
co.isUndefined = function (arg) {
	return typeof (arg) == "undefined";
};
co.isEmpty = function (arg) {
	return co.isUndefined(arg) || arg == null || arg.length == 0;
};
co.isNotEmpty = function (arg) {
	return !co.isEmpty(arg);
};
co.isTrue = function (arg) {
	var flag = false;
	if (arg != null) {
		if (co.isBoolean(arg)) {
			flag = arg;
		} else {
			if (arg == 'true' || arg == '1' || arg == 1) {
				flag = true;
			}
		}
	}
	return flag;
};
co.isFalse = function (arg) {

	return !co.isTrue(arg);
};
co.getKeys = function (arg) {
	if (arg == null) {
		return [];
	}
	return Object.keys(arg);
};
co.replaceList = function (arg1, arg2) {
	co.replaceArray(arg1, arg2);
};
co.replaceArray = function (arg1, arg2) {
	if (co.isArray(arg1) && co.isArray(arg2)) {
		co.trimList(arg1);
		arg2.forEach(function (one, i, array) {
			arg1.push(one);
		})
	}
};

co.trimList = function (arg) {
	co.trimArray(arg);
};
co.trimArray = function (arg) {
	if (co.isArray(arg)) {
		arg.splice(0, arg.length);
	}
};
co.trimDate = function (arg) {
	if (arg == null) {
		return arg;
	}
	let date = arg;
	if (arg instanceof Date) {
		date = co.formatDate(arg);
	}
	date = '' + date;
	return date.replace(/\D/g, '');
};
co.formatDate = function (arg, format) {
	if (arg == null) {
		return arg;
	}
	format = format || "yyyy-MM-dd hh:mm:ss";
	let date = arg;
	if (arg instanceof Date) {

	} else {
		date = co.toDate(arg);
	}
	var o = {
		"M+": date.getMonth() + 1, // month
		"d+": date.getDate(), // day
		"h+": date.getHours(), // hour
		"H+": date.getHours(), // hour
		"m+": date.getMinutes(), // minute
		"s+": date.getSeconds(), // second
		"q+": Math.floor((date.getMonth() + 3) / 3), // quarter
		"S+": date.getMilliseconds()
	};
	if (/(y+)/.test(format)) {
		format = format.replace(RegExp.$1, (date.getFullYear() + "").substr(4 - RegExp.$1.length));
	}
	for (var k in o) {
		if (new RegExp("(" + k + ")").test(format)) {
			format = format.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length));
		}
	}
	return format;

};
co.toDate = function (arg) {
	if (arg == null) {
		return arg;
	}
	if (arg instanceof Date) {
		return arg;
	}
	arg = co.trimDate(arg);
	if (arg.length == 4) {
		arg = arg + '01' + '01';
	}
	else if (arg.length == 6) {
		arg = arg + '01';
	}
	else if (arg.length == 10) {
		arg = arg + '00';
	}
	let date = '';
	let cs = arg.split('');
	for (let i = 0; i < cs.length; i++) {
		let c = cs[i];
		let a = '';
		if (i == 4 || i == 6) {
			a = '/';
		}
		if (i == 8) {
			a = ' ';
		}
		if (i == 10 || i == 12 || i == 14) {
			a = ':';
		}
		date += a + c;
	}
	return new Date(date);
};
co.getColorSuffix = function (color) {
	if (co.isNotEmpty(color) && co.startsWith('#')) {
		co.style.addColor(color);
		color = color.substring(1);
	}
	return color;
};
co.toURL = function (arg) {
	window.location.href = arg;
};
co.isRequired = function (input) {
	if (input == null) {
		return;
	}
	if (input.length != null) {
		if (input.length > 0) {
			input = input[0];
		} else {
			return;
		}
	}

	var required = input.getAttribute('isrequired');

	if (typeof (required) != 'undefined') {
		return co.isTrue(required);
	}
	required = input.getAttribute('required');
	if (typeof (required) == 'undefined') {
		return false;
	}
	if (required == 'required') {
		return true;
	}
	return co.isTrue(required);
};
co.addClass = function (dom, addClassName) {
	if (dom == null || co.isEmpty(addClassName)) {
		return;
	}
	let classNames = [];
	if (co.isNotEmpty(dom.className)) {
		classNames = dom.className.split(' ');
	}
	let addClassNames = addClassName.split(' ');

	addClassNames.forEach(addClassName => {
		addClassName = addClassName.trim();
		if (co.isNotEmpty(addClassName)) {
			let has = false;
			classNames.forEach(className => {
				className = className.trim();
				if (addClassName == className) {
					has = true;
				}
			});
			if (!has) {
				classNames.push(addClassName);
			}
		}
	});
	dom.className = classNames.join(' ');
};

co.removeClass = function (dom, removeClassName) {
	if (dom == null || co.isEmpty(removeClassName)) {
		return;
	}
	let classNames = [];
	if (co.isNotEmpty(dom.className)) {
		classNames = dom.className.split(' ');
	}
	let removeClassNames = removeClassName.split(' ');

	removeClassNames.forEach(removeClassName => {
		removeClassName = removeClassName.trim();
		if (co.isNotEmpty(removeClassName)) {
			classNames.forEach((className, index) => {
				className = className.trim();
				if (removeClassName == className) {
					classNames[index] = '';
				}
			});
		}
	});
	dom.className = classNames.join(' ');
};
co.startZIndex = 1000;
let index_number = 0;
co.getNextZIndex = function () {
	if (index_number > 0) {
		index_number++;
	}
	let index = co.startZIndex + index_number;
	if (index_number == 0) {
		index_number++;
	}
	return index;
};


co.getUrlParamValue = function (name) {
	if (name == null || name == 'undefined') {
		return null;
	}
	var searchStr = decodeURI(location.search);
	var infoIndex = searchStr.indexOf(name + "=");
	if (infoIndex == -1) {
		return null;
	}
	var searchInfo = searchStr.substring(infoIndex + name.length + 1);
	var tagIndex = searchInfo.indexOf("&");
	if (tagIndex != -1) {
		searchInfo = searchInfo.substring(0, tagIndex);
	}
	return searchInfo;
};


let zIndexsForMask = [];

co.showMask = function (zIndex) {
	if (zIndex == null) {
		zIndex = co.getNextZIndex();
	}
	if (window.COOS_MASK_DOM == null) {
		window.COOS_MASK_DOM = document.createElement("div");
		window.COOS_MASK_DOM.setAttribute("class", "tm-mask");
		document.body.appendChild(window.COOS_MASK_DOM);
	}
	co.addClass(window.COOS_MASK_DOM, 'tm-show');
	zIndexsForMask.push(zIndex);
	window.COOS_MASK_DOM.style['z-index'] = zIndex;
	return zIndexsForMask.length - 1;
};

co.hideMask = function (index) {
	if (index < 0 || index >= zIndexsForMask.length) {
		return;
	}
	zIndexsForMask.splice(index, 1);
	if (window.COOS_MASK_DOM != null) {
		if (zIndexsForMask.length == 0) {
			co.removeClass(window.COOS_MASK_DOM, 'tm-show');
		} else {
			var max = Math.max.apply(Math, zIndexsForMask);
			window.COOS_MASK_DOM.style['z-index'] = max;
		}
	}
};

export default co;