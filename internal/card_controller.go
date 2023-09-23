package internal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/serg666/card_storage/models"
)

type cardController struct {
	repo CardRepository
	keyR KeyRepository
	vm   *validateMethods
}

func (ctrl *cardController) Index(ctx *gin.Context) {
	v, httpErr := ctrl.vm.CardsList(ctx)
	if httpErr != nil {
		ctx.JSON(httpErr.getHttpStatus(), httpErr.toJson())
		return
	}

	if m, err := ctrl.repo.ListByPan(v.filtersModel(), v.paginator()); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data": ctrl.cardsAsJson(m),
		})
		return
	}

	ctx.JSON(httpErrModelNotFound.getHttpStatus(), httpErrModelNotFound.toJson())
}

func (ctrl *cardController) GetByToken(ctx *gin.Context) {
	v, httpErr := ctrl.vm.GetCardByToken(ctx)
	if httpErr != nil {
		ctx.JSON(httpErr.getHttpStatus(), httpErr.toJson())
		return
	}

	if m, err := ctrl.repo.GetByToken(v.Token); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data": ctrl.cardAsJson(m),
		})
		return
	}

	ctx.JSON(httpErrModelNotFound.getHttpStatus(), httpErrModelNotFound.toJson())
}

func (ctrl *cardController) Create(ctx *gin.Context) {
	v, httpErr := ctrl.vm.CardsCreate(ctx.Request.Body)
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

	if m, err := ctrl.repo.Create(v.toCardModel(randomKey[0])); err == nil {
		ctx.JSON(http.StatusCreated, ctrl.cardAsJson(m))
		return
	} else {
		logrus.Errorf("Error on Create card: %v", err)
	}

	ctx.JSON(httpErrInternalServer.getHttpStatus(), httpErrInternalServer.toJson())
}

func (ctrl *cardController) Delete(ctx *gin.Context) {
	cardId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.Errorf("Error on convert to int: %v", err)
		ctx.JSON(httpErrInternalServer.getHttpStatus(), httpErrInternalServer.toJson())
		return
	}

	if err := ctrl.repo.Delete(&models.Card{ID: uint(cardId)}); err != nil {
		logrus.Errorf("Error on delete card_id: %d, err=%v", cardId, err)
		ctx.JSON(httpErrInternalServer.getHttpStatus(), httpErrInternalServer.toJson())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (ctrl *cardController) cardsAsJson(cards []*models.Card) []gin.H {
	l := make([]gin.H, 0)
	for _, c := range cards {
		l = append(l, ctrl.cardAsJson(c))
	}
	return l
}

func (ctrl *cardController) cardAsJson(c *models.Card) gin.H {
	return gin.H{
		"id":       c.ID,
		"pan":      c.PAN,
		"holder":   c.Holder,
		"exp_date": c.FormatedExpDate(),
		"token":    c.Token,
	}
}

func NewCardController(repo CardRepository, keyR KeyRepository) *cardController {
	return &cardController{
		repo: repo,
		keyR: keyR,
		vm:   &validateMethods{},
	}
}
