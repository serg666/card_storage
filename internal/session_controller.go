package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/serg666/card_storage/models"
)

type sessionController struct {
	repo SessionRepository
	keyR KeyRepository
	vm   *validateMethods
}

func (ctrl *sessionController) Create(ctx *gin.Context) {
	v, httpErr := ctrl.vm.SessionsCreate(ctx.Request.Body)
	if httpErr != nil {
		ctx.JSON(httpErr.getHttpStatus(), httpErr.toJson())
		return
	}

	randomKey, err := ctrl.keyR.Query(NewKeySpecificationWithRandom())
	if err != nil {
		logrus.Errorf("Error on Query random crypto key: %v", err)
		ctx.JSON(httpErrInternalServer.getHttpStatus(), httpErrInternalServer.toJson())
		return
	}

	if len(randomKey) == 0 {
		logrus.Errorf("Random crypto key not found")
		ctx.JSON(httpErrInternalServer.getHttpStatus(), httpErrInternalServer.toJson())
		return
	}

	if m, err := ctrl.repo.Create(v.toSessionModel(randomKey[0])); err == nil {
		ctx.JSON(http.StatusCreated, ctrl.sessionAsJson(m))
		return
	}

	ctx.JSON(httpErrInternalServer.getHttpStatus(), httpErrInternalServer.toJson())
}

func (ctrl *sessionController) Get(ctx *gin.Context) {
	if m, err := ctrl.repo.Get(ctx.Param("key")); err == nil {
		ctx.JSON(http.StatusOK, ctrl.sessionAsJson(m))
		return
	}

	ctx.JSON(httpErrModelNotFound.getHttpStatus(), httpErrModelNotFound.toJson())
}

func (ctrl *sessionController) sessionAsJson(c *models.Session) gin.H {
	return gin.H{
		"id":        c.ID,
		"key":       c.Key,
		"body":      c.Body,
		"expire_at": c.ExpireAt,
	}
}

func NewSessionController(repo SessionRepository, keyR KeyRepository) *sessionController {
	return &sessionController{
		repo: repo,
		keyR: keyR,
		vm:   &validateMethods{},
	}
}
