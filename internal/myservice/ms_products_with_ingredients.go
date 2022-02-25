package myservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iivkis/pos-ninja-backend/internal/repository"
)

type PWIOutputModel struct {
	ID               uint    `json:"id"`
	CountTakeForSell float64 `json:"count_take_for_sell"`
	ProductID        uint    `json:"product_id"`
	IngredientID     uint    `json:"ingredient_id"`
	OutletID         uint    `json:"outlet_id"`
}

type ProductsWithIngredientsService struct {
	repo *repository.Repository
}

func newProductsWithIngredientsService(repo *repository.Repository) *ProductsWithIngredientsService {
	return &ProductsWithIngredientsService{
		repo: repo,
	}
}

type PWICreateInput struct {
	CountTakeForSell float64 `json:"count_take_for_sell"`
	ProductID        uint    `json:"product_id" binding:"min=1"`
	IngredientID     uint    `json:"ingredient_id" binding:"min=1"`
}

//@Summary Добавить связь продукта и ингридиента в точку
//@param type body PWICreateInput false "Принимаемый объект"
//@Accept json
//@Success 201 {object} object "возвращает пустой объект"
//@Router /pwis [post]
func (s *ProductsWithIngredientsService) Create(c *gin.Context) {
	var input PWICreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData(err.Error()))
		return
	}

	if !s.repo.Products.ExistsInOutlet(input.ProductID, c.MustGet("claims_outlet_id").(uint)) ||
		!s.repo.Ingredients.ExistsInOutlet(input.IngredientID, c.MustGet("claims_outlet_id").(uint)) {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData("not found product or ingredient with this `id` in outlet"))
		return
	}

	m := repository.ProductWithIngredientModel{
		CountTakeForSell: input.CountTakeForSell,
		IngredientID:     input.IngredientID,
		ProductID:        input.ProductID,
		OutletID:         c.MustGet("claims_outlet_id").(uint),
		OrgID:            c.MustGet("claims_org_id").(uint),
	}

	if err := s.repo.ProductsWithIngredients.Create(&m); err != nil {
		NewResponse(c, http.StatusInternalServerError, errUnknownDatabase(err.Error()))
		return
	}

	NewResponse(c, http.StatusCreated, nil)
}

type PWIGetAllForOrgOutput []PWIOutputModel

//@Summary Получить список связей продуктов и ингредиентов
//@Accept json
//@Success 200 {object} PWIGetAllForOrgOutput "Список связей продуктов и ингредиентов"
//@Router /pwis [get]
func (s *ProductsWithIngredientsService) GetAllForOrg(c *gin.Context) {
	pwis, err := s.repo.ProductsWithIngredients.GetAllForOrg(c.MustGet("claims_org_id").(uint))
	if err != nil {
		NewResponse(c, http.StatusInternalServerError, errUnknownDatabase(err.Error()))
		return
	}

	output := make(PWIGetAllForOrgOutput, len(pwis))
	for i, pwi := range pwis {
		output[i] = PWIOutputModel{
			ID:               pwi.ID,
			CountTakeForSell: pwi.CountTakeForSell,
			ProductID:        pwi.ProductID,
			IngredientID:     pwi.IngredientID,
			OutletID:         pwi.OutletID,
		}
	}
	NewResponse(c, http.StatusOK, output)
}

//@Summary Удалить связь
//@Accept json
//@Success 200 {object} object "пустой объект"
//@Router /pwis/:id [delete]
func (s *ProductsWithIngredientsService) Delete(c *gin.Context) {
	err := s.repo.ProductsWithIngredients.Delete(c.Param("id"), c.MustGet("claims_outlet_id"))
	if err != nil {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData(err.Error()))
		return
	}
	NewResponse(c, http.StatusOK, nil)
}

type PWIUpdateFields struct {
	CountTakeForSell float64 `json:"count_take_for_sell"`
}

//@Summary Обновить связь
//@param type body PWIUpdateFields false "Обновляемые поля"
//@Accept json
//@Success 200 {object} object "возвращает пустой объект"
//@Router /pwis/:id [put]
func (s *ProductsWithIngredientsService) UpdateFields(c *gin.Context) {
	var input PWICreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData(err.Error()))
		return
	}

	m := repository.ProductWithIngredientModel{
		CountTakeForSell: input.CountTakeForSell,
	}
	if err := s.repo.ProductsWithIngredients.Updates(&m, c.Param("id"), c.MustGet("claims_outlet_id").(uint)); err != nil {
		NewResponse(c, http.StatusInternalServerError, errUnknownDatabase(err.Error()))
		return
	}
	NewResponse(c, http.StatusOK, nil)
}