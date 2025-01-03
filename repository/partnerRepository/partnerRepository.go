package partnerrepository

import (
	"fmt"
	"log"
	"myapp/constants"
	"myapp/helpers"
	"myapp/models"
	"myapp/repository"
)

type PartnerRepository struct {
	repo repository.Repository
}

func NewPartnerRepository(repo repository.Repository) PartnerRepository {
	return PartnerRepository{
		repo: repo,
	}
}

const (
	partnerTable   = `partner`
	partnerColumns = `id, name, secret_key, public_key, domain, created_at, updated_at, deleted_at`
)

func partnerScanArgs(row *models.PartnerModels) []interface{} {
	return []interface{}{
		&row.ID,
		&row.Name,
		&row.SecretKey,
		&row.PublicKey,
		&row.Domain,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.DeletedAt,
	}
}

func (r PartnerRepository) PartnerCreate(req models.PartnerModels) (int64, error) {
	var ID int64

	query := fmt.Sprintf(
		`INSERT INTO %s
            (
                name, secret_key, public_key, domain, 
				created_at, updated_at, deleted_at
            )
         VALUES 
             (
                ?,?,?,?,
                ?,?,?
            )
         RETURNING id`,
		partnerTable,
	)
	query = helpers.ReplaceSQL(query, "?")
	err := r.repo.DB.QueryRow(query,
		req.Name,
		req.SecretKey,
		req.PublicKey,
		req.Domain,
		req.CreatedAt,
		req.UpdatedAt,
		req.DeletedAt).Scan(&ID)

	if err != nil {
		log.Println("Error querying create partner: ", err)
		return ID, err
	}
	return ID, nil
}

func (r PartnerRepository) PartnerUpdate(req models.PartnerModels) (int64, error) {
	var ID int64
	var args []interface{}

	query := fmt.Sprintf(`
        UPDATE %s 
        SET updated_at =?
    `, partnerTable)
	args = append(args, req.UpdatedAt)

	if req.Name != constants.EMPTY_STRING {
		query += `, name =? `
		args = append(args, req.Name)
	}
	if req.SecretKey != constants.EMPTY_STRING {
		query += `, secret_key =? `
		args = append(args, req.SecretKey)
	}
	if req.PublicKey != constants.EMPTY_STRING {
		query += `, public_key =? `
		args = append(args, req.PublicKey)
	}
	if req.Domain != constants.EMPTY_STRING {
		query += `, domain =? `
		args = append(args, req.Domain)
	}

	query += ` WHERE id =? AND deleted_at = ''
        RETURNING id`
	args = append(args, req.ID)

	query = helpers.ReplaceSQL(query, "?")
	err := r.repo.DB.QueryRow(query, args...).Scan(&ID)
	if err != nil {
		log.Println("Error querying PartnerUpdate: ", err)
		return 0, err
	}
	return ID, nil
}

func (r PartnerRepository) PartnerFindByPartnerIDandSecretKey(partnerID int64, secretKey string) (models.PartnerModels, error) {
	var partner models.PartnerModels
	query := fmt.Sprintf(`
        SELECT %s
        FROM %s
        WHERE id =? AND secret_key =? AND deleted_at = ''`, partnerColumns, partnerTable)

	query = helpers.ReplaceSQL(query, "?")
	row := r.repo.DB.QueryRow(query, partnerID, secretKey)

	err := row.Scan(partnerScanArgs(&partner)...)
	if err != nil {
		log.Println("Error querying PartnerFindByPartnerIDandSecretKey: ", err)
		return partner, err
	}

	return partner, nil
}
