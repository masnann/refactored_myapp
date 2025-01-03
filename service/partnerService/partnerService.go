package partnerservice

import (
	"errors"
	"log"
	"myapp/models"
	"myapp/service"
)

type PartnerService struct {
	service service.Service
}

func NewPartnerService(service service.Service) PartnerService {
	return PartnerService{
		service: service,
	}
}

func (s PartnerService) PartnerCreate(req models.PartnerCreateRequest) (int64, error) {

	encryptedSecretKey, err := s.service.Utils.EncryptPlainText(req.SecretKey)
	log.Println("!", encryptedSecretKey)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	newData := models.PartnerModels{
		Name:      req.Name,
		SecretKey: encryptedSecretKey,
	}
	result, err := s.service.PartnerRepo.PartnerCreate(newData)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s PartnerService) PartnerAssignKey(req models.PartnerAssignKeyRequest) (int64, error) {
	partner, err := s.service.PartnerRepo.PartnerFindByPartnerIDandSecretKey(req.PartnerID, req.SecretKey)
	if err != nil {
		log.Println("Partner not found", req.PartnerID, req.SecretKey)
		msg := "partner not found"
		return 0, errors.New(msg)
	}
	encryptedPublicKey, err := s.service.Utils.EncryptPlainText(req.PublicKey)
	if err != nil {
		return 0, err
	}
	newData := models.PartnerModels{
		ID:        partner.ID,
		PublicKey: encryptedPublicKey,
		Domain:    req.Domain,
	}
	result, err := s.service.PartnerRepo.PartnerUpdate(newData)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s PartnerService) PartnerFindByPartnerIDandSecretKey(partnerID int64, secretKey string) (models.PartnerModels, error) {
	result, err := s.service.PartnerRepo.PartnerFindByPartnerIDandSecretKey(partnerID, secretKey)
	if err != nil {
		return result, err
	}
	return result, nil

}

