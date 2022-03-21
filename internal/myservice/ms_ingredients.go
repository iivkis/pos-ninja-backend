package myservice

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iivkis/pos.7-era.backend/internal/repository"
)

type IngredientOutputModel struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Count         float64 `json:"count"`
	MeasureUnit   int     `json:"measure_unit"`
	PurchasePrice float64 `json:"purchase_price"`
	OutletID      uint    `json:"outlet_id"`
}

type IngredientsService struct {
	repo *repository.Repository
}

func newIngredientsService(repo *repository.Repository) *IngredientsService {
	return &IngredientsService{
		repo: repo,
	}
}

type IngredientCreateInput struct {
	Name          string  `json:"name" binding:"required"`
	Count         float64 `json:"count" binding:"min=0"`
	PurchasePrice float64 `json:"purchase_price" binding:"min=0"`
	MeasureUnit   int     `json:"measure_unit" binding:"min=1,max=3"`
}

//@Summary Добавить новый ингредиент в точку
//@param type body IngredientCreateInput false "Принимаемый объект"
//@Accept json
//@Produce json
//@Success 201 {object} DefaultOutputModel "возвращает id созданной записи"
//@Failure 400 {object} serviceError
//@Failure 500 {object} serviceError
//@Router /ingredients [post]
func (s *IngredientsService) Create(c *gin.Context) {
	var input IngredientCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData(err.Error()))
		return
	}

	claims := mustGetEmployeeClaims(c)
	stdQuery := mustGetStdQuery(c)

	ingredient := repository.IngredientModel{
		Name:          input.Name,
		Count:         input.Count,
		PurchasePrice: input.PurchasePrice,
		MeasureUnit:   input.MeasureUnit,
		OutletID:      claims.OutletID,
		OrgID:         claims.OrganizationID,
	}

	if claims.HasRole(repository.R_OWNER, repository.R_DIRECTOR) {
		if stdQuery.OutletID != 0 && s.repo.Outlets.ExistsInOrg(stdQuery.OutletID, claims.OrganizationID) {
			ingredient.OutletID = stdQuery.OutletID
		}
	}

	if err := s.repo.Ingredients.Create(&ingredient); err != nil {
		NewResponse(c, http.StatusInternalServerError, errUnknownDatabase(err.Error()))
		return
	}

	NewResponse(c, http.StatusCreated, DefaultOutputModel{ID: ingredient.ID})
}

type IngredientGetAllOutput []IngredientOutputModel

//@Summary Получить все ингредиенты точки
//@Accept json
//@Produce json
//@Success 200 {object} IngredientGetAllOutput "возвращает все ингредиенты текущей точки"
//@Failure 400 {object} serviceError
//@Failure 500 {object} serviceError
//@Router /ingredients [get]
func (s *IngredientsService) GetAll(c *gin.Context) {
	claims := mustGetEmployeeClaims(c)
	stdQuery := mustGetStdQuery(c)

	where := &repository.IngredientModel{
		OrgID:    claims.OrganizationID,
		OutletID: claims.OutletID,
	}

	if claims.HasRole(repository.R_OWNER, repository.R_DIRECTOR) {
		where.OutletID = stdQuery.OutletID
	}

	ingredients, err := s.repo.Ingredients.Find(where)
	if err != nil {
		NewResponse(c, http.StatusInternalServerError, errUnknownDatabase(err.Error()))
		return
	}

	var output IngredientGetAllOutput = make(IngredientGetAllOutput, len(*ingredients))
	for i, ingredient := range *ingredients {
		output[i] = IngredientOutputModel{
			ID:            ingredient.ID,
			Name:          ingredient.Name,
			Count:         ingredient.Count,
			MeasureUnit:   ingredient.MeasureUnit,
			PurchasePrice: ingredient.PurchasePrice,
			OutletID:      ingredient.OutletID,
		}
	}
	NewResponse(c, http.StatusOK, output)
}

type IngredientUpdateInput struct {
	Name          string  `json:"name"`
	Count         float64 `json:"count" binding:"min=0"`
	PurchasePrice float64 `json:"purchase_price" binding:"min=0"`
	MeasureUnit   int     `json:"measure_unit" binding:"min=0,max=3"`
}

//@Summary Обновить ингредиент
//@param type body IngredientUpdateInput false "Обновляемые поля"
//@Success 200 {object} object "возвращает пустой объект"
//@Accept json
//@Produce json
//@Failure 400 {object} serviceError
//@Failure 500 {object} serviceError
//@Router /ingredients [put]
func (s *IngredientsService) UpdateFields(c *gin.Context) {
	var input IngredientUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData(err.Error()))
		return
	}

	ingrID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData(err.Error()))
		return
	}

	claims := mustGetEmployeeClaims(c)
	stdQuery := mustGetStdQuery(c)

	where := &repository.IngredientModel{
		ID:       uint(ingrID),
		OrgID:    claims.OrganizationID,
		OutletID: claims.OutletID,
	}

	if claims.HasRole(repository.R_OWNER, repository.R_DIRECTOR) {
		where.OutletID = stdQuery.OutletID
	}

	updatedFields := &repository.IngredientModel{
		Name:          input.Name,
		PurchasePrice: input.PurchasePrice,
		Count:         input.Count,
		MeasureUnit:   input.MeasureUnit,
	}

	if err := s.repo.Ingredients.Updates(where, updatedFields); err != nil {
		NewResponse(c, http.StatusInternalServerError, errUnknownDatabase(err.Error()))
		return
	}
	NewResponse(c, http.StatusOK, nil)
}

//@Summary Удаляет ингридиент из точки
//@Accept json
//@Produce json
//@Success 201 {object} object "возвращает пустой объект"
//@Failure 400 {object} serviceError
//@Failure 500 {object} serviceError
//@Router /ingredients/:id [delete]
func (s *IngredientsService) Delete(c *gin.Context) {
	claims := mustGetEmployeeClaims(c)
	stdQuery := mustGetStdQuery(c)

	ingrID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewResponse(c, http.StatusBadRequest, errIncorrectInputData(err.Error()))
		return
	}

	where := &repository.IngredientModel{
		ID:       uint(ingrID),
		OrgID:    claims.OrganizationID,
		OutletID: claims.OutletID,
	}

	if claims.HasRole(repository.R_OWNER, repository.R_DIRECTOR) {
		where.OutletID = stdQuery.OutletID
	}

	if err := s.repo.Ingredients.Delete(where); err != nil {
		if dberr, ok := isDatabaseError(err); ok {
			switch dberr.Number {
			case 1451:
				NewResponse(c, http.StatusBadRequest, errForeignKey("the ingredient has not deleted communications"))
				return
			}
		}
		NewResponse(c, http.StatusBadRequest, errUnknownDatabase(err.Error()))
		return
	}
	NewResponse(c, http.StatusOK, nil)
}
