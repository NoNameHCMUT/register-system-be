package business

import (
	"register-system-be/model"
	"register-system-be/repo"
)

type AffiliationBusiness interface {
	ListAll() ([]model.AffiliationResponse, error)
}

type affiliationBusiness struct {
	affiliationRepo repo.AffiliationRepo
}

func NewAffiliationBusiness(ar repo.AffiliationRepo) AffiliationBusiness {
	return &affiliationBusiness{affiliationRepo: ar}
}

func (b *affiliationBusiness) ListAll() ([]model.AffiliationResponse, error) {
	list, err := b.affiliationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	res := make([]model.AffiliationResponse, 0, len(list))
	for _, aff := range list {
		if aff.ID == 0 {
			continue
		}
		res = append(res, model.AffiliationResponse{
			ID:          aff.ID,
			StdName:     aff.StdName,
			Description: aff.Description,
		})
	}

	return res, nil
}
