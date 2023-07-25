package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = sql.ErrNoRows
)

type HTTPError struct {
	Err string `json:"error"`
}

func errBadRequest(err error, c echo.Context) error {
	return errHandler(fmt.Errorf("%w: %w", ErrBadRequest, err), c)
}

func errHandler(err error, c echo.Context) error {
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		code = http.StatusBadRequest
	}

	return c.JSON(code, HTTPError{Err: err.Error()})
}
