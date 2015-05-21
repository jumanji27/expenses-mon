package main


import (
    "fmt"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    // "labix.org/v2/mgo"
)


// type PagesData struct {
//     Number      int
//     Description string
//     Status      string
// }


func setData() [1][3][12][5]int {
    var data [1][3][12][5]int;

    for i := 0; i < 3; i++ {
        for j := 0; j < 12; j++ {
            for k := 0; k < 5; k++ {
                data[0][i][j][k] = k;
            }
        }
    }

    return data;
}


func main() {
    app := martini.Classic()
    app.Use(render.Renderer())

    app.Get("/api/v1/index", func(render render.Render, params martini.Params) {
        render.JSON(200, setData())
    })

    fmt.Printf("App starting!\n")

    app.Run()
}

