const config = {};
config.distances = [];
for (var i = 0; i <= 100; i++) {
	if (i <= 10) {
		config.distances.push(i);
	} else {
		config.distances.push(i);
	}
}

config.sizes = [
	{
		name: "",
		text: "默认",
		"font-size": "15px",
		padding: "3px 13px",
		height: "34px",
		"line-height": "26px",
		"img-height": "26px",
		"circle-line-height": "32px"
	},
	{
		name: "xs",
		text: "最小",
		"font-size": "12px",
		padding: "1px 5px",
		height: "18px",
		"line-height": "14px",
		"img-height": "14px",
		"circle-line-height": "18px"
	},
	{
		name: "sm",
		text: "小",
		"font-size": "14px",
		padding: "2px 8px",
		height: "26px",
		"line-height": "20px",
		"img-height": "20px",
		"circle-line-height": "24px"
	},
	{
		name: "md",
		text: "大",
		"font-size": "18px",
		padding: "6px 18px",
		height: "42px",
		"line-height": "28px",
		"img-height": "28px",
		"circle-line-height": "40px"
	},
	{
		name: "lg",
		text: "最大",
		"font-size": "22px",
		padding: "8px 25px",
		height: "50px",
		"line-height": "32px",
		"img-height": "32px",
		"circle-line-height": "48px"
	}
];

config.cols = [{
	text: "0格",
	value: "0",
	width: "0%"
}, {
	text: "1格",
	value: "1",
	width: "8.33333333%"
}, {
	text: "2格",
	value: "2",
	width: "16.66666667%"
}, {
	text: "3格",
	value: "3",
	width: "25%"
}, {
	text: "4格",
	value: "4",
	width: "33.33333333%"
}, {
	text: "5格",
	value: "5",
	width: "41.66666667%"
}, {
	text: "6格",
	value: "6",
	width: "50%"
}, {
	text: "7格",
	value: "7",
	width: "58.33333333%"
}, {
	text: "8格",
	value: "8",
	width: "66.66666667%"
}, {
	text: "9格",
	value: "9",
	width: "75%"
}, {
	text: "10格",
	value: "10",
	width: "83.33333333%"
}, {
	text: "11格",
	value: "11",
	width: "91.66666667%"
}, {
	text: "12格",
	value: "12",
	width: "100%"
}, {
	text: "0.5格",
	value: "0.5",
	width: "4.166666665%"
}, {
	text: "1.5格",
	value: "1.5",
	width: "12.5%"
}, {
	text: "2.5格",
	value: "2.5",
	width: "20.833333335%"
}, {
	text: "3.5格",
	value: "3.5",
	width: "29.166666665%"
}, {
	text: "4.5格",
	value: "4.5",
	width: "37.5%"
}, {
	text: "5.5格",
	value: "5.5",
	width: "45.833333335%"
}, {
	text: "6.5格",
	value: "6.5",
	width: "54.166666665%"
}, {
	text: "7.5格",
	value: "7.5",
	width: "62.5%"
}, {
	text: "8.5格",
	value: "8.5",
	width: "70.833333335%"
}, {
	text: "9.5格",
	value: "9.5",
	width: "79.166666665%"
}, {
	text: "10.5格",
	value: "10.5",
	width: "87.5%"
}, {
	text: "11.5格",
	value: "11.5",
	width: "95.833333335%"
}];

