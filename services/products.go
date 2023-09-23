package services

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/suraj1294/go-gin-planetscale/db"
	"github.com/suraj1294/go-gin-planetscale/logger"
)

//var err error = nil

type Product struct {
	Id    int64  `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	Name  string `json:"name" db:"name"`
	Price int    `json:"price" db:"price, ,omitempty" `
}

type Option func(p Product) Product

func NewProduct(name string, rest ...Option) Product {
	p := Product{}
	p.Name = name

	for _, o := range rest {
		p = o(p)
	}
	return p
}

type ProductService struct {
	sqlx  *sqlx.DB
	mysql *goqu.Database
}

func (p ProductService) GetAll() (*[]Product, *error) {

	err := p.sqlx.DB.Ping()

	if err != nil {
		logger.Error("failed to connect to DB" + err.Error())
		return nil, &err
	}

	ds, _, err := p.mysql.From("products").ToSQL()

	if err != nil {
		logger.Error("failed to generate query all products" + err.Error())
		return nil, &err
	}

	products := []Product{}
	err = p.sqlx.Select(&products, ds)
	if err != nil {
		logger.Error("failed to get products" + err.Error())
		return nil, &err
	}

	return &products, nil
}

func (p ProductService) GetById(id int) (*Product, *error) {

	ds, _, _ := p.mysql.From("products").Select("*").Where(goqu.C("id").Eq(id)).ToSQL()

	product := Product{}
	err := p.sqlx.Get(&product, ds)
	if err != nil {
		logger.Error("failed to get product" + err.Error())
		return nil, &err
	}

	return &product, nil
}

func (p ProductService) Add(newProduct *Product) (*Product, *error) {

	ds := p.mysql.Insert("products").Rows(Product{Name: newProduct.Name, Price: newProduct.Price})

	addQuery, _, _ := ds.ToSQL()

	res, err := p.sqlx.Exec(addQuery)
	if err != nil {
		logger.Error("(CreateProduct) db.Exec")
		return nil, &err
	}
	newProduct.Id, err = res.LastInsertId()
	if err != nil {
		logger.Error("(CreateProduct) res.LastInsertId")
		return nil, &err
	}

	return newProduct, nil
}

func (p ProductService) Update(update *Product, id int) (*Product, *error) {

	ds := p.mysql.Update("products").Set(Product{Name: update.Name}).Where(goqu.Ex{"id": goqu.Op{"eq": id}})

	updateQuery, _, _ := ds.ToSQL()

	fmt.Println(updateQuery)

	_, err := p.sqlx.Exec(updateQuery)

	update.Id = int64(id)

	if err != nil {
		logger.Error("(UpdateProduct) db.Exec")
		return nil, &err
	}
	return update, nil
}

func GetProductService() *ProductService {

	dbCon := db.NewDatabaseConnection()

	return &ProductService{sqlx: dbCon, mysql: goqu.New("mysql", dbCon)}
}
