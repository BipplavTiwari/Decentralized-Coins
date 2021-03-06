package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type WalletHistoryResponse struct {
	History []interface{}	`json:"history"`
	RollNo	string			`json:"rollNo"`
}

func WalletHistory(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	queriedRollNo := r.URL.Query().Get("rollNo")

	if err := account.ValidateRollNo(queriedRollNo); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid rollNo")
	}

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollNo(requestorRollNo)
	if err != nil {
		return err
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollNo == queriedRollNo) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you are not authorized to read this account history")
	}

	userExists, err := account.UserExists(queriedRollNo)
	if err != nil {
		errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExists {
		return errors.NewHTTPError(err, http.StatusBadRequest, "account does not exist")
	}

	history, err := account.GetWalletHistoryByRollNo(queriedRollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	json.NewEncoder(w).Encode(&WalletHistoryResponse{
		History:	history,
		RollNo:		queriedRollNo,
	})
	return nil
}
