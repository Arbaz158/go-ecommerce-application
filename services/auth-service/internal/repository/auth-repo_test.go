package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-ecommerce-application/services/auth-service/internal/database"
	"github.com/go-ecommerce-application/services/auth-service/internal/models"
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

func TestCreateUser_Success(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	testUser := models.AuthUser{
		Id:       "user-123",
		Email:    "test@example.com",
		Password: "hashedPassword123",
		Role:     "user",
		Status:   "active",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `auth_users` (`id`,`email`,`password`,`role`,`status`) VALUES (?,?,?,?,?)")).
		WithArgs(testUser.Id, testUser.Email, testUser.Password, testUser.Role, testUser.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(testUser)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateUser_WithDifferentRoles(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	roles := []string{"admin", "user", "moderator"}

	for _, role := range roles {
		testUser := models.AuthUser{
			Id:       "user-" + role,
			Email:    "test-" + role + "@example.com",
			Password: "hashedPassword" + role,
			Role:     role,
			Status:   "active",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `auth_users` (`id`,`email`,`password`,`role`,`status`) VALUES (?,?,?,?,?)")).
			WithArgs(testUser.Id, testUser.Email, testUser.Password, role, "active").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateUser(testUser)

		if err != nil {
			t.Fatalf("unexpected error for role %s: %v", role, err)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateUser_WithDifferentStatuses(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	statuses := []string{"active", "inactive", "pending", "suspended"}

	for _, status := range statuses {
		testUser := models.AuthUser{
			Id:       "user-" + status,
			Email:    "test-" + status + "@example.com",
			Password: "hashedPassword" + status,
			Role:     "user",
			Status:   status,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `auth_users` (`id`,`email`,`password`,`role`,`status`) VALUES (?,?,?,?,?)")).
			WithArgs(testUser.Id, testUser.Email, testUser.Password, "user", status).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateUser(testUser)

		if err != nil {
			t.Fatalf("unexpected error for status %s: %v", status, err)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateUser_WithDuplicateEmail(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	firstUser := models.AuthUser{
		Id:       "user-1",
		Email:    "duplicate@example.com",
		Password: "hashedPassword123",
		Role:     "user",
		Status:   "active",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `auth_users` (`id`,`email`,`password`,`role`,`status`) VALUES (?,?,?,?,?)")).
		WithArgs(firstUser.Id, firstUser.Email, firstUser.Password, "user", "active").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err1 := repo.CreateUser(firstUser)

	if err1 != nil {
		t.Fatalf("first user creation failed: %v", err1)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreateUser_WithEmptyId(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	testUser := models.AuthUser{
		Id:       "", // Empty ID
		Email:    "test@example.com",
		Password: "hashedPassword123",
		Role:     "user",
		Status:   "active",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `auth_users` (`id`,`email`,`password`,`role`,`status`) VALUES (?,?,?,?,?)")).
		WithArgs("", testUser.Email, testUser.Password, "user", "active").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	err := repo.CreateUser(testUser)

	if err == nil {
		t.Fatalf("expected error for empty ID, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	testUser := models.AuthUser{
		Id:       "user-123",
		Email:    "success@example.com",
		Password: "hashedPassword123",
		Role:     "user",
		Status:   "active",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "status"}).
		AddRow(testUser.Id, testUser.Email, testUser.Password, testUser.Role, testUser.Status)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `auth_users` WHERE email = ? ORDER BY `auth_users`.`id` LIMIT ?")).
		WithArgs(testUser.Email, 1).
		WillReturnRows(rows)

	retrievedUser, err := repo.GetUserByEmail(testUser.Email)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedUser == nil {
		t.Fatalf("expected user, got nil")
	}

	if retrievedUser.Email != testUser.Email {
		t.Fatalf("expected email %s, got %s", testUser.Email, retrievedUser.Email)
	}

	if retrievedUser.Id != testUser.Id {
		t.Fatalf("expected id %s, got %s", testUser.Id, retrievedUser.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "status"})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `auth_users` WHERE email = ? ORDER BY `auth_users`.`id` LIMIT ?")).
		WithArgs("nonexistent@example.com", 1).
		WillReturnRows(rows)

	retrievedUser, err := repo.GetUserByEmail("nonexistent@example.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedUser != nil {
		t.Fatalf("expected nil user, got %v", retrievedUser)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByEmail_WithMultipleUsers(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	targetUser := models.AuthUser{
		Id:       "user-2",
		Email:    "bob@example.com",
		Password: "pass2",
		Role:     "admin",
		Status:   "active",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "status"}).
		AddRow(targetUser.Id, targetUser.Email, targetUser.Password, targetUser.Role, targetUser.Status)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `auth_users` WHERE email = ? ORDER BY `auth_users`.`id` LIMIT ?")).
		WithArgs(targetUser.Email, 1).
		WillReturnRows(rows)

	retrievedUser, err := repo.GetUserByEmail(targetUser.Email)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedUser == nil {
		t.Fatalf("expected user, got nil")
	}

	if retrievedUser.Email != "bob@example.com" {
		t.Fatalf("expected email bob@example.com, got %s", retrievedUser.Email)
	}

	if retrievedUser.Id != "user-2" {
		t.Fatalf("expected id user-2, got %s", retrievedUser.Id)
	}

	if retrievedUser.Role != "admin" {
		t.Fatalf("expected role admin, got %s", retrievedUser.Role)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByEmail_VerifyFieldValues(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	testUser := models.AuthUser{
		Id:       "verify-user-123",
		Email:    "verify@example.com",
		Password: "secureHashedPassword",
		Role:     "moderator",
		Status:   "pending",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "status"}).
		AddRow(testUser.Id, testUser.Email, testUser.Password, testUser.Role, testUser.Status)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `auth_users` WHERE email = ? ORDER BY `auth_users`.`id` LIMIT ?")).
		WithArgs(testUser.Email, 1).
		WillReturnRows(rows)

	retrievedUser, err := repo.GetUserByEmail(testUser.Email)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedUser == nil {
		t.Fatalf("expected user, got nil")
	}

	if retrievedUser.Id != testUser.Id {
		t.Fatalf("expected id %s, got %s", testUser.Id, retrievedUser.Id)
	}

	if retrievedUser.Email != testUser.Email {
		t.Fatalf("expected email %s, got %s", testUser.Email, retrievedUser.Email)
	}

	if retrievedUser.Password != testUser.Password {
		t.Fatalf("expected password %s, got %s", testUser.Password, retrievedUser.Password)
	}

	if retrievedUser.Role != testUser.Role {
		t.Fatalf("expected role %s, got %s", testUser.Role, retrievedUser.Role)
	}

	if retrievedUser.Status != testUser.Status {
		t.Fatalf("expected status %s, got %s", testUser.Status, retrievedUser.Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByEmail_WithSpecialCharacters(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	testUser := models.AuthUser{
		Id:       "user-special",
		Email:    "user+tag@example.co.uk",
		Password: "hashedPassword123",
		Role:     "user",
		Status:   "active",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "status"}).
		AddRow(testUser.Id, testUser.Email, testUser.Password, testUser.Role, testUser.Status)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `auth_users` WHERE email = ? ORDER BY `auth_users`.`id` LIMIT ?")).
		WithArgs(testUser.Email, 1).
		WillReturnRows(rows)

	retrievedUser, err := repo.GetUserByEmail(testUser.Email)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedUser == nil {
		t.Fatalf("expected user, got nil")
	}

	if retrievedUser.Email != "user+tag@example.co.uk" {
		t.Fatalf("expected email user+tag@example.co.uk, got %s", retrievedUser.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByEmail_WithEmptyEmail(t *testing.T) {

	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	database.DB = gormDB
	repo := NewAuthRepository()

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "status"})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `auth_users` WHERE email = ? ORDER BY `auth_users`.`id` LIMIT ?")).
		WithArgs("", 1).
		WillReturnRows(rows)

	retrievedUser, err := repo.GetUserByEmail("")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedUser != nil {
		t.Fatalf("expected nil user, got %v", retrievedUser)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
