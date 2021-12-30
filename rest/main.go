package main 

import (
	"strconv"
    "net/http"
    "fmt"
    "github.com/gin-gonic/gin"
)


var repo Repo

func main() {
    //http.HandleFunc("/", defaultHandler)
    //http.ListenAndServe(":8080", nil)
    
    repo = Repo{Database: "margus", Username: "margus", Password: "margus"}
    
    router := gin.Default()
    router.GET("/albums", getAlbums)
    router.POST("/albums", postAlbums)    
    router.GET("/albums/:id", getAlbumByID)
    
    router.Run("localhost:8080")
    
    repo.Close()
}

func defaultPage(c *gin.Context) {
	html := fmt.Sprintf("<h1>Welcome at %s!</h1>", c.Request.URL.Path[1:])
	
	c.Data(200, "text/html", []byte(html))
}

// getAlbums responds with the list of all albums as JSON.
//     curl http://localhost:8080/albums
// Here's a nice api design discussion https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/
//    paging, ranges, etc.
// here, price=gt:10&price=lt:20  is implemented
func getAlbums(c *gin.Context) {
	//albums := repo.FindAll()
	
fmt.Printf("blaah %s ", c.Query("price"))

	qInt := func (c *gin.Context, name string) int {
		i,_ := strconv.Atoi(c.Query(name))
		return i
	}
	
	qFloat := func (c *gin.Context, name string) float64 {
		f,_ := strconv.ParseFloat(c.Query(name), 64)
		return f
	} 

	q := QueryCriteria{Artist: c.Query("a"), Title: c.Query("t"),
		Price: []float64{qFloat(c, "price:gt"),qFloat(c, "price:lt")},
		Limit: qInt(c, "limit"), Offset: qInt(c, "offset")}
	
	albums := repo.query(q)
	
	// cache these in Redis
	key := q.calcKey()
	fmt.Println("key=" + key)
	
	
	
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
//     curl http://localhost:8080/albums     --include     --header "Content-Type: application/json"     --request "POST"     --data '{"title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
func postAlbums(c *gin.Context) {
    var newAlbum Album

    // Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newAlbum); err != nil {
    	fmt.Println(err)
        return
    }

    // Add the new album to the slice.
    repo.Add(newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
//   curl http://localhost:8080/albums/1
func getAlbumByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))

	album, err := repo.FindById(id)
	if err == nil {
        c.IndentedJSON(http.StatusOK, album)
    } else {
    	fmt.Println(err)
    	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
    }
}


