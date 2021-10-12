package model

/*
|page_size|否|int|分页大小|
|page_no|否|int|分页页数|
*/

// 分页请求
type PageReq struct {
	PageSize int64 `form:"page_size" json:"page_size"`
	PageNo   int64 `form:"page_no" json:"page_no"`
	offset   int64
	limit    int64
}

// 分页响应
type PageResp struct {
	Content interface{} `json:"content"`
	Total   int64       `json:"total"`
}

// 计算偏移量
func (pr *PageReq) OffsetLimit() (int64, int64) {
	if pr.limit > 0 {
		return pr.offset, pr.limit
	}

	if pr.PageSize == 0 {
		pr.PageSize = 10
	}
	if pr.PageNo < 1 {
		pr.PageNo = 1
	}
	pr.limit = pr.PageSize
	if pr.limit < 0 || pr.limit > 5000 {
		pr.limit = 10
	}
	pr.offset = pr.PageSize * (pr.PageNo - 1)

	return pr.offset, pr.limit
}
