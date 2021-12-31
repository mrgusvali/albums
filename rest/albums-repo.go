package main 

import (
	"fmt"
	"context"
	"strings"
    "github.com/go-pg/pg/v10"
)

// loosely around https://go.dev/doc/tutorial/web-service-gin
// Album represents data about a record Album.
type Album struct {
    Id     int  `json:"id" pg:",pk"` 
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

type QueryCriteria struct {
	Title string
	Artist string
	// expect ranges for:
	Year []int
	Price []float64
	
	Limit int
	Offset int
}

func (c QueryCriteria) calcKey() string {
	return fmt.Sprintf("t=%s,a=%s,p=%f:%f,o=%d", c.Title, c.Artist, c.Price[0], c.Price[1], c.Offset)
}

type AlbumsRepository interface {
	FindById(id string) (Album, error)
	FindAll() []Album
	Add(a Album) Album
}

type Repo struct {
	Database string	
	Username string 
	Password string
	db *pg.DB
}


type dbLogger struct { }

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
    return c,nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	qs,_ := q.FormattedQuery()
    fmt.Println(string(qs))
    return nil
}

func (r Repo) connect() *pg.DB {
	if r.db == nil {
		r.db = pg.Connect(&pg.Options{
	        User: r.Username,
	        Password: r.Password,
	        Database: r.Database,
	    })
		r.db.AddQueryHook(dbLogger{})
	}
	return r.db
}

func (r Repo) Close() {
	r.db.Close()
}

func (r Repo) FindById(id int) (Album, error) {
	db := r.connect()
    
    a := Album{Id: id}
    err := db.Model(&a).WherePK().Select()
    
    if err != nil {
    	panic(err)
    }
    
    return a, err
}

func (r Repo) FindAll() []Album {
	var Albums []Album
	err := r.connect().Model(&Albums).Select()
    if err != nil {
        panic(err)
    }
    return Albums
}

func (r Repo) Add(a Album) Album {
	_, err := r.connect().Model(&a).Insert()
	if err != nil {
        panic(err)
    }
	return a
}

func toLowcaseMatcher(s string) string {
	return fmt.Sprintf("%%%s%%", strings.ToLower(s))
}

func (r Repo) query(q QueryCriteria) []Album {
	var Albums []Album
	m := r.connect().Model(&Albums)
	if q.Title != "" {
		m = m.Where("lower(title) like ?", toLowcaseMatcher(q.Title))
	}
	if q.Artist != "" {
		m = m.Where("lower(artist) like ?", toLowcaseMatcher(q.Artist))
	}
	if len(q.Price) > 0 && q.Price[1] > 0{
		m = m.Where("price between ? and ?", q.Price[0], q.Price[1])
	}
	if q.Offset > 0 {
		m = m.Offset(q.Offset)
	}
	err := m.Select()
    if err != nil {
        panic(err)
    }
    return Albums	
}