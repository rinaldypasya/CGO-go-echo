package usecase

import (
	"context"
	"github.com/auth/identityserver"
	"strconv"
	"time"

	"github.com/misc/notif"
	"github.com/transactions/transaction"

	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/transactions/payment"
)

type paymentUsecase struct {
	isUsecase  identityserver.Usecase
	transactionRepo  transaction.Repository
	notificationRepo notif.Repository
	paymentRepo      payment.Repository
	userUsercase     user.Usecase
	bookingRepo      booking_exp.Repository
	userRepo         user.Repository
	contextTimeout   time.Duration
}

// NewPaymentUsecase will create new an paymentUsecase object representation of payment.Usecase interface
func NewPaymentUsecase(isUsecase identityserver.Usecase,t transaction.Repository, n notif.Repository, p payment.Repository, u user.Usecase, b booking_exp.Repository, ur user.Repository, timeout time.Duration) payment.Usecase {
	return &paymentUsecase{
		isUsecase:isUsecase,
		transactionRepo:  t,
		notificationRepo: n,
		paymentRepo:      p,
		userUsercase:     u,
		bookingRepo:      b,
		userRepo:         ur,
		contextTimeout:   timeout,
	}
}

func (p paymentUsecase) Insert(ctx context.Context, payment *models.Transaction, token string, points float64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	var userId string

	if payment.PaymentMethodId == "" {
		return "", models.PaymentMethodIdRequired
	}

	if payment.Currency == "" {
		payment.Currency = "IDR"
	}
	bookingCode := payment.OrderId
	if payment.BookingExpId != nil {
		bookingCode = payment.BookingExpId
	}
	createdBy, err := p.bookingRepo.GetEmailByID(ctx, *bookingCode)
	if err != nil {
		return "", err
	}
	if token != "" {
		currentUser, err := p.userUsercase.ValidateTokenUser(ctx, token)
		if err != nil {
			return "", err
		}
		createdBy = currentUser.UserEmail
		userId = currentUser.Id
	}

	newData := &models.Transaction{
		Id:                  "",
		CreatedBy:           createdBy,
		CreatedDate:         time.Now(),
		ModifiedBy:          nil,
		ModifiedDate:        nil,
		DeletedBy:           nil,
		DeletedDate:         nil,
		IsDeleted:           0,
		IsActive:            1,
		BookingType:         payment.BookingType,
		BookingExpId:        payment.BookingExpId,
		PromoId:             payment.PromoId,
		PaymentMethodId:     payment.PaymentMethodId,
		ExperiencePaymentId: payment.ExperiencePaymentId,
		Status:              payment.Status,
		TotalPrice:          payment.TotalPrice,
		Currency:            payment.Currency,
		OrderId:             payment.OrderId,
		ExChangeRates:payment.ExChangeRates,
		ExChangeCurrency:payment.ExChangeCurrency,
	}

	res, err := p.paymentRepo.Insert(ctx, newData)
	if err != nil {
		return "", models.ErrInternalServerError
	}

	expiredPayment := res.CreatedDate.Add(2 * time.Hour)
	err = p.bookingRepo.UpdateStatus(ctx, *bookingCode, expiredPayment)
	if err != nil {
		return "", err
	}

	if points != 0 {
		err = p.userRepo.UpdatePointByID(ctx, points, userId)
		if err != nil {
			return "", err
		}
	}


	return res.Id, nil
}

func (p paymentUsecase) ConfirmPayment(ctx context.Context, confirmIn *models.ConfirmPaymentIn) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	err := p.paymentRepo.ConfirmPayment(ctx, confirmIn)
	if err != nil {
		return err
	}
	getTransaction, err := p.transactionRepo.GetById(ctx, confirmIn.TransactionID)
	if err != nil {
		return err
	}
	notif := models.Notification{
		Id:           "",
		CreatedBy:    getTransaction.CreatedBy,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		MerchantId:   getTransaction.MerchantId,
		Type:         0,
		Title:        " New Order Receive: Order ID " + getTransaction.OrderIdBook,
		Desc:         "You got a booking for " + getTransaction.ExpTitle + " , booked by " + getTransaction.CreatedBy,
	}
	pushNotifErr := p.notificationRepo.Insert(ctx, notif)
	if pushNotifErr != nil {
		return nil
	}
	if confirmIn.TransactionStatus == 2 && confirmIn.BookingStatus == 1 {
		//confirm
		bookingDetail, err := p.bookingRepo.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "")
		if err != nil {
			return err
		}
		msg := "<h1>" + *bookingDetail.ExpTitle + "</h1>" +
			"<p>Trip Dates :" + bookingDetail.BookingDate.Format("2006-01-01") + "</p>" +
			"<p>Price :" + strconv.FormatFloat(*bookingDetail.TotalPrice, 'f', 6, 64) + "</p>"
		pushEmail := &models.SendingEmail{
			Subject:  "E-Ticket cGO",
			Message:  msg,
			From:     "CGO Indonesia",
			To:       getTransaction.CreatedBy,
			FileName: "Ticket.pdf",
		}
		if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
			return nil
		}
	}else if confirmIn.TransactionStatus == 3 && confirmIn.BookingStatus == 1 {
		//cancelled
		bookingDetail, err := p.bookingRepo.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "")
		if err != nil {
			return err
		}
		msg := "<h1>" + *bookingDetail.ExpTitle + "</h1>" +
			"<p>Trip Dates :" + bookingDetail.BookingDate.Format("2006-01-01") + "</p>" +
			"<p>Price :" + strconv.FormatFloat(*bookingDetail.TotalPrice, 'f', 6, 64) + "</p>"
		pushEmail := &models.SendingEmail{
			Subject:  "Failed Payment",
			Message:  msg,
			From:     "CGO Indonesia",
			To:        getTransaction.CreatedBy,
			FileName: "",
		}
		if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
			return nil
		}
	}

	return nil
}
