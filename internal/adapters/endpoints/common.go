package endpoints

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

type errorer interface {
	Error() error
}

func DefaultRequestEncoder(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(errorer); ok && err.Error() != nil {
		encodeError(ctx, err.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

type UserClaimable interface {
	SetUserClaim(claim UserClaim)
}

func DefaultRequestDecoder(decoder func(r *http.Request) (UserClaimable, error)) func(_ context.Context, r *http.Request) (interface{}, error) {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		userClaim, err := GetUserClaimFromRequest(r)
		if err != nil {
			return nil, err
		}

		request, err := decoder(r)
		if err != nil {
			return nil, err
		}
		request.SetUserClaim(*userClaim)
		return request, nil
	}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func GetDefaultHTTPOptions() []kithttp.ServerOption {
	return []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
}
