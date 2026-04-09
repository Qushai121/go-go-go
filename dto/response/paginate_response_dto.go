package response

type PaginateResponseDto[T any] struct {
	Data        T     `json:"data"`
	TotalRecord int64 `json:"totalRecord"`
	TotalPage   int
}
