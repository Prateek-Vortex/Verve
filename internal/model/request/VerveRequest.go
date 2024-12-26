package request

import (
	"fmt"
	"net/http"
)

type VerveRequest struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func SanitizeUrlParams(r *http.Request) (*VerveRequest, error) {
	query := r.URL.Query()
	id := query.Get("id")
	url := query.Get("url")

	if id == "" {
		return nil, fmt.Errorf("id parameter is required")
	}

	return &VerveRequest{
		Id:  id,
		Url: url,
	}, nil
}
