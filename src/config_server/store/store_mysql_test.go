package store_test

import (
	. "config_server/store"

	"config_server/store/fakes"
	"database/sql"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StoreMysql", func() {

	var (
		fakeDbProvider *fakes.FakeDbProvider
		fakeDb         *fakes.FakeIDb
		fakeRow        *fakes.FakeIRow
		fakeResult     *fakes.FakeResult

		store Store
	)

	BeforeEach(func() {
		fakeDbProvider = &fakes.FakeDbProvider{}
		fakeDb = &fakes.FakeIDb{}
		fakeRow = &fakes.FakeIRow{}
		fakeResult = &fakes.FakeResult{}

		store = NewMysqlStore(fakeDbProvider)
	})

	Describe("GetByName", func() {

		It("closes db connection on exit", func() {
			fakeDb.QueryRowReturns(&fakes.FakeIRow{})
			fakeDbProvider.DbReturns(fakeDb, nil)

			store.GetByName("Luke")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("queries the database for the latest entry for a given name", func() {
			fakeDb.QueryRowReturns(&fakes.FakeIRow{})
			fakeDbProvider.DbReturns(fakeDb, nil)

			_, err := store.GetByName("Luke")
			Expect(err).To(BeNil())
			query, _ := fakeDb.QueryRowArgsForCall(0)

			Expect(query).To(Equal("SELECT id, name, value FROM configurations WHERE name = ? ORDER BY id DESC LIMIT 1"))
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("returns value from db query", func() {
			fakeRow.ScanStub = func(dest ...interface{}) error {
				idPtr, ok := dest[0].(*string)
				Expect(ok).To(BeTrue())

				namePtr, ok := dest[1].(*string)
				Expect(ok).To(BeTrue())

				valuePtr, ok := dest[2].(*string)
				Expect(ok).To(BeTrue())

				*idPtr = "some_id"
				*valuePtr = "Skywalker"
				*namePtr = "Luke"

				return nil
			}

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := store.GetByName("Luke")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(Configuration{
				Id:    "some_id",
				Value: "Skywalker",
				Name:  "Luke",
			}))
		})

		It("returns empty configuration when no result is found", func() {
			fakeRow.ScanReturns(sql.ErrNoRows)

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := store.GetByName("luke")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(Configuration{}))
		})

		It("returns an error when db provider fails to return db", func() {
			dbError := errors.New("connection failure")
			fakeDbProvider.DbReturns(nil, dbError)

			_, err := store.GetByName("luke")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(dbError))
		})

		It("returns an error when db query fails", func() {
			scanError := errors.New("query failure")
			fakeRow.ScanReturns(scanError)

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			_, err := store.GetByName("luke")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(scanError))
		})
	})

	Describe("GetById", func() {

		It("closes db connection on exit", func() {
			fakeDb.QueryRowReturns(&fakes.FakeIRow{})
			fakeDbProvider.DbReturns(fakeDb, nil)

			store.GetById("1")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("queries the database for the latest entry for a given id", func() {
			fakeDb.QueryRowReturns(&fakes.FakeIRow{})
			fakeDbProvider.DbReturns(fakeDb, nil)

			_, err := store.GetById("1")
			Expect(err).To(BeNil())
			query, _ := fakeDb.QueryRowArgsForCall(0)

			Expect(query).To(Equal("SELECT id, name, value FROM configurations WHERE id = ?"))
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("returns value from db query", func() {
			fakeRow.ScanStub = func(dest ...interface{}) error {
				idPtr, ok := dest[0].(*string)
				Expect(ok).To(BeTrue())

				namePtr, ok := dest[1].(*string)
				Expect(ok).To(BeTrue())

				valuePtr, ok := dest[2].(*string)
				Expect(ok).To(BeTrue())

				*idPtr = "54"
				*valuePtr = "Skywalker"
				*namePtr = "Luke"

				return nil
			}

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := store.GetById("54")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(Configuration{
				Id:    "54",
				Value: "Skywalker",
				Name:  "Luke",
			}))
		})

		It("returns empty configuration when no result is found", func() {
			fakeRow.ScanReturns(sql.ErrNoRows)

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := store.GetById("54")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(Configuration{}))
		})

		It("returns an error when db provider fails to return db", func() {
			dbError := errors.New("connection failure")
			fakeDbProvider.DbReturns(nil, dbError)

			_, err := store.GetById("2")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(dbError))
		})

		It("returns an error when db query fails", func() {
			scanError := errors.New("query failure")
			fakeRow.ScanReturns(scanError)

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			_, err := store.GetById("7")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(scanError))
		})
	})

	Describe("Put", func() {

		It("closes db connection on exit", func() {
			fakeDbProvider.DbReturns(fakeDb, nil)

			store.Put("Luke", "Skywalker")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("does an insert when name does not exist in database", func() {
			fakeDbProvider.DbReturns(fakeDb, nil)

			err := store.Put("Luke", "Skywalker")
			Expect(err).To(BeNil())

			Expect(fakeDb.ExecCallCount()).To(Equal(1))

			query, values := fakeDb.ExecArgsForCall(0)
			Expect(query).To(Equal("INSERT INTO configurations (name, value) VALUES(?,?)"))

			Expect(values[0]).To(Equal("Luke"))
			Expect(values[1]).To(Equal("Skywalker"))
		})

		It("does an update when name exists in database", func() {
			fakeDb.ExecReturns(nil, errors.New("duplicate"))
			fakeDbProvider.DbReturns(fakeDb, nil)

			store.Put("Luke", "Skywalker")

			Expect(fakeDb.ExecCallCount()).To(Equal(2))

			query, values := fakeDb.ExecArgsForCall(0)
			Expect(query).To(Equal("INSERT INTO configurations (name, value) VALUES(?,?)"))
			Expect(values[0]).To(Equal("Luke"))
			Expect(values[1]).To(Equal("Skywalker"))

			query, values = fakeDb.ExecArgsForCall(1)
			Expect(query).To(Equal("UPDATE configurations SET value = ? WHERE name = ?"))
			Expect(values[0]).To(Equal("Skywalker"))
			Expect(values[1]).To(Equal("Luke"))
		})
	})

	Describe("Delete", func() {

		It("closes db connection on exit", func() {
			fakeDbProvider.DbReturns(fakeDb, nil)
			store.Delete("Luke")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		Context("Name exists", func() {

			BeforeEach(func() {
				fakeDbProvider.DbReturns(fakeDb, nil)
				fakeDb.ExecReturns(fakeResult, nil)

				fakeResult.RowsAffectedReturns(1, nil)
			})

			It("removes value", func() {
				store.Delete("Luke")

				Expect(fakeDb.ExecCallCount()).To(Equal(1))
				query, value := fakeDb.ExecArgsForCall(0)
				Expect(query).To(Equal("DELETE FROM configurations WHERE name = ?"))
				Expect(value[0]).To(Equal("Luke"))
			})

			It("returns true", func() {
				deleted, err := store.Delete("Luke")

				Expect(deleted).To(BeTrue())
				Expect(err).To(BeNil())
			})
		})

		Context("Name does not exist", func() {

			BeforeEach(func() {
				fakeDbProvider.DbReturns(fakeDb, nil)
				fakeDb.ExecReturns(fakeResult, nil)

				fakeResult.RowsAffectedReturns(0, nil)
			})

			It("returns false", func() {
				deleted, err := store.Delete("name")
				Expect(deleted).To(BeFalse())
				Expect(err).To(BeNil())
			})
		})
	})
})
