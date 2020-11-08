package register

import (
	"errors"
	"reflect"
	"testing"

	"github.com/CienciaArgentina/go-enigma/internal/domain"
	domain2 "github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestGetUserByIdOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT username FROM users WHERE user_id = ?"
	table := sqlmock.NewRows([]string{"user_id"})
	table.AddRow(123)

	mock.ExpectQuery(query).WillReturnRows(table)

	repo := NewRepository(sqlx.NewDb(db, "sqlmock"))
	got, err := repo.GetUserById(123)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	expected := domain.User{
		AuthId: 123,
	}

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Expected %+v got %+v", expected, got)
		return
	}
}

func TestGetUserByIdInternalError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT username FROM users WHERE user_id = ?"
	mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))

	repo := NewRepository(sqlx.NewDb(db, "sqlmock"))
	_, err = repo.GetUserById(123)
	if err == nil {
		t.Error("Expected error")
		return
	}
}

func TestGetUserByIdErrorNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT username FROM users WHERE user_id = ?"
	table := sqlmock.NewRows([]string{"user_id"})
	mock.ExpectQuery(query).WillReturnRows(table)

	repo := NewRepository(sqlx.NewDb(db, "sqlmock"))
	got, err := repo.GetUserById(123)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}
	if got != nil {
		t.Errorf("Expected nil got %+v", got)
	}
}

func Test_registerRepository_AddUser(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		usr *domain2.User
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Error creating tx %+v", err)
		return
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		want     int64
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				usr: &domain.User{AuthId: 123},
			},
			want:    123,
			wantErr: false,
			mockFunc: func() {
				query := "INSERT INTO users (username, normalized_username, password_hash,  date_created, verification_token, security_token) VALUES (?, ?, ?, now(), ?, ?)"

				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(123, 1))
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				usr: &domain.User{AuthId: 123},
			},
			want:    0,
			wantErr: true,
			mockFunc: func() {
				query := "INSERT INTO users (username, normalized_username, password_hash,  date_created, verification_token, security_token) VALUES (?, ?, ?, now(), ?, ?)"

				mock.ExpectExec(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := &registerRepository{
				db: tt.fields.db,
			}

			got, err := u.AddUser(&sqlx.Tx{Tx: tx}, tt.args.usr)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("registerRepository.AddUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerRepository_AddUserEmail(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		usrEmail *domain2.UserEmail
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Error creating tx %+v", err)
		return
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		want     int64
		wantErr  bool
		mockFunc func()
	}{
		{
			name: "ok",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				usrEmail: &domain.UserEmail{UserId: 123},
			},
			want:    123,
			wantErr: false,
			mockFunc: func() {
				query := "INSERT INTO users_email (user_id, email, normalized_email, date_created) VALUES (?, ?, ?, now())"

				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(123, 1))
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				usrEmail: &domain.UserEmail{UserId: 123},
			},
			want:    0,
			wantErr: true,
			mockFunc: func() {
				query := "INSERT INTO users_email (user_id, email, normalized_email, date_created) VALUES (?, ?, ?, now())"

				mock.ExpectExec(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := &registerRepository{
				db: tt.fields.db,
			}

			got, err := u.AddUserEmail(&sqlx.Tx{Tx: tx}, tt.args.usrEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("registerRepository.AddUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerRepository_DeleteUser(t *testing.T) {
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
			mockFunc: func() {
				query := "DELETE FROM users WHERE user_id = ?"

				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(123, 1))
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID: 123,
			},
			wantErr: true,
			mockFunc: func() {
				query := "DELETE FROM users WHERE user_id = ?"

				mock.ExpectExec(query).WillReturnError(errors.New("Internal error"))
			},
		},
		{
			name: "error_cannot_delete",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				userID: 123,
			},
			wantErr: true,
			mockFunc: func() {
				query := "DELETE FROM users WHERE user_id = ?"

				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := &registerRepository{
				db: tt.fields.db,
			}

			err := u.DeleteUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_registerRepository_CheckUsernameExists(t *testing.T) {
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
				username: "el_pepe",
			},
			wantErr:  false,
			expected: true,
			mockFunc: func() {
				query := "SELECT count(*) FROM users where username = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				username: "el_pepe",
			},
			wantErr:  false,
			expected: false,
			mockFunc: func() {
				query := "SELECT count(*) FROM users where username = ?"

				table := sqlmock.NewRows([]string{"user_id"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				username: "el_pepe",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
				query := "SELECT count(*) FROM users where username = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := &registerRepository{
				db: tt.fields.db,
			}

			got, err := u.CheckUsernameExists(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.expected {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}

func Test_registerRepository_CheckEmailExists(t *testing.T) {
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
				email: "el_pepe",
			},
			wantErr:  false,
			expected: true,
			mockFunc: func() {
				query := "SELECT count(*) FROM users_email WHERE normalized_email = ?"

				table := sqlmock.NewRows([]string{"user_id"})
				table.AddRow(123)

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "doesnt_exist",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe",
			},
			wantErr:  false,
			expected: false,
			mockFunc: func() {
				query := "SELECT count(*) FROM users_email WHERE normalized_email = ?"

				table := sqlmock.NewRows([]string{"user_id"})

				mock.ExpectQuery(query).WillReturnRows(table)
			},
		},
		{
			name: "internal_error",
			fields: fields{
				db: sqlx.NewDb(db, "sqlmock"),
			},
			args: args{
				email: "el_pepe",
			},
			wantErr:  true,
			expected: false,
			mockFunc: func() {
				query := "SELECT count(*) FROM users_email WHERE normalized_email = ?"

				mock.ExpectQuery(query).WillReturnError(errors.New("Internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			u := &registerRepository{
				db: tt.fields.db,
			}

			got, err := u.CheckEmailExists(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.expected {
				t.Errorf("registerRepository.AddUser() got = %v, expected %v", got, tt.expected)
				return
			}
		})
	}
}
