package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/salahfarzin/roja-shop/repositories"
	"github.com/salahfarzin/roja-shop/types"
)

var ErrInsufficientInventory = errors.New("insufficient inventory")

type Sale interface {
	Sell(productID string, quantity int) (*types.Product, error)
}

type sale struct {
	repo        repositories.Sale
	productRepo repositories.Product
}

func NewSale(repo repositories.Sale, productRepo repositories.Product) Sale {
	return &sale{repo: repo, productRepo: productRepo}
}

func (s *sale) Sell(productID string, quantity int) (*types.Product, error) {
	prod, err := s.productRepo.FetchOne(productID)
	if err != nil {
		return nil, err
	}
	if prod.Inventory < quantity {
		return nil, ErrInsufficientInventory
	}
	prod.Inventory -= quantity
	prod.SoldCount += quantity
	if err := s.productRepo.Update(productID, *prod); err != nil {
		return nil, err
	}
	sale := &types.Sale{
		ID:        uuid.New().String(),
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: time.Now(),
	}

	err = s.repo.Create(sale)
	if err != nil {
		return nil, err
	}

	return prod, nil
}
