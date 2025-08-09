package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Repository is a thin generic repository that works with any models.
// It provides CRUD, pagination, queries by field and transaction composition.
type Repository struct {
	db *gorm.DB
}

// New creates a Repository bound to the provided *gorm.DB (or a tx).
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts the given entity into DB.
func (r *Repository) Create(ctx context.Context, entity any) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update saves the provided entity.
func (r *Repository) Update(ctx context.Context, entity any) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete deletes the provided entity (or by primary key if entity is a model with ID set).
func (r *Repository) Delete(ctx context.Context, entity any) error {
	return r.db.WithContext(ctx).Delete(entity).Error
}

// DeleteByID deletes a model by primary key value.
func (r *Repository) DeleteByID(ctx context.Context, model any, id any) error {
	return r.db.WithContext(ctx).Delete(model, id).Error
}

// GetByID finds a single record by primary key. Returns (nil, nil) when not found.
func (r *Repository) GetByID(ctx context.Context, model any, id any, out any) error {
	if err := r.db.WithContext(ctx).Model(model).First(out, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}

// GetAll finds all records for model and scans into out.
func (r *Repository) GetAll(ctx context.Context, model any, out any) error {
	return r.db.WithContext(ctx).Model(model).Find(out).Error
}

// GetPaginated finds records with offset/limit and scans into out.
func (r *Repository) GetPaginated(ctx context.Context, model any, out any, page, pageSize int) (int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(model)
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	if page <= 0 || pageSize <= 0 {
		if err := q.Find(out).Error; err != nil {
			return 0, err
		}
		return total, nil
	}
	offset := (page - 1) * pageSize
	if err := q.Offset(offset).Limit(pageSize).Find(out).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetByField finds records where field = value and scans into out.
// model: a pointer to the model type or model instance for GORM's Model()
// field: column name (e.g. "key" or "email")
// value: value to match
func (r *Repository) GetByField(ctx context.Context, model any, field string, value any, out any) error {
	cond := fmt.Sprintf("%s = ?", field)
	return r.db.WithContext(ctx).Model(model).Where(cond, value).Find(out).Error
}

// Transaction runs the provided function inside a transaction. Commit is automatic when fn returns nil,
// rollback if fn returns an error. The txRepo provided uses the transactional *gorm.DB.
func (r *Repository) Transaction(ctx context.Context, fn func(txRepo *Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := New(tx)
		return fn(txRepo)
	})
}

// ManualTx returns a started transaction (*gorm.DB) so the caller can control Commit/Rollback.
// Caller must call tx.Commit() or tx.Rollback().
func (r *Repository) ManualTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	return tx, tx.Error
}
