let distanceCSS = `
/** 字体大小 **/
.ft-key,.hover-ft-key:hover{font-size: "value" !important;}
/** 宽度 **/
.wd-key,.hover-wd-key:hover{width: "value" !important;}
/** 高度 **/
.hg-key,.hover-hg-key:hover{height: "value" !important;}
/** 圆角 **/
.rd-key,.hover-rd-key:hover,.rdt-key,.hover-rdt-key:hover,.rdtl-key,.hover-rdtl-key:hover{border-top-left-radius: "value" !important;}
.rd-key,.hover-rd-key:hover,.rdt-key,.hover-rdt-key:hover,.rdtr-key,.hover-rdtr-key:hover{border-top-right-radius: "value" !important;}
.rd-key,.hover-rd-key:hover,.rdb-key,.hover-rdb-key:hover,.rdbl-key,.hover-rdbl-key:hover{border-bottom-left-radius: "value" !important;}
.rd-key,.hover-rd-key:hover,.rdb-key,.hover-rdb-key:hover,.rdbr-key,.hover-rdbr-key:hover{border-bottom-right-radius: "value" !important;}
/** 边框 **/
.bd-key,.bdl-key,.bdlr-key,.hover-bd-key:hover,.hover-bdl-key:hover,.hover-bdlr-key:hover{border-left-style: solid;border-left-width: "value" !important;}
.bd-key,.bdr-key,.bdlr-key,.hover-bd-key:hover,.hover-bdr-key:hover,.hover-bdlr-key:hover{border-right-style: solid;border-right-width: "value" !important;}
.bd-key,.bdt-key,.bdtb-key,.hover-bd-key:hover,.hover-bdt-key:hover,.hover-bdtb-key:hover{border-top-style: solid;border-top-width: "value" !important;}
.bd-key,.bdb-key,.bdtb-key,.hover-bd-key:hover,.hover-bdb-key:hover,.hover-bdtb-key:hover{border-bottom-style: solid;border-bottom-width: "value" !important;}
/** 内边距 **/
.pd-key,.pdl-key,.pdlr-key,.hover-pd-key:hover,.hover-pdl-key:hover,.hover-pdlr-key:hover{padding-left: "value" !important;}
.pd-key,.pdr-key,.pdlr-key,.hover-pd-key:hover,.hover-pdr-key:hover,.hover-pdlr-key:hover{padding-right: "value" !important;}
.pd-key,.pdt-key,.pdtb-key,.hover-pd-key:hover,.hover-pdt-key:hover,.hover-pdtb-key:hover{padding-top: "value" !important;}
.pd-key,.pdb-key,.pdtb-key,.hover-pd-key:hover,.hover-pdb-key:hover,.hover-pdtb-key:hover{padding-bottom: "value" !important;}
/** 外边距 **/
.mg-key,.mgl-key,.mglr-key,.hover-mg-key:hover,.hover-mgl-key:hover,.hover-mglr-key:hover{margin-left: "value" !important;}
.mg-key,.mgr-key,.mglr-key,.hover-mg-key:hover,.hover-mgr-key:hover,.hover-mglr-key:hover{margin-right: "value" !important;}
.mg-key,.mgt-key,.mgtb-key,.hover-mg-key:hover,.hover-mgt-key:hover,.hover-mgtb-key:hover{margin-top: "value" !important;}
.mg-key,.mgb-key,.mgtb-key,.hover-mg-key:hover,.hover-mgb-key:hover,.hover-mgtb-key:hover{margin-bottom: "value" !important;}

.mg--key,.mgl--key,.mglr--key,.hover-mg--key:hover,.hover-mgl--key:hover,.hover-mglr--key:hover{margin-left: -"value" !important;}
.mg--key,.mgr--key,.mglr--key,.hover-mg--key:hover,.hover-mgr--key:hover,.hover-mglr--key:hover{margin-right: -"value" !important;}
.mg--key,.mgt--key,.mgtb--key,.hover-mg--key:hover,.hover-mgt--key:hover,.hover-mgtb--key:hover{margin-top: -"value" !important;}
.mg--key,.mgb--key,.mgtb--key,.hover-mg--key:hover,.hover-mgb--key:hover,.hover-mgtb--key:hover{margin-bottom: -"value" !important;}
	`;
let colCSS = `
.col-key,.tm-window-xs .col-xs-key,.tm-window-sm .col-sm-key,.tm-window-md .col-md-key,.tm-window-lg .col-lg-key{float: left;width: "value" !important;}
.offset-key,.tm-window-xs .offset-xs-key,.tm-window-sm .offset-sm-key,.tm-window-md .offset-md-key,.tm-window-lg .offset-lg-key{margin-left: "value" !important;}
    `;
