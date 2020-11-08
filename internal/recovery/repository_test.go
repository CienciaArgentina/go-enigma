package recovery

import (
	"errors"
	"reflect"
	"testing"

	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_registerRepository_GetEmailByUserId(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		userId int64
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
		expected *domain.UserEmail
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userId: 123,
			},
			wantErr:  false,
			expected: &domain.UserEmail{UserId: 123},
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users_email where user_id = ?"

				table = sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userId: 123,
			},
			wantErr:  true,
			expected: nil,
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

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
				userId: 123,
			},
			wantErr:  true,
			expected: nil,
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
		{
			name: "email_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userId: 123,
			},
			wantErr:  true,
			expected: nil,
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users_email where user_id = ?"

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
				userId: 123,
			},
			wantErr:  true,
			expected: nil,
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users_email where user_id = ?"
				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			_, got, err := u.GetEmailByUserId(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_GetuserIdByEmail(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		email string
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
		expected int64
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe@gmail.com",
			},
			wantErr:  false,
			expected: 123,
			mockFunc: func() {
				query := "SELECT user_id FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe@gmail.com",
			},
			wantErr:  true,
			expected: 0,
			mockFunc: func() {
				query := "SELECT user_id FROM users_email where email = ?"

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
				email: "ete_sech@gmail.com",
			},
			wantErr:  true,
			expected: 0,
			mockFunc: func() {
				query := "SELECT user_id FROM users_email where email = ?"
				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			got, err := u.GetuserIdByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_GetUsernameByEmail(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		email string
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
		expected string
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe@gmail.com",
			},
			wantErr:  false,
			expected: "el_pepe_😎👊",
			mockFunc: func() {
				query := "SELECT user_id FROM users_email WHERE email = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT username FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"username"})
				table.AddRow("el_pepe_😎👊")

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "email_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT user_id FROM users_email WHERE email = ?"

				table := sqlmock.NewRows([]string{"user_id"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "email_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT user_id FROM users_email WHERE email = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
		{
			name: "username_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT user_id FROM users_email WHERE email = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT username FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"username"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "username_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT user_id FROM users_email WHERE email = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT username FROM users where user_id = ?"
				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			got, err := u.GetUsernameByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_GetSecurityToken(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		email string
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
		expected string
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "ete_sech@gmail.com",
			},
			wantErr:  false,
			expected: "token",
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"email", "verified_email"})
				table.AddRow("ete_sech@gmail.com", true)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT security_token FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"security_token"})
				table.AddRow("token")

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "email_not_verified",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "ete_sech@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"email", "verified_email"})
				table.AddRow("ete_sech@gmail.com", false)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "email_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "ete_sech@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"email", "verified_email"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "email_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "ete_sech@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
		{
			name: "token_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "ete_sech@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"email", "verified_email"})
				table.AddRow("ete_sech@gmail.com", true)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT security_token FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"security_token"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "token_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "ete_sech@gmail.com",
			},
			wantErr:  true,
			expected: "",
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"email", "verified_email"})
				table.AddRow("ete_sech@gmail.com", true)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT security_token FROM users where user_id = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			got, err := u.GetSecurityToken(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_GetUserByUserId(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		userID int64
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
		expected *domain.User
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID: 123,
			},
			wantErr: false,
			expected: &domain.User{
				Username: "test",
			},
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

				table := sqlmock.NewRows([]string{"username"})
				table.AddRow("test")

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID: 123,
			},
			wantErr:  true,
			expected: nil,
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

				table := sqlmock.NewRows([]string{"username"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "bad_request",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID: 0,
			},
			wantErr:  true,
			expected: nil,
			mockFunc: func() {
			},
		},
		{
			name: "user_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID: 123,
			},
			wantErr:  true,
			expected: nil,
			mockFunc: func() {
				query := "SELECT * FROM users where user_id = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			got, err := u.GetUserByUserId(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_UpdateSecurityToken(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		userID           int64
		newSecurityToken string
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
			name: "bad_request",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:           0,
				newSecurityToken: "",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
			},
		},
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:           123,
				newSecurityToken: "123",
			},
			wantErr:  false,
			expected: true,
			mockFunc: func() {
				query := "UPDATE users SET security_token = ? where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "no_affected_rows",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:           123,
				newSecurityToken: "123",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
				query := "UPDATE users SET security_token = ? where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			got, err := u.UpdateSecurityToken(tt.args.userID, tt.args.newSecurityToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_UpdatePasswordHash(t *testing.T) {
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
			name: "bad_request",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID:   0,
				passHash: "",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
			},
		},
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
				query := "UPDATE users SET password_hash = ?  where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "no_affected_rows",
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
				query := "UPDATE users SET password_hash = ?  where user_id = ?"
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			got, err := u.UpdatePasswordHash(tt.args.userID, tt.args.passHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_ConfirmUserEmail(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		email string
		token string
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
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: false,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})
				table.AddRow(false, 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"verification_token", "user_id"})
				table.AddRow("test", 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "UPDATE users_email SET verified_email = 1, verification_date = now() WHERE user_id = ?"

				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "email_not_found",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "email_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"
				mock.ExpectQuery(query).WillReturnError(errors.New("internal_error"))
			},
		},
		{
			name: "user_email_already_verified",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})
				table.AddRow(true, 123)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_not_found",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})
				table.AddRow(false, 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"verification_token", "user_id"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_not_verified",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})
				table.AddRow(false, 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"verification_token", "user_id"})
				table.AddRow("", 123)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "user_internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})
				table.AddRow(false, 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users where user_id = ?"
				mock.ExpectQuery(query).WillReturnError(errors.New("internal_error"))
			},
		},
		{
			name: "no_affected_rows",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})
				table.AddRow(false, 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"verification_token", "user_id"})
				table.AddRow("test", 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "UPDATE users_email SET verified_email = 1, verification_date = now() WHERE user_id = ?"

				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "error_updating_verified_email",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "test@gmail.com",
				token: "test",
			},
			wantErr: true,
			mockFunc: func() {
				query := "SELECT * FROM users_email where email = ?"

				table := sqlmock.NewRows([]string{"verified_email", "user_id"})
				table.AddRow(false, 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "SELECT * FROM users where user_id = ?"

				table = sqlmock.NewRows([]string{"verification_token", "user_id"})
				table.AddRow("test", 123)

				mock.ExpectQuery(query).WillReturnRows(table)

				query = "UPDATE users_email SET verified_email = 1, verification_date = now() WHERE user_id = ?"

				mock.ExpectExec(query).WillReturnError(errors.New("internal_error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := NewRepository(tt.fields.db)

			err := u.ConfirmUserEmail(tt.args.email, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
