package common

type successRes struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}

// Hàm tạo đối tượng successRes với dữ liệu, phân trang và bộ lọc
func NewSuccessResponse(data, paging, filter interface{}) *successRes {
	return &successRes{
		Data:   data,
		Paging: paging,
		Filter: filter,
	}
}

// Hàm trả về response thành công đơn giản (không có phân trang và bộ lọc)
func SimpleSuccessResponse(data interface{}) *successRes {
	return NewSuccessResponse(data, nil, nil)
}
