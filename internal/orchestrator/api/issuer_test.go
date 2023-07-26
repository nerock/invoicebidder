package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/nerock/invoicebidder/internal/invoice"

	"github.com/bojanz/currency"

	"github.com/labstack/echo/v4"
	"github.com/nerock/invoicebidder/internal/issuer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateIssuer(t *testing.T) {
	Convey("CreateIssuer", t, func() {
		Convey("with a valid server and service", func() {
			issSvc := &mockIssuerService{}
			srv := New(0, nil, nil, issSvc, nil)
			rec := httptest.NewRecorder()

			Convey("when request is invalid", func() {
				req := httptest.NewRequest(http.MethodPost, "/", nil)

				Convey("return bad request code and an error", func() {
					err := srv.CreateIssuer(srv.e.NewContext(req, rec))
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
				})
			})

			Convey("when request is valid", func() {
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"fullName":"manu"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				Convey("it tries to create the issuer", func() {
					Convey("when the service fails", func() {
						issSvc.createIssuerFunc = func(_ context.Context, _ string) (issuer.Issuer, error) {
							return issuer.Issuer{}, errors.New("error")
						}

						Convey("return internal error code and an error", func() {
							err := srv.CreateIssuer(srv.e.NewContext(req, rec))
							So(err, ShouldBeNil)
							So(rec.Code, ShouldEqual, http.StatusInternalServerError)
						})
					})

					Convey("when the service succeeds", func() {
						iss := issuer.Issuer{
							ID:       "id",
							FullName: "manu",
						}
						issSvc.createIssuerFunc = func(_ context.Context, _ string) (issuer.Issuer, error) {
							return iss, nil
						}

						Convey("return created code and the issuer", func() {
							err := srv.CreateIssuer(srv.e.NewContext(req, rec))
							So(err, ShouldBeNil)
							So(rec.Code, ShouldEqual, http.StatusCreated)

							js, err := json.Marshal(IssuerResponse{ID: iss.ID, FullName: iss.FullName})
							So(err, ShouldBeNil)
							So(rec.Body.String()[:rec.Body.Len()-1], ShouldEqual, string(js))
						})
					})
				})
			})
		})
	})
}

func TestGetIssuer(t *testing.T) {
	Convey("RetrieveIssuer", t, func() {
		Convey("with a valid server and service", func() {
			issSvc := &mockIssuerService{}
			invSvc := &mockInvoiceService{}
			srv := New(0, invSvc, nil, issSvc, nil)
			rec := httptest.NewRecorder()

			Convey("when request is invalid", func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)

				Convey("return bad request code and an error", func() {
					err := srv.RetrieveIssuer(srv.e.NewContext(req, rec))
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
				})
			})

			Convey("when request is valid", func() {
				issID := "issID"
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"fullName":"manu"}`))
				c := srv.e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues(issID)

				Convey("it tries to create the issuer", func() {
					Convey("when the issuer service fails", func() {
						Convey("because the issuer does not exist", func() {
							issSvc.getIssuerFunc = func(_ context.Context, id string) (issuer.Issuer, error) {
								So(id, ShouldEqual, issID)
								return issuer.Issuer{}, pgx.ErrNoRows
							}

							Convey("return not found error code and an error", func() {
								err := srv.RetrieveIssuer(c)
								So(err, ShouldBeNil)
								So(rec.Code, ShouldEqual, http.StatusNotFound)
							})
						})

						Convey("because an internal error", func() {
							issSvc.getIssuerFunc = func(_ context.Context, id string) (issuer.Issuer, error) {
								So(id, ShouldEqual, issID)
								return issuer.Issuer{}, errors.New("error")
							}

							Convey("return internal error code and an error", func() {
								err := srv.RetrieveIssuer(c)
								So(err, ShouldBeNil)
								So(rec.Code, ShouldEqual, http.StatusInternalServerError)
							})
						})
					})

					Convey("when the issuer service succeeds", func() {
						amount, err := currency.NewAmount("1230.45", "EUR")
						So(err, ShouldBeNil)
						iss := issuer.Issuer{
							ID:       "id",
							FullName: "manu",
							Balance:  amount,
						}
						issSvc.getIssuerFunc = func(_ context.Context, id string) (issuer.Issuer, error) {
							So(id, ShouldEqual, issID)
							return iss, nil
						}

						Convey("when the invoice service fails", func() {
							invSvc.getByIssuerIDFunc = func(_ context.Context, _ string) ([]invoice.Invoice, error) {
								return nil, errors.New("internal error")
							}

							Convey("return internal error code and an error", func() {
								err := srv.RetrieveIssuer(c)
								So(err, ShouldBeNil)
								So(rec.Code, ShouldEqual, http.StatusInternalServerError)
							})
						})

						Convey("when the invoice service succeeds", func() {
							inv := invoice.Invoice{
								ID:       "invID",
								IssuerID: "issID",
								Price:    amount,
								Bids: []invoice.Bid{
									{
										ID:        "bidID",
										InvoiceID: "invID",
										Amount:    amount,
										Active:    true,
									},
								},
								Status: invoice.TRADED,
							}
							invSvc.getByIssuerIDFunc = func(_ context.Context, _ string) ([]invoice.Invoice, error) {
								return []invoice.Invoice{inv}, nil
							}

							Convey("return ok code and the issuer", func() {
								err := srv.RetrieveIssuer(c)
								So(err, ShouldBeNil)
								So(rec.Code, ShouldEqual, http.StatusOK)

								js, err := json.Marshal(IssuerResponse{
									ID:       iss.ID,
									FullName: iss.FullName,
									Balance:  currFmt.Format(amount),
									Invoices: []IssuerInvoiceResponse{
										{
											ID:     inv.ID,
											Price:  currFmt.Format(amount),
											Status: "traded",
											Bids: []IssuerBidResponse{
												{
													ID:     inv.Bids[0].ID,
													Amount: currFmt.Format(amount),
												},
											},
										},
									},
								})
								So(err, ShouldBeNil)
								So(rec.Body.String()[:rec.Body.Len()-1], ShouldEqual, string(js))
							})
						})
					})
				})
			})
		})
	})
}
