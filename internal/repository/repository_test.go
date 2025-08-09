package repository

import (
	"context"
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCarRepository_DeleteByID(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	car := Car{Brand: "Mazda", Color: "Gray", Year: 2019, Model: "3"}
	if err := repo.Create(ctx, &car); err != nil {
		t.Fatalf("failed to create car: %v", err)
	}

	if err := repo.DeleteByID(ctx, &Car{}, car.ID); err != nil {
		t.Fatalf("DeleteByID failed: %v", err)
	}

	var got Car
	err := db.WithContext(ctx).First(&got, car.ID).Error
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound after DeleteByID, got %v", err)
	}
}

func TestCarRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	car := Car{Brand: "Nissan", Color: "Green", Year: 2017, Model: "Sentra"}
	if err := repo.Create(ctx, &car); err != nil {
		t.Fatalf("failed to create car: %v", err)
	}

	var got Car
	err := repo.GetByID(ctx, &Car{}, car.ID, &got)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Brand != car.Brand || got.Color != car.Color || got.Year != car.Year || got.Model != car.Model {
		t.Errorf("got %+v, want %+v", got, car)
	}
}

func TestCarRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	for _, c := range carGetAllTestData {
		if err := repo.Create(ctx, &Car{Brand: c.Brand, Color: c.Color, Year: c.Year, Model: c.Model}); err != nil {
			t.Fatalf("failed to create car: %v", err)
		}
	}

	var cars []Car
	err := repo.GetAll(ctx, &Car{}, &cars)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(cars) != len(carGetAllTestData) {
		t.Errorf("expected %d cars, got %d", len(carGetAllTestData), len(cars))
	}
}

func TestCarRepository_GetPaginated(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	for _, c := range carPaginatedTestData {
		if err := repo.Create(ctx, &Car{Brand: c.Brand, Color: c.Color, Year: c.Year, Model: c.Model}); err != nil {
			t.Fatalf("failed to create car: %v", err)
		}
	}

	var cars []Car
	total, err := repo.GetPaginated(ctx, &Car{}, &cars, 2, 2)
	if err != nil {
		t.Fatalf("GetPaginated failed: %v", err)
	}
	if total != int64(len(carPaginatedTestData)) {
		t.Errorf("expected total %d, got %d", len(carPaginatedTestData), total)
	}
	if len(cars) != 2 {
		t.Errorf("expected 2 cars in page, got %d", len(cars))
	}
}

func TestCarRepository_GetByField(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	for _, c := range carByFieldTestData {
		if err := repo.Create(ctx, &Car{Brand: c.Brand, Color: c.Color, Year: c.Year, Model: c.Model}); err != nil {
			t.Fatalf("failed to create car: %v", err)
		}
	}

	var cars []Car
	err := repo.GetByField(ctx, &Car{}, "brand", "Peugeot", &cars)
	if err != nil {
		t.Fatalf("GetByField failed: %v", err)
	}
	if len(cars) != 1 || cars[0].Brand != "Peugeot" {
		t.Errorf("expected 1 Peugeot, got %+v", cars)
	}
}

func TestCarRepository_Transaction(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	err := repo.Transaction(ctx, func(txRepo *Repository) error {
		car := Car{Brand: "Fiat", Color: "Yellow", Year: 2015, Model: "Uno"}
		return txRepo.Create(ctx, &car)
	})
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	var cars []Car
	if err := repo.GetByField(ctx, &Car{}, "brand", "Fiat", &cars); err != nil {
		t.Fatalf("failed to get cars by field: %v", err)
	}
	if len(cars) != 1 {
		t.Errorf("expected 1 Fiat after transaction, got %d", len(cars))
	}
}

func TestCarRepository_ManualTx(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	tx, err := repo.ManualTx(ctx)
	if err != nil {
		t.Fatalf("ManualTx failed: %v", err)
	}
	car := Car{Brand: "Chevrolet", Color: "Silver", Year: 2016, Model: "Onix"}
	if err := tx.Create(&car).Error; err != nil {
		t.Fatalf("failed to create car in tx: %v", err)
	}
	tx.Commit()

	var got Car
	db.First(&got, car.ID)
	if got.Brand != "Chevrolet" {
		t.Errorf("expected Chevrolet, got %+v", got)
	}
}

func TestCarRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	for _, tt := range carTestCases {
		t.Run(tt.name, func(t *testing.T) {
			car := tt.input
			if err := db.WithContext(ctx).Create(&car).Error; err != nil {
				t.Fatalf("failed to create car: %v", err)
			}
			if car.ID == 0 {
				t.Fatal("expected car ID to be set after create")
			}
		})
	}
}

func TestCarRepository_Read(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	for _, tt := range carTestCases {
		t.Run(tt.name, func(t *testing.T) {
			car := tt.input
			if err := db.WithContext(ctx).Create(&car).Error; err != nil {
				t.Fatalf("failed to create car: %v", err)
			}
			var got Car
			if err := db.WithContext(ctx).First(&got, car.ID).Error; err != nil {
				t.Fatalf("failed to fetch car: %v", err)
			}
			if got.Brand != tt.expected.Brand || got.Color != tt.expected.Color || got.Year != tt.expected.Year || got.Model != tt.expected.Model {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestCarRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	for _, tt := range carTestCases {
		t.Run(tt.name, func(t *testing.T) {
			car := tt.input
			if err := db.WithContext(ctx).Create(&car).Error; err != nil {
				t.Fatalf("failed to create car: %v", err)
			}
			if err := db.WithContext(ctx).Model(&car).Updates(Car{Color: "White", Year: 2022}).Error; err != nil {
				t.Fatalf("failed to update car: %v", err)
			}
			var updated Car
			if err := db.WithContext(ctx).First(&updated, car.ID).Error; err != nil {
				t.Fatalf("failed to fetch updated car: %v", err)
			}
			if updated.Color != "White" || updated.Year != 2022 {
				t.Errorf("update failed, got %+v", updated)
			}
		})
	}
}

func TestCarRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	for _, tt := range carTestCases {
		t.Run(tt.name, func(t *testing.T) {
			car := tt.input
			if err := db.WithContext(ctx).Create(&car).Error; err != nil {
				t.Fatalf("failed to create car: %v", err)
			}
			if err := db.WithContext(ctx).Delete(&Car{}, car.ID).Error; err != nil {
				t.Fatalf("failed to delete car: %v", err)
			}
			var deleted Car
			err := db.WithContext(ctx).First(&deleted, car.ID).Error
			if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
				t.Errorf("expected ErrRecordNotFound after delete, got %v", err)
			}
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&Car{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}
