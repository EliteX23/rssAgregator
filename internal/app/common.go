package app

import (
	"errors"
	"github.com/gocraft/dbr"
)

//QueryFilters - фильтр для запросов. Можно расширить на еще несколько параметров
type QueryFilters struct {
	Limit  int    `form:"limit" query:"limit"`
	Offset int    `form:"offset" query:"offset"`
	Title  string `form:"title" query:"title"`
}

func (q QueryFilters) IsValid() bool {
	return q.Limit > 0 && q.Offset >= 0
}
func (q QueryFilters) MakeStmt(selectStmt *dbr.SelectStmt) (*dbr.SelectStmt,error) {
	if q.IsValid() {
		return selectStmt.
			Limit(uint64(q.Limit)).
			Offset(uint64(q.Offset)),nil
	}else{
		return nil, errors.New("bad pagination settings")
	}

	return selectStmt,nil
}

type PagedList struct {
	Total  int64       `json:"total"`
	Result interface{} `json:"result,omitempty"`
}
