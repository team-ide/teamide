package common

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
	base2 "teamide/pkg/application/base"
	"teamide/pkg/application/model"
	"time"

	"github.com/Chain-Zhang/pinyin"
	uuid "github.com/satori/go.uuid"
)

type IScript interface {
	TrimSuffix(s, suffix string) (res string)
	TrimPrefix(s, prefix string) (res string)
	GetLock(key string) (res Locker, err error)

	Length(obj interface{}) int                       //
	NotMatch(pattern string, value interface{}) bool  //
	Match(pattern string, value interface{}) bool     //
	ThrowError(code string, msg string) (bool, error) //

	IsEmpty(obj interface{}) bool    // 是否为空  null或空字符串
	IsNotEmpty(obj interface{}) bool // 是否不为空  取 是否为空 反

	IsTrue(obj interface{}) bool  // 是否为真 入：true、"true"、1、"1"为真
	IsFalse(obj interface{}) bool // 是否为假  取 是否为真 反

	GetErrCode(err error) string // 获取异常中的错误码
	GetErrMsg(err error) string  // 获取异常中的信息

	WrapTableName(database string, table string) string // 包装表名
	WrapColumnName(alias string, column string) string  // 包装列名

	DataToJSON(data interface{}) string
	JSONToData(json string) (interface{}, error)

	AesECBEncrypt(data, key string) (string, error) // AES ECB模式加密
	AesECBDecrypt(data, key string) (string, error) // AES ECB模式解密

	AesCBCEncrypt(origData, key string) (string, error) // AES CBC模式加密
	AesCBCDecrypt(origData, key string) (string, error) // AES CBC模式解密

	Base64Encode(data string) (string, error)     // Base64 加密
	Base64Decode(origData string) (string, error) // Base64 解密

	MD5(data string) (string, error) // MD5 加密

	UUID() (string, error) // MD5 加密
	GetReader(obj interface{}) (res io.Reader, err error)
	GetFileTypeSuffix(obj interface{}) (res string)

	Now() time.Time                                  // 当前时间
	NowFormat(format string) string                  // 当前时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS
	NowTime() int64                                  // 当前时间戳
	DateFormat(date time.Time, format string) string // 获取传入时间的时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS
	DateTime(ate time.Time) int64                    // 获取传入时间的时间戳
	ToDate(time string, format string) time.Time     // 根据格式化的日期和格式转为日期
	ToDateByTime(time int64) time.Time               // 根据时间戳转为日期

	GetID() int64                    // 生成ID，生成不重复的ID
	GetIDByType(idType string) int64 // 根据ID类型生成ID，生成不重复的ID

	GetStringID() string                    // 生成字符串ID，生成不重复的ID
	GetStringIDByType(idType string) string // 根据ID类型生成字符串ID，生成不重复的ID

	RandString(minLen int, maxLen int) string   // 生成随机字符串
	RandUserName(minLen int, maxLen int) string // 生成随机用户名
	ToPinYin(name string) (string, error)       // 转换为拼音
}

type ScriptDefault struct {
}

func (this_ *ScriptDefault) Length(obj interface{}) (res int) { //
	if obj == nil {
		return 0
	}
	switch v := obj.(type) {
	case string:
		return len(v)
	case map[interface{}]interface{}:
		return len(v)
	case []interface{}:
		return len(v)
	default:
		return len(fmt.Sprint(v))
	}
}
func (this_ *ScriptDefault) TrimPrefix(s, prefix string) (res string) { //
	res = strings.TrimPrefix(s, prefix)
	return
}

func (this_ *ScriptDefault) TrimSuffix(s, suffix string) (res string) { //
	res = strings.TrimSuffix(s, suffix)
	return
}

//
func (this_ *ScriptDefault) GetLock(key string) (res Locker, err error) { //

	return
}

