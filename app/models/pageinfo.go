package models

import(
	"fmt"
)

type PageInfo struct{
	BeginPage 	int64
	EndPage		int64
	TotalPageCount	int64
}

func (p PageInfo) String() string{
	return fmt.Sprintf("PageInfo (Begin : %d , EndPage : %d , TotalPageCount : %d",
	p.BeginPage, p.EndPage, p.TotalPageCount)
}

func (p PageInfo) PrevBeginPage() int64{
	return p.BeginPage - int64(1)
}
func (p PageInfo) NextEndPage() int64{
	return p.EndPage +  int64(1)
}

func (p PageInfo) FirstPage() int64{
	return p.BeginPage
}

func (p PageInfo) LastPage() int64{
	return p.EndPage
}


func (p PageInfo) Pagenation() []int64{
	var pageRow = make([]int64, p.EndPage- p.BeginPage +1)
	for i:=0;p.BeginPage + int64(i) <= p.EndPage; i++{
		pageRow[i] = p.BeginPage + int64(i)
	}
	return pageRow
}