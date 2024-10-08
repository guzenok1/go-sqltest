package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	store "github.com/guzenok1/go-sqltest/sample"
)

func initTestDb(dbUrl string) (db *sql.DB, err error) {
	db, err = sql.Open(driverName, dbUrl)
	if err != nil {
		return
	}

	err = Migrate(db)
	if err != nil {
		db.Close()
		return
	}

	err = loadFixtures(db, "users")
	if err != nil {
		db.Close()
		return
	}

	return
}

func testStoreUsers(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	s := wrap(db)

	user := &store.User{
		ID:       1,
		Login:    "user01",
		Password: "first-P",
	}
	newPassword := "third-P"

	_ = true &&

		t.Run("already exists", func(t *testing.T) {
			var err error
			_, err = s.CreateUser(ctx, user)
			assert.Equal(t,
				`pq: duplicate key value violates unique constraint "users_pkey"`,
				err.Error())
		}) &&

		t.Run("delete", func(t *testing.T) {
			err := s.DeleteUser(ctx, user.ID)
			assert.NoError(t, err)
		}) &&

		t.Run("create", func(t *testing.T) {
			var err error
			user.Password = "second-P"
			user, err = s.CreateUser(ctx, user)
			assert.NoError(t, err)
		}) &&

		t.Run("set password", func(t *testing.T) {
			err := s.SetPassword(ctx, user.ID, newPassword)
			assert.NoError(t, err)
		}) &&

		t.Run("get by login", func(t *testing.T) {
			var err error
			user, err = s.GetUserByLogin(ctx, user.Login)
			_ = assert.NoError(t, err) &&
				assert.NotNil(t, user) &&
				assert.Equal(t, newPassword, user.Password)
		}) &&

		t.Run("get by id", func(t *testing.T) {
			var err error
			user, err = s.GetUserByID(ctx, user.ID)
			_ = assert.NoError(t, err) &&
				assert.NotNil(t, user) &&
				assert.Equal(t, newPassword, user.Password)
		})
}
