package main

// handle all the imports
import (
	"encoding/csv"
	"fmt"
	"os"
    "time"
) 



// define a Product struct, like a data object
// for test reasons just string datapyes
type Product struct {
    id string   
    organisation_number  string
    image_id string
    erp_id  string
    erp_numberd  string
    brand string
    season string
    is_active string
    product_id string
    gln string
    sku string
    erp_branch string
    size string
    ean string
    stock string
    purchase_price string
    original_price string
    original_sale_price  string
    price string
    sale_price string
}


// main func gets executed on exec the compiled source
func main() {
    // start output + timer
    fmt.Println("Application started")
    start := time.Now()
    // input file
    filename := "lastImported-full.csv"
    // incrementor to check against the map size to see if all datasets 
    // get parsed correctly
    i := 0
    // Open CSV file
    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close() // this is cool shit it gets executed when the method returns

    // Read File into a Variable
    reader := csv.NewReader(f)
    reader.Comma = ';' // change the default deliumiter (,) to ; 
    lines, err := reader.ReadAll()
    // error handling
    if err != nil {
        panic(err)
    }
    // initialize a store map with int pointer and Product struct contents
    var store = make(map[int]Product)
    // Loop through lines & turn into object
    for _, line := range lines {
        data := Product{
            id: line[0],
            organisation_number  : line[1],
            image_id : line[2],
            erp_id : line[3],
            erp_numberd : line[4],
            brand: line[5],
            season : line[6],
            is_active: line[7],
            product_id: line[8],
            gln : line[9],
            sku: line[10],
            erp_branch : line[11],
            size : line[12],
            ean: line[13],
            stock : line[14],
            purchase_price : line[15],
            original_price: line[16],
            original_sale_price : line[17],
            price: line[18],
            sale_price : line[19],
        }
        i++
        // store the filled struct into our map
        store[i] = data

    }
    fmt.Println("Application finished")
    // timer shit & stat output
    elapsed := time.Since(start)
    fmt.Println("Line amount parsed", i , " | Done in ",elapsed, " Seconds | Storage Map Size",len(store) )
}
