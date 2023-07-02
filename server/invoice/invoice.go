package invoice

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"transac_api/lib/structs"
	"transac_api/service/invoice"
)

type Router struct {
	engine         *gin.Engine
	invoiceService *invoice.Service
	logger         *log.Logger
}

func MustMakeRouter(engine *gin.Engine, invoiceService *invoice.Service, logger *log.Logger) *Router {
	return &Router{
		engine:         engine,
		invoiceService: invoiceService,
		logger:         logger,
	}
}

func (router Router) AddRoutes() {
	group := router.engine.Group("")
	{
		group.POST("/invoice", router.AddInvoice)
	}
}

type AddInvoiceRequest struct {
	UserId int     `json:"user_id"`
	Amount float64 `json:"amount"`
	Label  string  `json:"label"`
}

func (router Router) AddInvoice(c *gin.Context) {
	var request AddInvoiceRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, structs.GetBadRequestError())
		return
	}
	if err := router.invoiceService.CreateInvoice(request.UserId, request.Label, request.Amount); err != nil {
		router.logger.Printf("Could not create invoice for user %d, label %s, amount %f, err: %s", request.UserId, request.Label, request.Amount, err.Error())
		c.JSON(http.StatusInternalServerError, structs.GetInternalServerError())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
