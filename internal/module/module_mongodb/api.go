package module_mongodb

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/mongodb"
	"github.com/team-ide/go-tool/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"reflect"
	"strconv"
	"strings"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
)

type api struct {
	toolboxService *module_toolbox.ToolboxService
}

func NewApi(toolboxService *module_toolbox.ToolboxService) *api {
	return &api{
		toolboxService: toolboxService,
	}
}

var (
	Power = base.AppendPower(&base.PowerAction{Action: "mongodb", Text: "Mongodb", ShouldLogin: true, StandAlone: true})
	check = base.AppendPower(&base.PowerAction{Action: "check", Text: "Mongodb测试", ShouldLogin: true, StandAlone: true, Parent: Power})
	info  = base.AppendPower(&base.PowerAction{Action: "info", Text: "Mongodb信息", ShouldLogin: true, StandAlone: true, Parent: Power})

	database         = base.AppendPower(&base.PowerAction{Action: "database", Text: "库", ShouldLogin: true, StandAlone: true, Parent: Power})
	databaseList     = base.AppendPower(&base.PowerAction{Action: "list", Text: "列表", ShouldLogin: true, StandAlone: true, Parent: database})
	databaseDelete   = base.AppendPower(&base.PowerAction{Action: "delete", Text: "删除", ShouldLogin: true, StandAlone: true, Parent: database})
	databaseDataTrim = base.AppendPower(&base.PowerAction{Action: "dataTrim", Text: "清空数据", ShouldLogin: true, StandAlone: true, Parent: database})

	collection         = base.AppendPower(&base.PowerAction{Action: "collection", Text: "集合", ShouldLogin: true, StandAlone: true, Parent: Power})
	collectionList     = base.AppendPower(&base.PowerAction{Action: "list", Text: "列表", ShouldLogin: true, StandAlone: true, Parent: collection})
	collectionDelete   = base.AppendPower(&base.PowerAction{Action: "delete", Text: "删除", ShouldLogin: true, StandAlone: true, Parent: collection})
	collectionCreate   = base.AppendPower(&base.PowerAction{Action: "create", Text: "创建", ShouldLogin: true, StandAlone: true, Parent: collection})
	collectionDataTrim = base.AppendPower(&base.PowerAction{Action: "dataTrim", Text: "清空数据", ShouldLogin: true, StandAlone: true, Parent: collection})

	index       = base.AppendPower(&base.PowerAction{Action: "index", Text: "索引", ShouldLogin: true, StandAlone: true, Parent: Power})
	indexList   = base.AppendPower(&base.PowerAction{Action: "list", Text: "列表", ShouldLogin: true, StandAlone: true, Parent: index})
	indexDelete = base.AppendPower(&base.PowerAction{Action: "delete", Text: "删除", ShouldLogin: true, StandAlone: true, Parent: index})
	indexCreate = base.AppendPower(&base.PowerAction{Action: "create", Text: "创建", ShouldLogin: true, StandAlone: true, Parent: index})

	insert     = base.AppendPower(&base.PowerAction{Action: "insert", Text: "插入", ShouldLogin: true, StandAlone: true, Parent: Power})
	update     = base.AppendPower(&base.PowerAction{Action: "update", Text: "更新", ShouldLogin: true, StandAlone: true, Parent: Power})
	delete_    = base.AppendPower(&base.PowerAction{Action: "delete", Text: "删除", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteById = base.AppendPower(&base.PowerAction{Action: "deleteById", Text: "删除", ShouldLogin: true, StandAlone: true, Parent: Power})
	queryPage  = base.AppendPower(&base.PowerAction{Action: "queryPage", Text: "分页查询", ShouldLogin: true, StandAlone: true, Parent: Power})

	closePower = base.AppendPower(&base.PowerAction{Action: "close", Text: "关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: check, Do: this_.check})
	apis = append(apis, &base.ApiWorker{Power: info, Do: this_.info})

	apis = append(apis, &base.ApiWorker{Power: databaseList, Do: this_.databases})
	apis = append(apis, &base.ApiWorker{Power: databaseDelete, Do: this_.databaseDelete})
	apis = append(apis, &base.ApiWorker{Power: databaseDataTrim, Do: this_.databaseDataTrim})

	apis = append(apis, &base.ApiWorker{Power: collectionList, Do: this_.collections})
	apis = append(apis, &base.ApiWorker{Power: collectionCreate, Do: this_.collectionCreate})
	apis = append(apis, &base.ApiWorker{Power: collectionDelete, Do: this_.collectionDelete})
	apis = append(apis, &base.ApiWorker{Power: collectionDataTrim, Do: this_.collectionDataTrim})

	apis = append(apis, &base.ApiWorker{Power: indexList, Do: this_.indexList})
	apis = append(apis, &base.ApiWorker{Power: indexDelete, Do: this_.indexDelete})
	apis = append(apis, &base.ApiWorker{Power: indexCreate, Do: this_.indexCreate})

	apis = append(apis, &base.ApiWorker{Power: insert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Power: update, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Power: delete_, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Power: deleteById, Do: this_.deleteById})
	apis = append(apis, &base.ApiWorker{Power: queryPage, Do: this_.queryPage})

	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *mongodb.Config, err error) {
	config = &mongodb.Config{}
	_, err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(config *mongodb.Config) (res mongodb.IService, err error) {
	key := "mongodb-" + config.Address
	if config.Username != "" {
		key += "-" + base.GetMd5String(key+config.Username)
	}
	if config.Password != "" {
		key += "-" + base.GetMd5String(key+config.Password)
	}
	if config.CertPath != "" {
		key += "-" + base.GetMd5String(key+config.CertPath)
	}

	var serviceInfo *base.ServiceInfo
	serviceInfo, err = base.GetService(key, func() (res *base.ServiceInfo, err error) {
		var s mongodb.IService
		s, err = mongodb.New(config)
		if err != nil {
			util.Logger.Error("getService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}
		_, err = s.Count("_check_for_service_", "_check_for_service_", &map[string]interface{}{})
		if err != nil {
			util.Logger.Error("getService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}
		res = &base.ServiceInfo{
			WaitTime:    10 * 60 * 1000,
			LastUseTime: util.GetNowMilli(),
			Service:     s,
			Stop:        s.Close,
		}
		return
	})
	if err != nil {
		return
	}
	res = serviceInfo.Service.(mongodb.IService)
	serviceInfo.SetLastUseTime()
	return
}

type BaseRequest struct {
	WorkerId           string                 `json:"workerId"`
	DatabaseName       string                 `json:"databaseName"`
	CollectionName     string                 `json:"collectionName"`
	IndexName          string                 `json:"indexName"`
	Keys               bson.D                 `json:"keys"`
	Filter             map[string]interface{} `json:"filter"`
	WhereDoc           string                 `json:"whereDoc"`
	WhereList          []*Where               `json:"whereList"`
	OrderList          []*Order               `json:"orderList"`
	Id                 string                 `json:"id"`
	IdType             string                 `json:"idType"`
	PageIndex          int64                  `json:"pageIndex"`
	PageSize           int64                  `json:"pageSize"`
	Doc                string                 `json:"doc"`
	IsObjectID         bool                   `json:"isObjectID"`
	ObjectIDKey        string                 `json:"objectIDKey"`
	IndexType          string                 `json:"indexType"`
	ExpireAfterSeconds int32                  `json:"expireAfterSeconds"`
}

type Where struct {
	Name                 string `json:"name"`
	Value                string `json:"value"`
	Before               string `json:"before"`
	After                string `json:"after"`
	CustomSql            string `json:"customSql"`
	ConditionalOperation string `json:"conditionalOperation"`
	AndOr                string `json:"andOr"`
	DataType             string `json:"dataType"`
}

type Order struct {
	Name    string `json:"name"`
	AscDesc string `json:"ascDesc"`
}

func (this_ *api) check(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	_, err = getService(config)
	if err != nil {
		return
	}

	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}
func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	_, err = getService(config)
	if err != nil {
		return
	}

	return
}

func (this_ *api) databases(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	res, _, err = service.Databases()
	if err != nil {
		return
	}

	return
}

func (this_ *api) databaseDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.DatabaseDelete(request.DatabaseName)
	if err != nil {
		return
	}

	return
}

func (this_ *api) databaseDataTrim(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	collections, err := service.Collections(request.DatabaseName)
	if err != nil {
		return
	}

	for _, one := range collections {
		_, err = service.DeleteMany(request.DatabaseName, one.Name, bson.M{})
		if err != nil {
			return
		}
	}

	return
}

func (this_ *api) collections(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Collections(request.DatabaseName)
	if err != nil {
		return
	}

	return
}

func (this_ *api) collectionCreate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.CollectionCreate(request.DatabaseName, request.CollectionName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) collectionDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.CollectionDelete(request.DatabaseName, request.CollectionName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) collectionDataTrim(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	_, err = service.DeleteMany(request.DatabaseName, request.CollectionName, bson.M{})
	if err != nil {
		return
	}

	return
}

func (this_ *api) indexList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Indexes(request.DatabaseName, request.CollectionName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) indexDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.IndexDelete(request.DatabaseName, request.CollectionName, request.IndexName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) indexCreate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	i := mongo.IndexModel{}
	i.Keys = request.Keys
	i.Options = &options.IndexOptions{}
	if request.IndexName != "" {
		i.Options.SetName(request.IndexName)
	}
	if request.IndexType == "unique" {
		i.Options.SetUnique(true)
	} else if request.IndexType == "expireAfterSeconds" {
		i.Options.SetExpireAfterSeconds(request.ExpireAfterSeconds)
	}

	res, err = service.IndexCreate(request.DatabaseName, request.CollectionName, i)
	if err != nil {
		return
	}
	return
}

func (this_ *api) insert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	data := map[string]interface{}{}
	err = util.JSONDecodeUseNumber([]byte(request.Doc), &data)
	if err != nil {
		return
	}
	res, err = service.Insert(request.DatabaseName, request.CollectionName, data)
	if err != nil {
		return
	}
	return
}

func (this_ *BaseRequest) getID() interface{} {
	switch this_.IdType {
	case "M", "map":
		data := map[string]interface{}{}
		err := util.JSONDecodeUseNumber([]byte(this_.Id), &data)
		if err == nil {
			return data
		}
		break
	case "primitive.ObjectID", "ObjectID":
		if v, e := primitive.ObjectIDFromHex(this_.Id); e == nil {
			return v
		}
		break
	case "int64", "int32", "int16", "int8":
		if v, e := strconv.ParseInt(this_.Id, 10, 64); e == nil {
			switch this_.IdType {
			case "int32":
				return int32(v)
			case "int16":
				return int16(v)
			case "int8":
				return int8(v)
			default:
				return v
			}
		}
		break
	case "float64", "float32":
		if v, e := strconv.ParseFloat(this_.Id, 64); e == nil {
			switch this_.IdType {
			case "float32":
				return float32(v)
			default:
				return v
			}
		}
		break
	}
	return this_.Id
}

func (this_ *api) update(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	var id = request.getID()
	data := map[string]interface{}{}
	err = util.JSONDecodeUseNumber([]byte(request.Doc), &data)
	if err != nil {
		return
	}

	delete(data, "_id")

	res, err = service.Update(request.DatabaseName, request.CollectionName, id, bson.M{
		"$set": data,
	})
	if err != nil {
		return
	}
	return
}
func getNumberValue(s string) (interface{}, error) {
	n := json.Number(s)
	if strings.Contains(s, ".") {
		return n.Float64()
	}
	return n.Int64()
}
func (this_ *api) queryPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.IsObjectID && request.ObjectIDKey != "" && request.Filter[request.ObjectIDKey] != nil {
		v, e := primitive.ObjectIDFromHex(util.GetStringValue(request.Filter[request.ObjectIDKey]))
		if e == nil {
			request.Filter[request.ObjectIDKey] = v
		}
	}

	page := &mongodb.Page{
		PageSize: request.PageSize,
		PageNo:   request.PageIndex,
	}

	filter := bson.M{}
	if request.WhereDoc != "" {
		err = util.JSONDecodeUseNumber([]byte(request.WhereDoc), &filter)
		if err != nil {
			return
		}
	} else {

		for _, where := range request.WhereList {
			if where.Name == "" {
				continue
			}
			dataType := where.DataType
			var v interface{}
			v = where.Value
			if dataType == "number" && util.StringIndexOf([]string{
				"=", "<>", "<", "<=", ">", ">=",
			}, where.ConditionalOperation) >= 0 {
				v, err = getNumberValue(where.Value)
				if err != nil {
					return
				}
			}
			if where.ConditionalOperation == "=" && dataType == "objectID" {
				if v, err = primitive.ObjectIDFromHex(where.Value); err != nil {
					return
				}
			}
			filter[where.Name] = v
			switch where.ConditionalOperation {
			case "like":
				filter[where.Name] = bson.M{
					"$regex": where.Value,
				}
				break
			case "like start":
				filter[where.Name] = bson.M{
					"$regex": "^" + where.Value,
				}
				break
			case "like end":
				filter[where.Name] = bson.M{
					"$regex": where.Value + "$",
				}
				break
			case "<>":
				filter[where.Name] = bson.M{
					"$ne": v,
				}
				break
			case ">":
				filter[where.Name] = bson.M{
					"$gt": v,
				}
				break
			case ">=":
				filter[where.Name] = bson.M{
					"$gte": v,
				}
				break
			case "<":
				filter[where.Name] = bson.M{
					"$lt": v,
				}
				break
			case "<=":
				filter[where.Name] = bson.M{
					"$lte": v,
				}
				break
			case "between":
				var b interface{} = where.Before
				var a interface{} = where.After
				if dataType == "number" {
					b, err = getNumberValue(where.Before)
					if err != nil {
						return
					}
					a, err = getNumberValue(where.After)
					if err != nil {
						return
					}
				}
				filter[where.Name] = bson.M{
					"$gte": b,
					"$lte": a,
				}
				break
			case "in":
				var vs []interface{}
				ss := strings.Split(where.Value, ",")
				for _, s := range ss {
					if dataType == "number" {
						var n interface{}
						n, err = getNumberValue(where.Before)
						if err != nil {
							return
						}
						vs = append(vs, n)
					} else {
						vs = append(vs, s)
					}
				}
				filter[where.Name] = bson.M{
					"$in": vs,
				}
				break
			}
		}
	}
	sort := bson.D{}
	for _, order := range request.OrderList {
		if order.Name == "" {
			continue
		}
		if strings.EqualFold(order.AscDesc, "asc") {
			sort = append(sort, bson.E{
				Key:   order.Name,
				Value: 1,
			})
		} else {
			sort = append(sort, bson.E{
				Key:   order.Name,
				Value: -1,
			})
		}
	}

	opts := options.Find()
	opts.SetSort(sort)
	result, err := service.QueryMapPageResult(request.DatabaseName, request.CollectionName, filter, page, opts)
	if err != nil {
		return
	}

	res = result
	var list []interface{}

	var bs []byte
	for _, one := range result.List {
		item := one.(map[string]interface{})
		d := map[string]interface{}{}
		bs, err = json.MarshalIndent(item, "", "  ")
		if err != nil {
			return
		}
		if _id, ok := item["_id"]; ok {
			typeName := reflect.TypeOf(_id).Name()
			if typeName == "ObjectID" {
				id := _id.(primitive.ObjectID)
				d["_id"] = id.Hex()
			} else {
				d["_id"] = util.GetStringValue(_id)
			}
			d["_id_type"] = typeName
		}
		d["value"] = string(bs)
		list = append(list, d)
	}
	result.List = list
	return
}

func (this_ *api) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.IsObjectID && request.ObjectIDKey != "" && request.Filter[request.ObjectIDKey] != nil {
		v, e := primitive.ObjectIDFromHex(util.GetStringValue(request.Filter[request.ObjectIDKey]))
		if e == nil {
			request.Filter[request.ObjectIDKey] = v
		}
	}
	res, err = service.DeleteOne(request.DatabaseName, request.CollectionName, request.Filter)
	if err != nil {
		return
	}
	return
}

func (this_ *api) deleteById(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.DeleteOne(request.DatabaseName, request.CollectionName, bson.M{
		"_id": request.getID(),
	})
	if err != nil {
		return
	}
	return
}