//
func (this_ *ScriptDefault) ThrowError(code string, msg string) (res bool, err error) { //
	err = base2.NewError(code, msg)
	return
}
func (this_ *ScriptDefault) Match(pattern string, value interface{}) (res bool) { //
	var re *regexp.Regexp = regexp.MustCompile(pattern)
	res = re.MatchString(fmt.Sprint(value))
	return
}
func (this_ *ScriptDefault) NotMatch(pattern string, value interface{}) (res bool) { //
	res = !this_.Match(pattern, value)
	return
}

func (this_ *ScriptDefault) IsEmpty(obj interface{}) (res bool) { // 是否为空  null或空字符串
	if obj == nil || obj == "" || obj == 0 {
		res = true
	}
	return
}

func (this_ *ScriptDefault) IsNotEmpty(obj interface{}) (res bool) { // 是否不为空  取 是否为空 反
	return !this_.IsEmpty(obj)
}

func (this_ *ScriptDefault) IsTrue(obj interface{}) (res bool) { // 是否为真 入：true、"true"、1、"1"为真
	if obj == true || obj == "true" || obj == 1 || obj == "1" {
		res = true
	}
	return
}

func (this_ *ScriptDefault) IsFalse(obj interface{}) (res bool) { // 是否为假  取 是否为真 反
	return !this_.IsTrue(obj)
}

func (this_ *ScriptDefault) GetErrCode(err error) (res string) { // 获取异常中的错误码
	res = "-1"
	return
}

func (this_ *ScriptDefault) GetErrMsg(err error) (res string) { // 获取异常中的信息
	res = err.Error()
	return
}

func (this_ *ScriptDefault) WrapTableName(database string, table string) (res string) { // 包装表名
	if this_.IsNotEmpty(database) {
		res = "" + database + "."
	}
	res += "" + table + ""
	return
}
func (this_ *ScriptDefault) WrapColumnName(tableAlias string, column string) (res string) { // 包装列名
	if this_.IsNotEmpty(tableAlias) {
		res = "" + tableAlias + "."
	}
	res += "" + column + ""
	return
}

func (this_ *ScriptDefault) DataToJSON(data interface{}) (res string) {
	if data == nil {
		return
	}
	res = base2.ToJSON(data)
	return
}
func (this_ *ScriptDefault) JSONToData(json string) (data interface{}, err error) {
	if json == "" {
		return
	}
	data = make(map[string]interface{})
	err = base2.ToBean([]byte(json), &data)
	return
}

