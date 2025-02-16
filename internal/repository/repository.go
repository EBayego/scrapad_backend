package repository

import (
	"database/sql"

	"github.com/EBayego/scrapad-backend/internal/domain"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteConnection crea la conexión a la base de datos
func NewSQLiteConnection(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// NewSQLiteRepository crea una instancia de repositorio
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// ----------------
// ORGANIZATIONS
// ----------------

func (r *SQLiteRepository) GetOrganizationByID(orgID string) (*domain.Organization, error) {
	query := `SELECT id, country, created_date FROM organizations WHERE id = ?`
	row := r.db.QueryRow(query, orgID)

	var o domain.Organization
	err := row.Scan(&o.ID, &o.Country, &o.CreatedDate)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

// ----------------
// ADS
// ----------------

func (r *SQLiteRepository) GetAdByID(adID string) (*domain.Ad, error) {
	query := `SELECT id, amount, price, org_id FROM ads WHERE id = ?`
	row := r.db.QueryRow(query, adID)

	var ad domain.Ad
	err := row.Scan(&ad.ID, &ad.Amount, &ad.Price, &ad.OrgID)
	if err != nil {
		return nil, err
	}

	return &ad, nil
}

// Suma la cantidad total de todos los ads de una org
func (r *SQLiteRepository) GetSumAdsPublishedByOrg(orgID string) (int, error) {
	query := `SELECT SUM(amount * price) FROM ads WHERE org_id = ?`
	row := r.db.QueryRow(query, orgID)

	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// ----------------
// FINANCING PROVIDERS
// ----------------

func (r *SQLiteRepository) GetAllFinancingProviders() ([]domain.FinancingProvider, error) {
	query := `SELECT id, slug, payment_method, financing_percentage FROM fnancing_providers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []domain.FinancingProvider
	for rows.Next() {
		var p domain.FinancingProvider
		if err := rows.Scan(&p.ID, &p.Slug, &p.PaymentMethod, &p.FinancingPercentage); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}

	return providers, nil
}

func (r *SQLiteRepository) GetFinancingProviderBySlug(slug string) (*domain.FinancingProvider, error) {
	query := `SELECT id, slug, payment_method, financing_percentage FROM fnancing_providers WHERE slug = ?`
	row := r.db.QueryRow(query, slug)

	var p domain.FinancingProvider
	err := row.Scan(&p.ID, &p.Slug, &p.PaymentMethod, &p.FinancingPercentage)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// ----------------
// OFFERS
// ----------------

func (r *SQLiteRepository) CreateOffer(o domain.Offer) (*domain.Offer, error) {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	query := `INSERT INTO offers(id, payment_method, financing_privder, amount, accepted, price)
              VALUES(?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, o.ID, o.PaymentMethod, o.FinancingProvider, o.Amount, o.Accepted, o.Price)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *SQLiteRepository) UpdateOffer(o domain.Offer) error {
	query := `UPDATE offers SET 
                payment_method = ?, 
                financing_privder = ?, 
                amount = ?, 
                accepted = ?, 
                price = ?
              WHERE id = ?`
	_, err := r.db.Exec(query, o.PaymentMethod, o.FinancingProvider, o.Amount, o.Accepted, o.Price, o.ID)
	return err
}

func (r *SQLiteRepository) GetOfferByID(offerID string) (*domain.Offer, error) {
	query := `SELECT id, payment_method, financing_privder, amount, accepted, price
              FROM offers
              WHERE id = ?`
	row := r.db.QueryRow(query, offerID)

	var o domain.Offer
	err := row.Scan(&o.ID, &o.PaymentMethod, &o.FinancingProvider, &o.Amount, &o.Accepted, &o.Price)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *SQLiteRepository) GetOffersByOrgID(orgID string) ([]domain.Offer, error) {
	// Asumimos que la relación Offer->Ad->Org la deberíamos resolver con JOIN,
	// pero en este ejemplo simplificamos y guardamos orgID en la misma tabla,
	// o lo resolvemos uniendo con Ads. Dependerá del esquema final.

	// Un ejemplo posible:
	query := `
        SELECT offers.id, offers.payment_method, offers.financing_privder, offers.amount, offers.accepted, offers.price
        FROM offers
        JOIN ads ON ads.id = offers.id  -- OJOIN con la tabla 'ads' si coincidiera en algo
        -- Deberíamos tener un campo en 'offers' que relacione con 'ads', p. ej. 'ad_id'
        WHERE ads.org_id = ?
    `
	rows, err := r.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []domain.Offer
	for rows.Next() {
		var o domain.Offer
		err := rows.Scan(&o.ID, &o.PaymentMethod, &o.FinancingProvider, &o.Amount, &o.Accepted, &o.Price)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}
	return offers, nil
}
