package transactions

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/princecee/escrow-api/cmd/app/pkg/response"
	"github.com/princecee/escrow-api/config"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/pkg/apis/paystack"
	"github.com/princecee/escrow-api/pkg/json"
	"github.com/princecee/escrow-api/pkg/push"
	"github.com/princecee/escrow-api/pkg/utils"
	"github.com/princecee/escrow-api/pkg/validator"
	"github.com/rs/zerolog"
)

type transactionHandler struct {
	c config.IConfig
}

func (t *transactionHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(createTransactionDto)

	userRepo := t.c.GetUserRepository()
	businessRepo := t.c.GetBusinessRepository()
	transactionRepo := t.c.GetTransactionRepository()
	transactionTimelineRepo := t.c.GetTransactionTimelineRepository()

	tx, _ := t.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	err := json.ReadJSON(r.Body, body)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	var seller *models.Business
	var buyer *models.User

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	if body.CreatedBy == models.TransactionCreatedByBuyer {
		buyer = user
		seller, err = businessRepo.GetById(body.SellerID, tx)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				resp.Message = response.ErrNotFound.Error()
				status = http.StatusNotFound
			default:
				resp.Message = err.Error()
				status = http.StatusBadRequest
			}

			response.SendErrorResponse(w, resp, status)
		}
	} else {
		seller = user.Business
		buyer, err = userRepo.GetById(body.BuyerID, tx)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				resp.Message = response.ErrNotFound.Error()
				status = http.StatusNotFound
			default:
				resp.Message = err.Error()
				status = http.StatusBadRequest
			}

			response.SendErrorResponse(w, resp, status)
		}
	}

	if seller == nil || buyer == nil {
		resp.Message = "seller or buyer not found"
		response.SendErrorResponse(w, resp, http.StatusNotFound)
		return
	}

	var totalAmount, totalCost, charges, receivableAmount int

	productDetails := []models.ProductDetail{}
	for _, v := range body.ProductDetails {
		detail := v

		if body.Type == models.TransactionTypeProduct {
			totalCost += detail.Price * detail.Quantity
		} else {
			totalCost += detail.Price
		}

		product := models.ProductDetail{
			Name:        detail.Name,
			Quantity:    detail.Quantity,
			Description: detail.Description,
			Price:       detail.Price,
		}

		productDetails = append(productDetails, product)
	}

	charges = int(math.Ceil(0.03 * float64(totalCost)))
	totalAmount = charges + totalCost
	receivableAmount = totalCost - int(math.Floor(float64(charges)*float64(body.ChargeConfiguration.SellerCharges/100)))

	transaction := &models.Transaction{
		Status:              models.TransactionStatusAwaiting,
		Type:                body.Type,
		CreatedBy:           body.CreatedBy,
		BuyerID:             body.BuyerID,
		SellerID:            body.SellerID,
		DeliveryDuration:    body.DeliveryDuration,
		Currency:            body.Currency,
		ChargeConfiguration: models.ChargeConfiguration(body.ChargeConfiguration),
		ProductDetails:      productDetails,
		TotalAmount:         totalAmount,
		TotalCost:           totalCost,
		Charges:             charges,
		ReceivableAmount:    receivableAmount,
	}

	err = transactionRepo.Create(transaction, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	timeline := &models.TransactionTimeline{
		Name:          models.TimelineCreated,
		TransactionID: transaction.ID,
	}
	err = transactionTimelineRepo.Create(timeline, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	transaction.Timeline = []*models.TransactionTimeline{timeline}

	utils.Background(func() {
		var email string
		if transaction.CreatedBy == models.TransactionCreatedByBuyer {
			email = seller.Email
		} else {
			email = buyer.Email
		}

		err = t.c.GetPush().SendEmail(&push.Email{
			To:      []string{email},
			Subject: "You have been invited to a new transaction",
			Text:    fmt.Sprintf("You have been invited to a new transaction: %s", transaction.ID),
			Html:    fmt.Sprintf("<p>You have been invited to a new transaction: %s</p>", transaction.ID),
		})

		if err != nil {
			t.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
		}
	})

	if err != nil {
		t.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
	}

	resp.Message = "transaction created successfully"
	resp.Data = map[string]any{
		"transaction": transaction,
	}

	response.SendResponse(w, resp)
}

func (t *transactionHandler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(updateTransactionDto)

	transactionRepo := t.c.GetTransactionRepository()
	transactionTimelineRepo := t.c.GetTransactionTimelineRepository()

	err := json.ReadJSON(r.Body, body)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	transactionId := chi.URLParam(r, "transaction_id")

	tx, _ := t.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	transaction, err := transactionRepo.GetById(transactionId, tx)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = response.ErrNotFound.Error()
			status = http.StatusBadRequest
		default:
			resp.Message = err.Error()
			status = http.StatusInternalServerError
		}

		response.SendErrorResponse(w, resp, status)
		return
	}

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	if transaction.SellerID != *user.BusinessID && transaction.BuyerID != user.ID {
		resp.Message = response.ErrForbidden.Error()
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	if body.DeliveryDuration != nil {
		transaction.DeliveryDuration = *body.DeliveryDuration
	}
	if body.Currency != nil {
		transaction.Currency = *body.Currency
	}
	if body.ChargeConfiguration != nil {
		transaction.ChargeConfiguration = models.ChargeConfiguration(*body.ChargeConfiguration)
	}
	if body.ProductDetails != nil {
		var totalAmount, totalCost, charges, receivableAmount int

		productDetails := []models.ProductDetail{}
		for _, v := range body.ProductDetails {
			detail := v

			if transaction.Type == models.TransactionTypeProduct {
				totalCost += detail.Price * detail.Quantity
			} else {
				totalCost += detail.Price
			}

			product := models.ProductDetail{
				Name:        detail.Name,
				Quantity:    detail.Quantity,
				Description: detail.Description,
				Price:       detail.Price,
			}

			productDetails = append(productDetails, product)
		}

		charges = int(math.Ceil(0.03 * float64(totalCost)))
		totalAmount = charges + totalCost
		receivableAmount = totalCost - int(math.Floor(float64(charges)*float64(body.ChargeConfiguration.SellerCharges/100)))

		transaction.Charges = charges
		transaction.TotalAmount = totalAmount
		transaction.TotalCost = totalCost
		transaction.ReceivableAmount = receivableAmount
		transaction.ProductDetails = productDetails
	}

	// pay from the wallet
	// check appropriate status and timeline

	timeline := &models.TransactionTimeline{
		TransactionID: transactionId,
	}
	if body.Status != nil {
		var text, html string

		transaction.Status = *body.Status

		switch *body.Status {
		case models.TransactionStatusPendingPayment:
			text = "Transaction accepted"
			timeline.Name = models.TimelineApproved
		case models.TransactionStatusCanceled:
			text = "Transaction canceled"
			timeline.Name = models.TImelineCanceled
		case models.TransactionStatusCompleted:
			text = "Transaction completed"
			timeline.Name = models.TimelineCompleted
		}

		html = fmt.Sprintf("<p>%s</p>", text)

		utils.Background(func() {
			err = t.c.GetPush().SendEmail(&push.Email{
				To:      []string{transaction.Seller.Email, transaction.Buyer.Email},
				Subject: "Transaction updated",
				Text:    text,
				Html:    html,
			})

			if err != nil {
				t.c.GetLogger().Log(zerolog.InfoLevel, push.ErrSendingEmailMsg, nil, err)
			}
		})
	}

	err = transactionTimelineRepo.Create(timeline, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = transactionRepo.Update(transaction, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	where, args := utils.GenerateANDWhereFromArgs([]utils.WhereArgs{{
		Name:  "transaction_id",
		Value: transactionId,
	}})
	timelines, err := transactionTimelineRepo.GetMany(args, where, nil)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			transaction.Timeline = []*models.TransactionTimeline{}
		default:
			resp.Message = err.Error()
			status = http.StatusInternalServerError
			response.SendErrorResponse(w, resp, status)
			return
		}
	} else {
		transaction.Timeline = timelines
	}

	err = tx.Commit(context.Background())
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "transaction updated successfully"
	resp.Data = map[string]any{
		"transaction": transaction,
	}
	response.SendResponse(w, resp)
}

