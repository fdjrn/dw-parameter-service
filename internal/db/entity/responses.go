package entity

// normal responses result
// -------------------------------------------

type Responses struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// paginated responses result
// -------------------------------------------

type ResponsePaginatedData struct {
	Result      interface{} `json:"results,omitempty"`
	Total       int64       `json:"total,omitempty"`
	PerPage     int64       `json:"perPage,omitempty"`
	CurrentPage int64       `json:"currentPage,omitempty"`
	LastPage    int64       `json:"lastPage,omitempty"`
}

type ResponsePaginated struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    ResponsePaginatedData `json:"data,omitempty"`
}
