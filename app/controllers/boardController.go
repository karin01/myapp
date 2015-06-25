package controllers

import (
	"log"
	"myapp/app/models"
	"myapp/app/routes"
	"github.com/revel/revel"
	"strconv"
)

const(
	COUNT_PER_PAGE = 10;
)

type Board struct {
	App
}

func (c Board) Index(RequestPage int64) revel.Result {
//////////////////////////////////////////////
//**	처음 화면을 무조건 1페이지로 뜨게 하기 위해
//**	 RequestPage를 설정 후 사용하지 않게 된 소스
//
//	results, err := c.Txn.Select(models.Board{}, `select * from board order by post_id desc`)
//	if err != nil {panic(err)}
//
//	var articles []*models.Board
//	for _ , r := range results{
//		b := r.(*models.Board)
//		articles = append(articles, b)
//	}
//	log.Println(articles)
//////////////////////////////////////////////
	if RequestPage == 0 {RequestPage++}
	articles, pageinfo := c.Page(RequestPage)
	log.Println("게시글들과 페이지 정보\n", articles, pageinfo)
	return c.Render(articles, pageinfo)
} 

/*페이지 계산 mysql에서 공통*/ 
func (c Board) Page(requestPage int64) ([]*models.Board, models.PageInfo){
	//페이지 요청후 페이징 계산 알고리즘..복붙복붙;;
	totalArticleCount, err := c.Txn.SelectInt("select count(*) from board")
    checkErr(err, "select count(*) failed")
    totalPageCount := totalArticleCount / COUNT_PER_PAGE;
    if (totalPageCount % COUNT_PER_PAGE) != 0{
    	totalPageCount++;
    }
    beginPage  := (requestPage - 1) / COUNT_PER_PAGE * COUNT_PER_PAGE + 1
    endPage := beginPage + (COUNT_PER_PAGE-1)
	if endPage > totalPageCount{
		endPage = totalPageCount
	}
	firstRow := (requestPage - 1) * COUNT_PER_PAGE 
	endRow := firstRow + COUNT_PER_PAGE 
	if endRow > totalArticleCount{
		endRow = totalArticleCount
	}
	//여기서부터는 sql
	results, err := c.Txn.Select(models.Board{}, 
		`select * from board order by post_id desc limit ?, ?`, firstRow, endRow-firstRow)
	if err != nil {
		panic(err)
	}
	
	var articles []*models.Board
	for _, r := range results {
		b := r.(*models.Board)
		articles = append(articles, b)
	}
	var pageinfo  models.PageInfo
	pageinfo.BeginPage = beginPage
	pageinfo.EndPage = endPage
	pageinfo.TotalPageCount = totalPageCount
	
	return articles, pageinfo
}

/*글읽기*/
func (c Board) Read(Id int64) revel.Result {
	article := c.loadBoardById(Id)
	return c.Render(article)
}

/*글쓰기 폼*/
func (c Board) FormWrite() revel.Result {
	return c.Render()
}

/*글쓰기*/
func (c Board) Post(board models.Board) revel.Result{
	log.Println(board)
	err := c.Txn.Insert(&board)
	if err != nil{ panic(err)}
	return c.Redirect(routes.Board.Index(1))
}

/*글 삭제*/
func (c Board) Delete (Id int64) revel.Result{
	_, err := c.Txn.Delete( &models.Board{Id:Id} )
	if err != nil{ panic(err) }
	return c.Redirect(routes.Board.Index(1))
}

/*수정 폼*/
func (c Board) FormUpdate(Id int64) revel.Result{
	article := c.loadBoardById(Id)
	return c.Render(article)
}

/*글 수정*/
func (c Board) Update(Id int64, BookName, Writer, Publisher, Body string) revel.Result{
	_ , err := c.Txn.Exec("update board set BookName=?, Writer=?, Publisher=?, Body=? where post_id = ?", BookName, Writer, Publisher, Body, Id)
	if err != nil {panic(err)}
	return c.Redirect(routes.Board.Read(Id))
}


/*글 읽어오기*/
func (c Board) loadBoardById(id int64) *models.Board {
	b, err := c.Txn.Get(models.Board{}, id)
	if err != nil{ panic(err) }
	if b==nil { return nil }
	return b.(*models.Board)
	
} 

func (c Board) Dummy() revel.Result{
	log.Println("더미를 작성합니다.")
	var board models.Board
	for i:=0;i<200;i++{
		board.BookName = "test"+strconv.Itoa(i)
		board.Body ="testtest"
		board.Writer = "루나스타"
		board.Publisher ="Texter"
		err := c.Txn.Insert(&board)
		if err != nil {
			panic(err)
		}
	}
	return c.Redirect(routes.Board.Index(1))
}