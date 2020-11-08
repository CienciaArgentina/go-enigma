package login

import (
	"errors"
	"reflect"
	"testing"

	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_registerRepository_LockAccount(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		userID   int64
		passHash string
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected bool
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  false,
			expected: true,
			mockFunc: func() {
				query := "UPDATE users SET lockout_enabled = 1, lockout_date = ? where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
				query := "UPDATE users SET lockout_enabled = 1, lockout_date = ? where user_id = ?"
				mock.ExpectExec(query).WillReturnError(errors.New("internal_error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			err := u.LockAccount(tt.args.userID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_registerRepository_UnlockAccount(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		userID   int64
		passHash string
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected bool
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  false,
			expected: true,
			mockFunc: func() {
				query := "UPDATE users SET lockout_enabled = 0, lockout_date = null, failed_login_attempts = 0 where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
				query := "UPDATE users SET lockout_enabled = 0, lockout_date = null, failed_login_attempts = 0 where user_id = ?"
				mock.ExpectExec(query).WillReturnError(errors.New("internal_error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			err := u.UnlockAccount(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_registerRepository_ResetLoginFails(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		userID   int64
		passHash string
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected bool
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  false,
			expected: true,
			mockFunc: func() {
				query := "UPDATE users SET failed_login_attempts = 0 where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
				query := "UPDATE users SET failed_login_attempts = 0 where user_id = ?"
				mock.ExpectExec(query).WillReturnError(errors.New("internal_error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			err := u.ResetLoginFails(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_registerRepository_IncrementLoginFailAttempt(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		userID   int64
		passHash string
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected bool
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  false,
			expected: true,
			mockFunc: func() {
				query := "UPDATE users SET failed_login_attempts = failed_login_attempts + 1 where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   123,
				passHash: "123",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
				query := "UPDATE users SET failed_login_attempts = failed_login_attempts + 1 where user_id = ?"
				mock.ExpectExec(query).WillReturnError(errors.New("internal_error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			err := u.IncrementLoginFailAttempt(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_registerRepository_GetUserByUsername(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		username string
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		fields        fields
		args          args
		expectedUser  *domain.User
		expectedEmail *domain.UserEmail
		wantErr       bool
		mockFunc      func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				username: "test",
			},
			wantErr:       false,
			expectedEmail: &domain.UserEmail{UserId: 123},
			expectedUser:  &domain.User{AuthId: 123},
			mockFunc: func() {
				query := "SELECT * FROM users where username = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users_email WHERE user_id = ?"

				table = sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_not_found",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				username: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users where username = ?"

				table := sqlmock.NewRows([]string{"user_id"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				username: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users where username = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("internal_error"))
			},
		},
		{
			name: "email_not_found",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				username: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users where username = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users_email WHERE user_id = ?"

				table = sqlmock.NewRows([]string{"user_id"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "email_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				username: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users where username = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users_email WHERE user_id = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("internal_error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			user, email, err := u.GetUserByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(user, tt.expectedUser) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", user, tt.expectedUser)
				return
			}

			if !reflect.DeepEqual(email, tt.expectedEmail) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", email, tt.expectedEmail)
				return
			}
		})
	}
}