func (this_ *ScriptDefault) AesECBEncrypt(data, key string) (origData string, err error) { // AES ECB模式加密
	var bs []byte
	bs, err = base2.AesECBEncrypt([]byte(data), []byte(key))
	if err != nil {
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	origData = base64.StdEncoding.EncodeToString(bs)
	return
}

func (this_ *ScriptDefault) AesECBDecrypt(origData, key string) (data string, err error) { // AES ECB模式解密
	var bs []byte
	// 经过一次base64 否则 直接转字符串乱码
	bs, err = base64.StdEncoding.DecodeString(origData)
	if err != nil {
		return
	}
	bs, err = base2.AesECBDecrypt(bs, []byte(key))
	if err != nil {
		return
	}
	data = string(bs)
	return
}

func (this_ *ScriptDefault) AesCBCEncrypt(data, key string) (origData string, err error) { // AES CBC模式加密
	var bs []byte
	bs, err = base2.AesCBCEncrypt([]byte(data), []byte(key))
	if err != nil {
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	origData = base64.StdEncoding.EncodeToString(bs)
	return
}
func (this_ *ScriptDefault) AesCBCDecrypt(origData, key string) (data string, err error) { // AES CBC模式解密
	var bs []byte
	// 经过一次base64 否则 直接转字符串乱码
	bs, err = base64.StdEncoding.DecodeString(origData)
	if err != nil {
		return
	}
	bs, err = base2.AesCBCDecrypt(bs, []byte(key))
	if err != nil {
		return
	}
	data = string(bs)
	return
}
func (this_ *ScriptDefault) Base64Encode(data string) (origData string, err error) { // Base64 加密

	return
}
func (this_ *ScriptDefault) Base64Decode(origData string) (data string, err error) { // Base64 解密

	return
}

func (this_ *ScriptDefault) MD5(data string) (origData string, err error) { // MD5 加密
	m := md5.New()
	_, err = io.WriteString(m, data)
	if err != nil {
		return
	}
	arr := m.Sum(nil)
	origData = fmt.Sprintf("%x", arr)
	return
}

func (this_ *ScriptDefault) GetReader(obj interface{}) (res io.Reader, err error) { // MD5 加密
	if obj == nil {
		return
	}
	switch v := obj.(type) {
	case io.Reader:
		res = v
	case *model.FileInfo:
		res, err = v.GetReader()
	case map[string]interface{}:
		fileInfo := &model.FileInfo{}
		if v["name"] != nil {
			fileInfo.Name = v["name"].(string)
		}
		if v["path"] != nil {
			fileInfo.Path = v["path"].(string)
		}
		if v["absolutePath"] != nil {
			fileInfo.AbsolutePath = v["absolutePath"].(string)
		}
		if v["fileHeader"] != nil {
			fileInfo.FileHeader = v["fileHeader"].(*multipart.FileHeader)
		}
		if v["file"] != nil {
			fileInfo.File = v["file"].(*os.File)
		}
		res, err = fileInfo.GetReader()
	}
	return
}

func (this_ *ScriptDefault) GetFileTypeSuffix(obj interface{}) (res string) { // MD5 加密
	if obj == nil {
		return
	}

	var fileName string
	var fileType string
	switch v := obj.(type) {
	case *model.FileInfo:
		fileName = v.Name
		fileType = v.Type
	case map[string]interface{}:
		if v["name"] != nil {
			fileName = v["name"].(string)
		}
		if v["type"] != nil {
			fileType = v["type"].(string)
		}
	}

	if fileType == "" && strings.Contains(fileName, ".") {
		fileType = fileName[strings.LastIndex(fileName, ".")+1:]
	}
	if fileType != "" {
		res = "." + fileType
	}
	return
}

func (this_ *ScriptDefault) UUID() (res string, err error) { // MD5 加密
	uuid := uuid.NewV4().String()
	slice := strings.Split(uuid, "-")
	var uuidNew string
	for _, str := range slice {
		uuidNew += str
	}
	res = uuidNew
	return
}
func (this_ *ScriptDefault) Now() (res time.Time) { // 当前时间
	res = time.Now()
	return
}
func (this_ *ScriptDefault) NowFormat(format string) (res string) { // 当前时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS
	if format == "" {
		format = "2006-01-02 15:04:05.000"
	}
	res = this_.Now().Format(format)
	return
}
func (this_ *ScriptDefault) NowTime() (res int64) { // 当前时间戳
	res = time.Now().UnixNano() / 1e6
	return
}
func (this_ *ScriptDefault) DateFormat(date time.Time, format string) (res string) { // 获取传入时间的时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS

	return
}
func (this_ *ScriptDefault) DateTime(date time.Time) (res int64) { // 获取传入时间的时间戳

	return
}
func (this_ *ScriptDefault) ToDate(time string, format string) (res time.Time) { // 根据格式化的日期和格式转为日期

	return
}
func (this_ *ScriptDefault) ToDateByTime(time int64) (res time.Time) { // 根据时间戳转为日期

	return
}

func (this_ *ScriptDefault) GetID() (res int64) { // 生成ID，生成不重复的ID
	res = int64(randInt(1, 999999999))
	return
}
func (this_ *ScriptDefault) GetIDByType(idType string) (res int64) { // 根据ID类型生成ID，生成不重复的ID

	return
}

func (this_ *ScriptDefault) GetStringID() (res string) { // 生成字符串ID，生成不重复的ID

	return
}
func (this_ *ScriptDefault) GetStringIDByType(idType string) (res string) { // 根据ID类型生成字符串ID，生成不重复的ID

	return
}

func (this_ *ScriptDefault) RandString(minLen int, maxLen int) string { // 生成随机字符串
	len := minLen
	if maxLen > minLen {
		len = randInt(minLen, maxLen)
	}
	var i int = 0
	var str = ""
	for i = 0; i < len; i++ {
		randNum := randInt(0, randChatsSize*3)
		str += randChats[randNum%randChatsSize]

	}
	return str
}

func (this_ *ScriptDefault) RandUserName(minLen int, maxLen int) string { // 生成随机用户名
	len := minLen - 1
	if maxLen > minLen {
		len = randInt(minLen-1, maxLen-1)
	}
	str := firstName[randInt(0, firstNameLen*3)%firstNameLen]
	for i := 0; i < len; i++ { //随机产生2位或者3位的名
		str += lastName[randInt(0, lastNameLen*3+i)%lastNameLen]
	}
	return str
}

func (this_ *ScriptDefault) ToPinYin(name string) (string, error) { // 转换为拼音
	str, err := pinyin.New(name).Split("").Mode(pinyin.WithoutTone).Convert()
	if err != nil {
		// 错误处理
		return "", err
	}
	return str, nil
}

var (
	// randMutex sync.Mutex
	//设置随机数种子
	rand_ = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 获取随机数
func randInt(min int, max int) int {
	if max <= min {
		return min
	}
	// randMutex.Lock()
	// defer randMutex.Unlock()
	return min + rand_.Intn(max-min)
}

var (
	randChats = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u",
		"v", "w", "z", "y", "z",
		"A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U",
		"V", "W", "Z", "Y", "Z",
		"_",
	}
	randChatsSize = len(randChats)

	firstName = []string{
		"赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈", "褚", "卫", "蒋",
		"沈", "韩", "杨", "朱", "秦", "尤", "许", "何", "吕", "施", "张", "孔", "曹", "严", "华", "金", "魏",
		"陶", "姜", "戚", "谢", "邹", "喻", "柏", "水", "窦", "章", "云", "苏", "潘", "葛", "奚", "范", "彭",
		"郎", "鲁", "韦", "昌", "马", "苗", "凤", "花", "方", "任", "袁", "柳", "鲍", "史", "唐", "费", "薛",
		"雷", "贺", "倪", "汤", "滕", "殷", "罗", "毕", "郝", "安", "常", "傅", "卞", "齐", "元", "顾", "孟",
		"平", "黄", "穆", "萧", "尹", "姚", "邵", "湛", "汪", "祁", "毛", "狄", "米", "伏", "成", "戴", "谈",
		"宋", "茅", "庞", "熊", "纪", "舒", "屈", "项", "祝", "董", "梁", "杜", "阮", "蓝", "闵", "季", "贾",
		"路", "娄", "江", "童", "颜", "郭", "梅", "盛", "林", "钟", "徐", "邱", "骆", "高", "夏", "蔡", "田",
		"樊", "胡", "凌", "霍", "虞", "万", "支", "柯", "管", "卢", "莫", "柯", "房", "裘", "缪", "解", "应",
		"宗", "丁", "宣", "邓", "单", "杭", "洪", "包", "诸", "左", "石", "崔", "吉", "龚", "程", "嵇", "邢",
		"裴", "陆", "荣", "翁", "荀", "于", "惠", "甄", "曲", "封", "储", "仲", "伊", "宁", "仇", "甘", "武",
		"符", "刘", "景", "詹", "龙", "叶", "幸", "司", "黎", "溥", "印", "怀", "蒲", "邰", "从", "索", "赖",
		"卓", "屠", "池", "乔", "胥", "闻", "莘", "党", "翟", "谭", "贡", "劳", "逄", "姬", "申", "扶", "堵",
		"冉", "宰", "雍", "桑", "寿", "通", "燕", "浦", "尚", "农", "温", "别", "庄", "晏", "柴", "瞿", "阎",
		"连", "习", "容", "向", "古", "易", "廖", "庾", "终", "步", "都", "耿", "满", "弘", "匡", "国", "文",
		"寇", "广", "禄", "阙", "东", "欧", "利", "师", "巩", "聂", "关", "荆", "司马", "上官", "欧阳", "夏侯",
		"诸葛", "闻人", "东方", "赫连", "皇甫", "尉迟", "公羊", "澹台", "公冶", "宗政", "濮阳", "淳于", "单于",
		"太叔", "申屠", "公孙", "仲孙", "轩辕", "令狐", "徐离", "宇文", "长孙", "慕容", "司徒", "司空"}
	lastName = []string{
		"伟", "刚", "勇", "毅", "俊", "峰", "强", "军", "平", "保", "东", "文", "辉", "力", "明", "永", "健", "世", "广", "志", "义",
		"兴", "良", "海", "山", "仁", "波", "宁", "贵", "福", "生", "龙", "元", "全", "国", "胜", "学", "祥", "才", "发", "武", "新",
		"利", "清", "飞", "彬", "富", "顺", "信", "子", "杰", "涛", "昌", "成", "康", "星", "光", "天", "达", "安", "岩", "中", "茂",
		"进", "林", "有", "坚", "和", "彪", "博", "诚", "先", "敬", "震", "振", "壮", "会", "思", "群", "豪", "心", "邦", "承", "乐",
		"绍", "功", "松", "善", "厚", "庆", "磊", "民", "友", "裕", "河", "哲", "江", "超", "浩", "亮", "政", "谦", "亨", "奇", "固",
		"之", "轮", "翰", "朗", "伯", "宏", "言", "若", "鸣", "朋", "斌", "梁", "栋", "维", "启", "克", "伦", "翔", "旭", "鹏", "泽",
		"晨", "辰", "士", "以", "建", "家", "致", "树", "炎", "德", "行", "时", "泰", "盛", "雄", "琛", "钧", "冠", "策", "腾", "楠",
		"榕", "风", "航", "弘", "秀", "娟", "英", "华", "慧", "巧", "美", "娜", "静", "淑", "惠", "珠", "翠", "雅", "芝", "玉", "萍",
		"红", "娥", "玲", "芬", "芳", "燕", "彩", "春", "菊", "兰", "凤", "洁", "梅", "琳", "素", "云", "莲", "真", "环", "雪", "荣",
		"爱", "妹", "霞", "香", "月", "莺", "媛", "艳", "瑞", "凡", "佳", "嘉", "琼", "勤", "珍", "贞", "莉", "桂", "娣", "叶", "璧",
		"璐", "娅", "琦", "晶", "妍", "茜", "秋", "珊", "莎", "锦", "黛", "青", "倩", "婷", "姣", "婉", "娴", "瑾", "颖", "露", "瑶",
		"怡", "婵", "雁", "蓓", "纨", "仪", "荷", "丹", "蓉", "眉", "君", "琴", "蕊", "薇", "菁", "梦", "岚", "苑", "婕", "馨", "瑗",
		"琰", "韵", "融", "园", "艺", "咏", "卿", "聪", "澜", "纯", "毓", "悦", "昭", "冰", "爽", "琬", "茗", "羽", "希", "欣", "飘",
		"育", "滢", "馥", "筠", "柔", "竹", "霭", "凝", "晓", "欢", "霄", "枫", "芸", "菲", "寒", "伊", "亚", "宜", "可", "姬", "舒",
		"影", "荔", "枝", "丽", "阳", "妮", "宝", "贝", "初", "程", "梵", "罡", "恒", "鸿", "桦", "骅", "剑", "娇", "纪", "宽", "苛",
		"灵", "玛", "媚", "琪", "晴", "容", "睿", "烁", "堂", "唯", "威", "韦", "雯", "苇", "萱", "阅", "彦", "宇", "雨", "洋", "忠",
		"宗", "曼", "紫", "逸", "贤", "蝶", "菡", "绿", "蓝", "儿", "翠", "烟", "小", "轩"}
	firstNameLen = len(firstName)
	lastNameLen  = len(lastName)
)
