package apihttp

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	controllerV1 "github.com/iivkis/pos.7-era.backend/internal/api/http/controllers/v1"
	"github.com/iivkis/pos.7-era.backend/internal/repository"
	"github.com/iivkis/pos.7-era.backend/pkg/authjwt"
	"github.com/iivkis/pos.7-era.backend/pkg/mailagent"
	"github.com/iivkis/strcode"
)

type apihttp struct {
	engine     *gin.Engine
	repo       *repository.Repository
	strcode    *strcode.Strcode
	postman    *mailagent.MailAgent
	tokenMaker *authjwt.AuthJWT
}

func New(repo *repository.Repository, strcode *strcode.Strcode, postman *mailagent.MailAgent, tokenMaker *authjwt.AuthJWT) *apihttp {
	api := &apihttp{
		repo:       repo,
		strcode:    strcode,
		postman:    postman,
		tokenMaker: tokenMaker,
	}

	api.engine = gin.Default()
	api.init()

	return api
}

func (api *apihttp) init() {
	api.engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
	}))

	controllerV1.AddController(api.engine, api.repo, api.strcode, api.postman, api.tokenMaker)
}

func (api *apihttp) Engine() *gin.Engine {
	return api.engine
}