package repository

import (
	"fmt"

	"github.com/iivkis/pos-ninja-backend/internal/config"
	"github.com/iivkis/pos-ninja-backend/pkg/authjwt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	Organizations           *OrganizationsRepo
	Employees               *EmployeesRepo
	Outlets                 *OutletsRepo
	Sessions                *SessionsRepo
	Categories              *CategoriesRepo
	Products                *ProductsRepo
	Ingredients             *IngredientsRepo
	OrdersList              *OrderListRepo
	OrdersInfo              *OrderInfoRepo
	ProductsWithIngredients *ProductsWithIngredientsRepo
}

func NewRepository(authjwt *authjwt.AuthJWT) *Repository {
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True", config.Env.DatabaseLogin, config.Env.DatabasePassword, config.Env.DatabaseIP, config.Env.DatabaseLogin)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&OrganizationModel{},
		&EmployeeModel{},
		&OutletModel{},
		&SessionModel{},
		&ProductModel{},
		&OrderInfoModel{},
		&OrderListModel{},
		&CategoryModel{},
		&IngredientModel{},
		&ProductWithIngredientModel{},
	); err != nil {
		panic(err)
	}

	return &Repository{
		Organizations:           newOrganizationsRepo(db),
		Employees:               newEmployeesRepo(db),
		Outlets:                 newOutletsRepo(db),
		Sessions:                newSessionsRepo(db),
		Categories:              newCategoriesRepo(db),
		Products:                newProductsRepo(db),
		Ingredients:             newIngredientsRepo(db),
		OrdersList:              newOrderListRepo(db),
		OrdersInfo:              newOrderInfoRepo(db),
		ProductsWithIngredients: newProductsWithIngredientsRepo(db),
	}
}
