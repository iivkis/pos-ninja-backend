package myservice

import (
	"github.com/iivkis/pos.7-era.backend/internal/repository"
	"github.com/iivkis/pos.7-era.backend/pkg/authjwt"
	"github.com/iivkis/pos.7-era.backend/pkg/mailagent"
	"github.com/iivkis/strcode"
)

type DefaultOutputModel struct {
	ID uint `json:"id"`
}

type MyService struct {
	Mware                   *MiddlewareService
	Authorization           *AuthorizationService
	Employees               *EmployeesService
	Outlets                 *OutletsService
	Sessions                *SessionsService
	Categories              *CategoriesService
	Products                *ProductsService
	Ingredients             *IngredientsService
	OrdersList              *OrdersListService
	OrdersInfo              *OrdersInfoService
	ProductsWithIngredients *ProductsWithIngredientsService
	CashChages              *CashChangesService
	InventoryHistory        *InventoryHistoryService
	InventoryList           *InventoryListService
	ReportRevenue           *ReportRevenueService
}

func NewMyService(repo *repository.Repository, strcode *strcode.Strcode, mailagent *mailagent.MailAgent, authjwt *authjwt.AuthJWT) MyService {
	return MyService{
		Mware:                   newMiddlewareService(repo, authjwt),
		Authorization:           newAuthorizationService(repo, strcode, mailagent, authjwt),
		Employees:               newEmployeesService(repo),
		Outlets:                 newOutletsService(repo),
		Sessions:                newSessionsService(repo),
		Categories:              newCategoriesService(repo),
		Products:                newProductsService(repo),
		Ingredients:             newIngredientsService(repo),
		OrdersList:              newOrderListService(repo),
		OrdersInfo:              newOrdersInfoService(repo),
		ProductsWithIngredients: newProductsWithIngredientsService(repo),
		CashChages:              newCashChangesService(repo),
		InventoryHistory:        newInventoryHistoryService(repo),
		InventoryList:           newInventoryListService(repo),
		ReportRevenue:           newReportRevenueService(repo),
	}
}