let sizeCSS = `
.font-key{font-size: "font-size";}

.tm-btn-key{line-height: "line-height";padding: "padding";font-size: "font-size";}
.tm-btn-key *{height: "line-height";display: inline-block;vertical-align: bottom;}
.tm-btn-key img{border-radius: 100px;}
.tm-btn-circle.tm-btn-key{height: "height";width: "height";padding:0px;line-height: "circle-line-height";}

`;
let colorCSS = `
/** color **/
.color-key,.active-color-key.tm-active,.tm-window-xs .color-xs-key,.tm-window-sm .color-sm-key,.tm-window-md .color-md-key,.tm-window-lg .color-lg-key{color: "color" !important;}
@media (hover: hover) {
	.active-color-key:hover{color: "color" !important;}
	.hover-color-key:hover{color: "color" !important;}
}
/** background color **/
.bg-key,.active-bg-key.tm-active,.tm-window-xs .bg-xs-key,.tm-window-sm .bg-sm-key,.tm-window-md .bg-md-key,.tm-window-lg .bg-lg-key{color: "whiteColor";background-color: "color" !important;}
@media (hover: hover) {
	.active-bg-key:hover{ color: "whiteColor";background-color: "color" !important;}
	.hover-bg-key:hover{ color: "whiteColor";background-color: "color" !important;}
}
/** border color **/
.bd-key,.hover-bd-key:hover{border-color: "color" !important;}
.bdl-key,.bdlr-key,.hover-bdl-key:hover,.hover-bdlr-key:hover{border-left-color: "color" !important;}
.bdr-key,.bdlr-key,.hover-bdr-key:hover,.hover-bdlr-key:hover{border-right-color: "color" !important;}
.bdt-key,.bdtb-key,.hover-bdt-key:hover,.hover-bdtb-key:hover{border-top-color: "color" !important;}
.bdb-key,.bdtb-key,.hover-bdb-key:hover,.hover-bdtb-key:hover{border-bottom-color: "color" !important;}

`;
let componentCSS = `
/** btn color **/
.tm-btn.color-key,
.tm-btn.bd-key,
.tm-btn-color-key
{color: "color" !important;border-color: "color" !important;} /**background-color: #FFFFFF !important;**/
/** btn background color **/
.tm-btn.bg-key,
.tm-btn-bg-key{color: #FFFFFF !important;border-color: "color" !important;background-color: "color" !important;}


`;
let componentActiveCSS = `
/** btn color **/
.tm-btn.color-key.tm-active,.tm-btn.color-key.tm-hover,
.tm-btn.bd-key.tm-active,.tm-btn.bd-key.tm-hover,
.tm-btn-color-key.tm-active,.tm-btn-color-key.tm-hover
{color: #FFFFFF !important;border-color: "color" !important;background-color: "color" !important;}
@media (hover: hover) {
	.tm-btn.color-key:hover,
	.tm-btn.bd-key:hover,
	.tm-btn-color-key:hover
	{color: #FFFFFF !important;border-color: "color" !important;background-color: "color" !important;}
}
/** btn background color **/
.tm-btn.bg-key.tm-active,.tm-btn.bg-key.tm-hover,
.tm-btn-bg-key.tm-active,.tm-btn-bg-key.tm-hover
{border-color: "activeColor" !important;background-color: "activeColor" !important;}
@media (hover: hover) {
	.tm-btn.bg-key:hover,
	.tm-btn-bg-key:hover
	{border-color: "activeColor" !important;background-color: "activeColor" !important;}
}
.tm-btn.active-color-key.tm-active,.tm-btn.active-color-key.tm-hover,
.tm-btn-active-color-key.tm-active,.tm-btn-active-color-key.tm-hover
{background-color: #FFFFFF !important;border-color: "color" !important;color: "color" !important;}
@media (hover: hover) {
	.tm-btn.active-color-key:hover,
	.tm-btn-active-color-key:hover
	{background-color: #FFFFFF !important;border-color: "color" !important;color: "color" !important;}
}
.tm-btn.active-bg-key.tm-active,.tm-btn.active-bg-key.tm-hover,
.tm-btn-active-bg-key.tm-active,.tm-btn-active-bg-key.tm-hover
{color: #FFFFFF !important;background-color: "color" !important;border-color: "color" !important;}
@media (hover: hover) {
	.tm-btn.active-bg-key:hover,
	.tm-btn-active-bg-key:hover
	{color: #FFFFFF !important;background-color: "color" !important;border-color: "color" !important;}
}

/** link color **/
.tm-link.color-key.tm-active,.tm-link.color-key.tm-hover,
.tm-link-color-key.tm-active,.tm-link-color-key.tm-hover
{color: "color";background-color: transparent;border-bottom-color: "color";}
@media (hover: hover) {
	.tm-link.color-key:hover,
	.tm-link-color-key:hover
	{color: "color";background-color: transparent;border-bottom-color: "color";}
}
.tm-link.active-color-key.tm-active,.tm-link.active-color-key.tm-hover,
.tm-link-active-color-key.tm-active,.tm-link-active-color-key.tm-hover
{color: "color" !important;background-color: transparent !important;border-bottom-color: "color" !important;}
@media (hover: hover) {
	.tm-link.active-color-key:hover,
	.tm-link-active-color-key:hover
	{color: "color" !important;background-color: transparent !important;border-bottom-color: "color" !important;}
}

`;

export default { distanceCSS, colCSS, sizeCSS, colorCSS, componentCSS, componentActiveCSS }