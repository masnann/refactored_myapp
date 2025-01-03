package models

type PartnerModels struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	SecretKey string `json:"secretKey"`
	PublicKey string `json:"-"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"-"`
}

type PartnerCreateRequest struct {
	Name      string `json:"name" validate:"required"`
	SecretKey string `json:"secretKey" validate:"required"`
}

type PartnerAssignKeyRequest struct {
	PartnerID int64  `json:"partnerID" validate:"required"`
	SecretKey string `json:"secretKey" validate:"required"`
	PublicKey string `json:"publicKey"`
	Domain    string `json:"domain"`
}
