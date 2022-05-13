package zorm

//Page 分页对象
//Page Pagination object
type Page struct {
	//当前页码,从1开始
	//Current page number, starting from 1
	PageNo int

	//每页多少条,默认20条
	//How many items per page, 20 items by default
	PageSize int

	//数据总条数
	//Total number of data
	TotalCount int

	//共多少页
	//How many pages
	PageCount int

	//是否是第一页
	//Is it the first page
	FirstPage bool

	//是否有上一页
	//Whether there is a previous page
	HasPrev bool

	//是否有下一页
	//Is there a next page
	HasNext bool

	//是否是最后一页
	//Is it the last page
	LastPage bool
}

//NewPage 创建Page对象
//NewPage Create Page object
func NewPage() *Page {
	page := Page{}
	page.PageNo = 1
	page.PageSize = 20
	return &page
}

//setTotalCount 设置总条数,计算其他值
//setTotalCount Set the total number of bars, calculate other values
func (page *Page) setTotalCount(total int) {
	page.TotalCount = total
	page.PageCount = (page.TotalCount + page.PageSize - 1) / page.PageSize
	if page.PageNo >= page.PageCount {
		page.LastPage = true
	} else {
		page.HasNext = true
	}
	if page.PageNo > 1 {
		page.HasPrev = true
	} else {
		page.FirstPage = true
	}

}
