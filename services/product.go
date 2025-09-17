package services

import (
	"github.com/salahfarzin/roja-shop/repositories"
	"github.com/salahfarzin/roja-shop/types"
)

type Product interface {
	GetAll(perPage, offset int) ([]types.Product, error)
	GetOne(id string) (*types.Product, error)
	Create(product types.Product, file *types.File) (string, error)
	Update(id string, input types.Product) error
}

type product struct {
	repo repositories.Product
}

func NewProduct(repo repositories.Product) Product {
	return &product{repo: repo}
}

func (p *product) Create(product types.Product, file *types.File) (string, error) {
	return p.repo.CreateWithFile(product, file)
}

func (p *product) Update(id string, input types.Product) error {
	return p.repo.Update(id, input)
}

func (p *product) GetOne(id string) (*types.Product, error) {
	return p.repo.FetchOne(id)
}

func (p *product) GetAll(perPage, offset int) ([]types.Product, error) {
	return p.repo.FetchAll(perPage, offset)
}
