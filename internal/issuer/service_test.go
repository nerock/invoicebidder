package issuer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/bojanz/currency"
	. "github.com/smartystreets/goconvey/convey"
)

type mockStorage struct {
	createIssuerFunc   func(context.Context, Issuer) error
	retrieveIssuerFunc func(context.Context, string) (Issuer, error)
	updateBalanceFunc  func(context.Context, string, currency.Amount) error
}

func (m *mockStorage) CreateIssuer(ctx context.Context, issuer Issuer) error {
	return m.createIssuerFunc(ctx, issuer)
}

func (m *mockStorage) RetrieveIssuer(ctx context.Context, s string) (Issuer, error) {
	return m.retrieveIssuerFunc(ctx, s)
}

func (m *mockStorage) UpdateBalance(ctx context.Context, s string, amount currency.Amount) error {
	return m.updateBalanceFunc(ctx, s, amount)
}

func TestService_CreateIssuer(t *testing.T) {
	Convey("CreateIssuer", t, func() {
		st := &mockStorage{}
		svc := NewService(st)

		Convey("when storage fails", func() {
			st.createIssuerFunc = func(ctx context.Context, issuer Issuer) error {
				return errors.New("some error")
			}

			Convey("return an empty issuer and an error", func() {
				iss, err := svc.CreateIssuer(context.Background(), "name")
				So(err, ShouldNotBeNil)
				So(iss, ShouldEqual, Issuer{})
			})
		})

		Convey("when storage succeeds", func() {
			st.createIssuerFunc = func(ctx context.Context, issuer Issuer) error {
				So(issuer.FullName, ShouldEqual, "name")
				return nil
			}

			Convey("return a valid issuer and no error", func() {
				iss, err := svc.CreateIssuer(context.Background(), "name")
				So(err, ShouldBeNil)
				So(iss.FullName, ShouldEqual, "name")
				So(iss.Balance.IsZero(), ShouldBeTrue)
				_, err = uuid.Parse(iss.ID)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestService_ApproveTrade(t *testing.T) {
	Convey("ApproveTrade", t, func() {
		st := &mockStorage{}
		svc := NewService(st)

		Convey("when storage fails to retrieve the issuer", func() {
			st.retrieveIssuerFunc = func(ctx context.Context, id string) (Issuer, error) {
				So(id, ShouldEqual, "id")
				return Issuer{}, errors.New("some error")
			}

			Convey("return an error", func() {
				err := svc.ApproveTrade(context.Background(), "id", currency.Amount{})
				So(err, ShouldNotBeNil)
			})
		})

		Convey("when storage successfully retrieves the issuer", func() {
			st.retrieveIssuerFunc = func(ctx context.Context, id string) (Issuer, error) {
				So(id, ShouldEqual, "id")
				b, _ := currency.NewAmount("1000", "EUR")
				return Issuer{
					ID:       id,
					FullName: "name",
					Balance:  b,
				}, nil
			}
			amount, _ := currency.NewAmount("1000", "EUR")

			Convey("when storage fails to update issuer balance", func() {
				st.updateBalanceFunc = func(ctx context.Context, id string, a currency.Amount) error {
					So(id, ShouldEqual, "id")
					amount, _ := currency.NewAmount("2000", "EUR")
					So(a, ShouldEqual, amount)
					return errors.New("some error")
				}

				Convey("return no error", func() {
					err := svc.ApproveTrade(context.Background(), "id", amount)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("when storage successfully updates issuer balance", func() {
				st.updateBalanceFunc = func(ctx context.Context, id string, a currency.Amount) error {
					So(id, ShouldEqual, "id")
					amount, _ := currency.NewAmount("2000", "EUR")
					So(a, ShouldEqual, amount)
					return nil
				}

				Convey("return no error", func() {
					err := svc.ApproveTrade(context.Background(), "id", amount)
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
