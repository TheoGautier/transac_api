package transaction

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"transac_api/lib/structs"
	"transac_api/service/invoice"
	"transac_api/service/user"
)

type Router struct {
	engine         *gin.Engine
	db             *sql.DB
	invoiceService *invoice.Service
	userService    *user.Service
	logger         *log.Logger
}

func MustMakeRouter(engine *gin.Engine, db *sql.DB, invoiceService *invoice.Service, userService *user.Service, logger *log.Logger) *Router {
	return &Router{
		engine:         engine,
		db:             db,
		invoiceService: invoiceService,
		userService:    userService,
		logger:         logger,
	}
}

func (router Router) AddRoutes() {
	group := router.engine.Group("")
	{
		group.POST("/transaction", router.AddTransaction)
	}
}

type AddTransactionRequest struct {
	InvoiceId int     `json:"invoice_id"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

func (router Router) AddTransaction(c *gin.Context) {
	var request AddTransactionRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, structs.GetBadRequestError())
		return
	}
	i, err := router.invoiceService.GetInvoiceById(request.InvoiceId)
	if err != nil {
		if err == invoice.NotFoundError {
			c.JSON(http.StatusNotFound, structs.GetBadRequestError())
			return
		}
		router.logger.Printf("Could not get invoice with id %d, err: %s", request.InvoiceId, err.Error())
		c.JSON(http.StatusInternalServerError, structs.GetBadRequestError())
		return
	}
	if i.Amount != request.Amount {
		c.JSON(http.StatusNotFound, structs.GetBadRequestError())
		return
	}
	if i.Status == invoice.Paid {
		c.JSON(http.StatusUnprocessableEntity, structs.GetBadRequestError())
		return
	}
	tx, err := router.db.Begin()
	if err != nil {
		router.logger.Printf("Could not add transaction for invoice %d, err: %s", i.Id, err.Error())
		c.JSON(http.StatusInternalServerError, structs.GetInternalServerError())
		return
	}
	if err := router.userService.CreditUserBalanceWithTx(i.UserId, i.Amount, tx); err != nil {
		router.logger.Printf("Could not credit user %d balance of %f, err: %s", i.UserId, i.Amount, err.Error())
		c.JSON(http.StatusInternalServerError, structs.GetInternalServerError())
		return
	}
	if err := router.invoiceService.UpdateInvoiceStatusWithTx(i.Id, invoice.Paid, tx); err != nil {
		_ = tx.Rollback()
		router.logger.Printf("Could not mark invoice %d as paid, err: %s", i.Id, err.Error())
		c.JSON(http.StatusInternalServerError, structs.GetInternalServerError())
		return
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		router.logger.Printf("Could not add transaction for invoice %d, err: %s", i.Id, err.Error())
		c.JSON(http.StatusInternalServerError, structs.GetInternalServerError())
	}

	c.JSON(http.StatusNoContent, nil)
}
