package services

import (
	"github.com/salahfarzin/roja-shop/repositories"
	"github.com/salahfarzin/roja-shop/types"
)

type Product interface {
	Store(product types.Product, file *types.File) (string, error)
	GetAll() ([]types.Product, error)
}

type product struct {
	repo repositories.Product
}

func NewProduct(repo repositories.Product) Product {
	return &product{repo: repo}
}

func (p *product) Store(product types.Product, file *types.File) (string, error) {
	return p.repo.CreateWithFile(product, file)
}

func (p *product) GetAll() ([]types.Product, error) {
	return p.repo.FetchAll()
}