func (t *transactionHandler) makePayment(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}
	body := new(makePaymentDto)

	err := json.ReadJSON(r.Body, body)
	defer r.Body.Close()
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	user := r.Context().Value(utils.ContextKey{}).(*models.User)

	tx, _ := t.c.GetDB().Begin(context.Background())
	defer tx.Rollback(context.Background())

	walletRepo := t.c.GetWalletRepository()
	transactionRepo := t.c.GetTransactionRepository()
	transactionTimelineRepo := t.c.GetTransactionTimelineRepository()
	paystackAPI := t.c.GetAPIs().GetPaystack()

	transaction, err := transactionRepo.GetById(body.TransactionID, tx)
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	amount := transaction.TotalCost + int(math.Floor(float64(transaction.Charges)*float64(transaction.ChargeConfiguration.BuyerCharges/100)))
	if body.IsUseWallet {
		wallet, err := walletRepo.GetByIdentifier(user.ID, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		receivable := wallet.Receivable
		diff := receivable - amount
		if diff < 0 {
			resp.Message = "insufficient wallet balance"
			response.SendErrorResponse(w, resp, http.StatusBadRequest)
			return
		}

		wallet.Receivable = diff
		wallet.Balance -= amount

		err = walletRepo.Update(wallet, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		timeline := &models.TransactionTimeline{
			TransactionID: body.TransactionID,
			Name:          models.TimelinePaymentSubmitted,
		}
		err = transactionTimelineRepo.Create(timeline, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		transaction.Status = models.TransactionStatusPendingDelivery
		err = transactionRepo.Update(transaction, tx)
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		err = tx.Commit(context.Background())
		if err != nil {
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}

		resp.Message = "payment made successfully"
		resp.Data = map[string]any{
			"transaction": transaction,
		}
		response.SendResponse(w, resp)
		return
	}

	paystackResponse, err := paystackAPI.InitiateTransaction(paystack.InitiateTransactionDto{
		Email:     user.Email,
		Amount:    strconv.FormatInt(int64(amount), 10),
		MetaData:  map[string]any{"is_for_transaction": true},
		Reference: body.TransactionID,
	})

	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		resp.Message = err.Error()
		response.SendErrorResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Message = "wallet funded successfully"
	resp.Data = map[string]any{
		"transaction":  transaction,
		"payment_data": paystackResponse,
	}
	response.SendResponse(w, resp)
}

func (t *transactionHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	transactionId := chi.URLParam(r, "transaction_id")

	transactionRepo := t.c.GetTransactionRepository()
	transactionTimelineRepo := t.c.GetTransactionTimelineRepository()

	transaction, err := transactionRepo.GetById(transactionId, nil)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = response.ErrNotFound.Error()
			status = http.StatusNotFound
		default:
			resp.Message = err.Error()
			status = http.StatusInternalServerError
		}

		response.SendErrorResponse(w, resp, status)
		return
	}

	user := r.Context().Value(utils.ContextKey{}).(*models.User)
	if transaction.SellerID != *user.BusinessID && transaction.BuyerID != user.ID {
		resp.Message = response.ErrForbidden.Error()
		response.SendErrorResponse(w, resp, http.StatusForbidden)
		return
	}

	where, args := utils.GenerateANDWhereFromArgs([]utils.WhereArgs{{
		Name:  "transaction_id",
		Value: transactionId,
	}})
	timelines, err := transactionTimelineRepo.GetMany(args, where, nil)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			transaction.Timeline = []*models.TransactionTimeline{}
		default:
			resp.Message = err.Error()
			status = http.StatusInternalServerError
			response.SendErrorResponse(w, resp, status)
			return
		}
	} else {
		transaction.Timeline = timelines
	}

	resp.Message = "transaction fetched successfully"
	resp.Data = map[string]any{
		"transaction": transaction,
	}
	response.SendResponse(w, resp)
}

