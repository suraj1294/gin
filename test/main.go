package main

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/suraj1294/go-gin-planetscale/services"
)

func main() {

	dialect := goqu.Dialect("mysql")

	product := services.Product{Id: 1, Name: "Suraj", Price: 34}

	ds := dialect.Insert("product").Rows(
		&product,
	)
	insertSQL, args, _ := ds.ToSQL()
	fmt.Println(insertSQL, args)

	updateProduct := services.Product{Name: "Suraj"}

	ds1 := dialect.Update("products").Set(&updateProduct).Where(goqu.Ex{"id": goqu.Op{"eq": 123}})

	fmt.Println(updateProduct)

	updateSQL, args, _ := ds1.ToSQL()
	fmt.Println(updateSQL, args)

}
