package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gooferOrm/goofer/test"
	"github.com/gooferOrm/gooferest/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestTransactionReturns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "transaction",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			createUserA := client.User.CreateOne(
				User.Email.Set("a"),
				User.ID.Set("a"),
			).Tx()

			createUserB := client.User.CreateOne(
				User.Email.Set("b"),
				User.ID.Set("b"),
			).Tx()

			if err := client.Prisma.Transaction(createUserA, createUserB).Exec(ctx); err != nil {
				t.Fatal(err)
			}

			expectedA := &UserModel{
				InnerUser: InnerUser{
					ID:    "a",
					Email: "a",
				},
			}

			expectedB := &UserModel{
				InnerUser: InnerUser{
					ID:    "b",
					Email: "b",
				},
			}

			massert.Equal(t, expectedA, createUserA.Result())
			massert.Equal(t, expectedB, createUserB.Result())

			// --

			actual, err := client.User.FindMany().Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:    "a",
					Email: "a",
				},
			}, {
				InnerUser: InnerUser{
					ID:    "b",
					Email: "b",
				},
			}}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "transaction find many",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			createUser := client.User.CreateOne(
				User.Email.Set("a"),
				User.ID.Set("a"),
			).Tx()

			update := client.User.FindMany().Update(
				User.Email.Set("new"),
			).Tx()

			if err := client.Prisma.Transaction(createUser, update).Exec(ctx); err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, &BatchResult{Count: 1}, update.Result())

			// --

			actual, err := client.User.FindMany().Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:    "a",
					Email: "new",
				},
			}}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "transaction result caching",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			createUserA := client.User.CreateOne(
				User.Email.Set("a"),
				User.ID.Set("a"),
			).Tx()

			createUserB := client.User.CreateOne(
				User.Email.Set("b"),
				User.ID.Set("b"),
			).Tx()

			if err := client.Prisma.Transaction(createUserA, createUserB).Exec(ctx); err != nil {
				t.Fatal(err)
			}

			expectedA := &UserModel{
				InnerUser: InnerUser{
					ID:    "a",
					Email: "a",
				},
			}

			expectedB := &UserModel{
				InnerUser: InnerUser{
					ID:    "b",
					Email: "b",
				},
			}

			massert.Equal(t, expectedA, createUserA.Result())
			massert.Equal(t, expectedB, createUserB.Result())
			massert.Equal(t, expectedA, createUserA.Result())
			massert.Equal(t, expectedB, createUserB.Result())
			massert.Equal(t, expectedA, createUserA.Result())
			massert.Equal(t, expectedB, createUserB.Result())

			// --

			actual, err := client.User.FindMany().Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:    "a",
					Email: "a",
				},
			}, {
				InnerUser: InnerUser{
					ID:    "b",
					Email: "b",
				},
			}}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "rollback tx",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "exists",
					email: "email",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			// this will fail...
			aOp := client.User.FindUnique(
				User.ID.Equals("does-not-exist"),
			).Update(
				User.Email.Set("foo"),
			).Tx()

			// ...so this should be roll-backed
			bOp := client.User.FindUnique(
				User.ID.Equals("exists"),
			).Update(
				User.Email.Set("new"),
			).Tx()

			err := client.Prisma.Transaction(aOp, bOp).Exec(ctx)
			assert.Errorf(t, err, "should error")

			assert.Panics(t, func() {
				aOp.Result()
			})

			assert.Panics(t, func() {
				bOp.Result()
			})

			// make sure the existing record wasn't touched

			actual, err := client.User.FindMany().Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:    "exists",
					Email: "email",
				},
			}}

			massert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, test.Databases, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
