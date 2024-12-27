package entity

import "Verve/internal/model/request"

type VerveEntity struct {
	Id string `json:"id"`
}

func GetEntityFromRequest(request request.VerveRequest) VerveEntity {
	return VerveEntity{
		Id: request.Id,
	}
}
