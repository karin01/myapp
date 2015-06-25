package models

import (
	"fmt"	
	"time"
	"log"
	"github.com/coopernurse/gorp"
)

const (
	SQL_DATE_FORMAT ="2006-01-02"
)

type Board struct {
	Id int64 `db:"post_id"`		// `를 사용해야만 db에 입력될 때 post_id가 된다.
	BookName string				// 대문자를 사용하면 public, 소문자 사용하면 private
	Writer string
	Publisher	string
	Body	string
	
	DayWriterStr	string
	//trasient
	DayWrite time.Time
}

func (c Board) String() string{
	return fmt.Sprintf("게시물 번호 : %d, 책이름 : %s, 저자 : %s, 출판사 : %s,\n", c.Id, c.BookName, c.Writer, c.Publisher)
}

func (b *Board) PreInsert(_ gorp.SqlExecutor) error {
	b.DayWrite = time.Now()
	b.DayWriterStr = b.DayWrite.Format(SQL_DATE_FORMAT)
	log.Println(b)
	return nil
}

func (b *Board) PostGet(_ gorp.SqlExecutor) error {
	var (
		err error
	)
	if b.DayWrite, err = time.Parse(SQL_DATE_FORMAT, b.DayWriterStr); err!=nil{
		log.Println("작성일 에러... 검증에러...")
		return fmt.Errorf("Error parsing check in date '%s':", b.DayWriterStr, err)
	}
	return nil
}