config.colors = [
	{
		text: "白色",
		value: "white",
		colors: ["#FFFFFF"]
	}, {
		text: "黑色",
		value: "black",
		colors: ["#000000"]
	},
	{
		text: "主要",
		value: "primary",
		colors: ["#2196F3", "#BBDEFB", "#90CAF9", "#64B5F6", "#42A5F5", "#2196F3", "#1E88E5", "#1976D2", "#1565C0", "#0D47A1"]
	},
	{
		text: "信息",
		value: "info",
		colors: ["#607D8B", "#CFD8DC", "#B0BEC5", "#90A4AE", "#78909C", "#607D8B", "#546E7A", "#455A64", "#37474F", "#263238"]
	},
	{
		text: "成功",
		value: "success",
		colors: ["#4CAF50", "#C8E6C9", "#A5D6A7", "#81C784", "#66BB6A", "#4CAF50", "#43A047", "#388E3C", "#2E7D32", "#1B5E20"]
	},
	{
		text: "警告",
		value: "warn",
		colors: ["#FF9800", "#FFE0B2", "#FFCC80", "#FFB74D", "#FFA726", "#FF9800", "#FB8C00", "#F57C00", "#EF6C00", "#E65100"]
	},
	{
		text: "危险",
		value: "danger",
		colors: ["#F44336", "#FFCDD2", "#EF9A9A", "#E57373", "#EF5350", "#F44336", "#E53935", "#D32F2F", "#C62828", "#B71C1C"]
	},
	{
		text: "错误",
		value: "error",
		colors: ["#F44336", "#FFCDD2", "#EF9A9A", "#E57373", "#EF5350", "#F44336", "#E53935", "#D32F2F", "#C62828", "#B71C1C"]
	},
	{
		text: "确定",
		value: "define",
		colors: ["#2196F3", "#BBDEFB", "#90CAF9", "#64B5F6", "#42A5F5", "#2196F3", "#1E88E5", "#1976D2", "#1565C0", "#0D47A1"]
	},
	{
		text: "取消",
		value: "cancel",
		colors: ["#9E9E9E", "#F5F5F5", "#EEEEEE", "#E0E0E0", "#BDBDBD", "#9E9E9E", "#757575", "#616161", "#424242", "#212121"]
	}, {
		text: "红色",
		value: "red",
		colors: ["#F44336", "#FFCDD2", "#EF9A9A", "#E57373", "#EF5350", "#F44336", "#E53935", "#D32F2F", "#C62828", "#B71C1C"]
	}, {
		text: "粉色",
		value: "pink",
		colors: ["#E91E63", "#F8BBD0", "#F48FB1", "#F06292", "#EC407A", "#E91E63", "#D81B60", "#C2185B", "#AD1457", "#880E4F"]
	}, {
		text: "紫色",
		value: "purple",
		colors: ["#9C27B0", "#E1BEE7", "#CE93D8", "#BA68C8", "#AB47BC", "#9C27B0", "#8E24AA", "#7B1FA2", "#6A1B9A", "#4A148C"]
	}, {
		text: "深紫色",
		value: "deep-purple",
		colors: ["#673AB7", "#D1C4E9", "#B39DDB", "#9575CD", "#7E57C2", "#673AB7", "#5E35B1", "#512DA8", "#4527A0", "#311B92"]
	}, {
		text: "靛青色",
		value: "indigo",
		colors: ["#3F51B5", "#C5CAE9", "#9FA8DA", "#7986CB", "#5C6BC0", "#3F51B5", "#3949AB", "#303F9F", "#283593", "#1A237E"]
	}, {
		text: "蓝色",
		value: "blue",
		colors: ["#2196F3", "#BBDEFB", "#90CAF9", "#64B5F6", "#42A5F5", "#2196F3", "#1E88E5", "#1976D2", "#1565C0", "#0D47A1"]
	}, {
		text: "浅蓝色",
		value: "light-blue",
		colors: ["#03A9F4", "#B3E5FC", "#81D4FA", "#4FC3F7", "#29B6F6", "#03A9F4", "#039BE5", "#0288D1", "#0277BD", "#01579B"]
	}, {
		text: "青色",
		value: "cyan",
		colors: ["#00BCD4", "#B2EBF2", "#80DEEA", "#4DD0E1", "#26C6DA", "#00BCD4", "#00ACC1", "#0097A7", "#00838F", "#006064"]
	}, {
		text: "兰绿色",
		value: "teal",
		colors: ["#009688", "#B2DFDB", "#80CBC4", "#4DB6AC", "#26A69A", "#009688", "#00897B", "#00796B", "#00695C", "#004D40"]
	}, {
		text: "绿色",
		value: "green",
		colors: ["#4CAF50", "#C8E6C9", "#A5D6A7", "#81C784", "#66BB6A", "#4CAF50", "#43A047", "#388E3C", "#2E7D32", "#1B5E20"]
	}, {
		text: "浅绿",
		value: "light-green",
		colors: ["#8BC34A", "#DCEDC8", "#C5E1A5", "#AED581", "#9CCC65", "#8BC34A", "#7CB342", "#689F38", "#558B2F", "#33691E"]
	}, {
		text: "青柠色",
		value: "lime",
		colors: ["#CDDC39", "#F0F4C3", "#E6EE9C", "#DCE775", "#D4E157", "#CDDC39", "#C0CA33", "#AFB42B", "#9E9D24", "#827717"]
	}, {
		text: "黄色",
		value: "yellow",
		colors: ["#FFEB3B", "#FFF9C4", "#FFF59D", "#FFF176", "#FFEE58", "#FFEB3B", "#FDD835", "#FBC02D", "#F9A825", "#F57F17"]
	}, {
		text: "琥珀色",
		value: "amber",
		colors: ["#FFC107", "#FFECB3", "#FFE082", "#FFD54F", "#FFCA28", "#FFC107", "#FFB300", "#FFA000", "#FF8F00", "#FF6F00"]
	}, {
		text: "橙色",
		value: "orange",
		colors: ["#FF9800", "#FFE0B2", "#FFCC80", "#FFB74D", "#FFA726", "#FF9800", "#FB8C00", "#F57C00", "#EF6C00", "#E65100"]
	}, {
		text: "深橙色",
		value: "deep-orange",
		colors: ["#FF5722", "#FFCCBC", "#FFAB91", "#FF8A65", "#FF7043", "#FF5722", "#F4511E", "#E64A19", "#D84315", "#BF360C"]
	}, {
		text: "棕色",
		value: "brown",
		colors: ["#795548", "#D7CCC8", "#BCAAA4", "#A1887F", "#8D6E63", "#795548", "#6D4C41", "#5D4037", "#4E342E", "#3E2723"]
	}, {
		text: "灰色",
		value: "grey",
		colors: ["#9E9E9E", "#F5F5F5", "#EEEEEE", "#E0E0E0", "#BDBDBD", "#9E9E9E", "#757575", "#616161", "#424242", "#212121"]
	}, {
		text: "灰色",
		value: "gray",
		colors: ["#9E9E9E", "#F5F5F5", "#EEEEEE", "#E0E0E0", "#BDBDBD", "#9E9E9E", "#757575", "#616161", "#424242", "#212121"]
	}, {
		text: "蓝灰色",
		value: "blue-grey",
		colors: ["#607D8B", "#CFD8DC", "#B0BEC5", "#90A4AE", "#78909C", "#607D8B", "#546E7A", "#455A64", "#37474F", "#263238"]
	}, {
		text: "蓝灰色",
		value: "blue-gray",
		colors: ["#607D8B", "#CFD8DC", "#B0BEC5", "#90A4AE", "#78909C", "#607D8B", "#546E7A", "#455A64", "#37474F", "#263238"]
	}];

export default config;