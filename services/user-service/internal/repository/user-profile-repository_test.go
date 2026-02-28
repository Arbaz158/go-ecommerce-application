package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-ecommerce-application/services/user-service/internal/domain/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to create gorm db: %v", err)
	}

	return db, mock, gormDB
}

func TestGetMe_Success(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := &userProfileRepository{}

	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "email", "first_name", "last_name", "contact_number"}).
		AddRow(1, "ewuhiwj23200", "john@gmail.com", "John", "Doe", "123-456-7890")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_profiles` WHERE user_id = ? ORDER BY `user_profiles`.`id` LIMIT ?")).
		WithArgs("ewuhiwj23200", 1).
		WillReturnRows(expectedRows)

	repo.db = gormDB

	profile, err := repo.GetUserProfileByUserID("ewuhiwj23200")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile.FirstName != "John" || profile.LastName != "Doe" {
		t.Fatalf("expected profile with name John Doe, got %v %v", profile.FirstName, profile.LastName)
	}
}

func TestGetMe_EmptyResult(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := &userProfileRepository{}

	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "email", "first_name", "last_name", "contact_number"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_profiles` WHERE user_id = ? ORDER BY `user_profiles`.`id` LIMIT ?")).
		WithArgs("ewuhiwj23200", 1).WillReturnRows(expectedRows)

	repo.db = gormDB

	profile, err := repo.GetUserProfileByUserID("ewuhiwj23200")
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if profile != nil {
		t.Fatalf("expected empty profile but got %v", profile)
	}

}

func TestSaveUserAdress_Success(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := &userProfileRepository{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `addresses` (`user_id`,`street`,`city`,`state`,`postal_code`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(1, "123 Main St", "Cityville", "Stateville", "12345", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo.db = gormDB

	address := &models.Address{
		UserID:     1,
		Street:     "123 Main St",
		City:       "Cityville",
		State:      "Stateville",
		PostalCode: "12345",
	}

	err := repo.SaveAddress(address)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetUserAddresses_Success(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := &userProfileRepository{}

	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "street", "city", "state", "postal_code", "created_at", "updated_at"}).
		AddRow(1, 1, "123 Main St", "Cityville", "Stateville", "12345", 1769077351, 1769077351).
		AddRow(2, 1, "456 Side St", "Townsville", "Regionville", "67890", 1769077351, 1769077351)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `addresses` WHERE user_id = ?")).
		WithArgs(1).
		WillReturnRows(expectedRows)

	repo.db = gormDB

	addresses, err := repo.GetUserAddresses(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(addresses) != 2 {
		t.Fatalf("expected 2 addresses, got %d", len(addresses))
	}
}

func TestGetUserAddresses_EmptyResult(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := &userProfileRepository{}

	expectedRows := sqlmock.NewRows([]string{"id", "user_id", "street", "city", "state", "zip_code", "country"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `addresses` WHERE user_id = ?")).
		WithArgs(1).
		WillReturnRows(expectedRows)

	repo.db = gormDB

	addresses, err := repo.GetUserAddresses(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(addresses) != 0 {
		t.Fatalf("expected 0 addresses, got %d", len(addresses))
	}
}
