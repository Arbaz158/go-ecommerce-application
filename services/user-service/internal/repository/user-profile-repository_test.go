package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	expectedRows := sqlmock.NewRows([]string{"id", "name", "phone", "email"}).
		AddRow(1, "John Doe", "123-456-7890", "john@gmail.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_profiles` WHERE `user_profiles`.`id` = ? ORDER BY `user_profiles`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(expectedRows)

	repo.db = gormDB

	profile, err := repo.GetUserProfileByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile.Name != "John Doe" || profile.Email != "john@gmail.com" {
		t.Fatalf("expected profile with name John Doe, got %v %v", profile.Name, profile.Email)
	}
}

func TestGetMe_EmptyResult(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := &userProfileRepository{}

	expectedRows := sqlmock.NewRows([]string{"id", "name", "phone", "email"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_profiles` WHERE `user_profiles`.`id` = ? ORDER BY `user_profiles`.`id` LIMIT ? ")).
		WithArgs(1, 1).WillReturnRows(expectedRows)

	repo.db = gormDB

	profile, err := repo.GetUserProfileByID(1)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if profile != nil {
		t.Fatalf("expected empty profile but got %v", profile)
	}

}