func (t *transactionHandler) getTransactions(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{}

	query := r.URL.Query()

	transactionRepo := t.c.GetTransactionRepository()
	transactionTimelineRepo := t.c.GetTransactionTimelineRepository()

	var page, pageSize int64
	if query.Get("page") != "" {
		page, _ = strconv.ParseInt(query.Get("page"), 10, 64)
	}
	if query.Get("page_size") != "" {
		pageSize, _ = strconv.ParseInt(query.Get("page_size"), 10, 64)
	}

	body := &getTransactionsQueryDto{
		Page:      utils.GetPage(int(page)),
		PageSize:  utils.GetPageSize(int(pageSize)),
		Status:    query.Get("status"),
		Type:      query.Get("type"),
		CreatedBy: query.Get("created_by"),
		BuyerId:   query.Get("buyer_id"),
		SellerId:  query.Get("seller_id"),
	}

	validationErrors := validator.ValidateData(body)
	if validationErrors != nil {
		resp.Message = response.ErrBadRequest.Error()
		resp.Data = validationErrors
		response.SendErrorResponse(w, resp, http.StatusBadRequest)
		return
	}

	pagination := utils.GetPagination(body.Page, body.PageSize)
	where, args := utils.GenerateANDWhereFromArgs([]utils.WhereArgs{
		{
			Name:  "t.status",
			Value: body.Status,
		},
		{
			Name:  "t.type",
			Value: body.Type,
		},
		{
			Name:  "t.created_by",
			Value: body.CreatedBy,
		},
		{
			Name:  "t.seller_id",
			Value: body.SellerId,
		},
		{
			Name:  "t.buyer_id",
			Value: body.BuyerId,
		},
	})

	var total int
	_ = t.c.GetDB().
		QueryRow(
			context.Background(),
			fmt.Sprintf(
				`SELECT COUNT(*) FROM transactions t
				INNER JOIN businesses b ON b.id = t.seller_id
				INNER JOIN users u ON u.id= t.buyer_id
				%s`,
				where,
			),
			args...,
		).
		Scan(&total)

	args = append(args, pagination.Offset, pagination.Limit)
	transactions, err := transactionRepo.GetMany(args, where, nil)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			resp.Message = "transactions fetched successfully"
			resp.Data = map[string]any{
				"transactions": []any{},
			}
			response.SendResponse(w, resp)
			return
		default:
			resp.Message = err.Error()
			response.SendErrorResponse(w, resp, http.StatusInternalServerError)
			return
		}
	}

	for _, v := range transactions {
		t := v

		where, args := utils.GenerateANDWhereFromArgs([]utils.WhereArgs{{
			Name:  "transaction_id",
			Value: t.ID,
		}})
		timelines, err := transactionTimelineRepo.GetMany(args, where, nil)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				t.Timeline = []*models.TransactionTimeline{}
			default:
				resp.Message = err.Error()
				status = http.StatusInternalServerError
				response.SendErrorResponse(w, resp, status)
				return
			}
		} else {
			t.Timeline = timelines
		}
	}

	resp.Message = "transactions fetched successfully"
	resp.Data = map[string]any{
		"transactions": transactions,
	}
	resp.Meta.Page = body.Page
	resp.Meta.PageSize = body.PageSize
	resp.Meta.Total = total
	resp.Meta.TotalPages = utils.GetTotalPages(total, body.PageSize)

	response.SendResponse(w, resp)
}
