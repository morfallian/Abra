package Slices_and_maps

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"
)

//
func TestCalculateOrder(t *testing.T) {
	dbOrders := make(dbOrders)
	testCases := []struct {
		name               string
		shop               Products
		order              order
		expectedTotalPrice float64
		expectedError      error
	}{
		{
			"basic case. nil slice",
			Products{},
			nil,
			0,
			nil,
		},
		{
			"basic case. empty slice",
			Products{},
			order{},
			0,
			listOrderInitial,
		},
		{
			"basic case. single element slice",
			Products{
				"a": 1.0,
				"b": 10.0,
			},
			order{"a"},
			1.0,
			nil,
		},
		{
			"basic case. two element slice",
			Products{
				"a": 1.0,
				"b": 10.0,
			},
			order{"b", "a"},
			11,
			nil,
		},
		{
			"basic case. single unknown item",
			Products{
				"a": 1.0,
				"b": 10.0,
			},
			order{"xxx"},
			0,
			errItemNoFound,
		},
		{
			"basic case. single unknown item inbetween ",
			Products{
				"a": 1.0,
				"b": 10.0,
			},
			order{"a", "xxx", "b"},
			0,
			errItemNoFound,
		},
		{
			"partial match",
			Products{
				"a": 1.0,
				"b": 10.0,
			},
			order{"aa"},
			0,
			errItemNoFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := dbOrders.calculateOrdersCost(tc.order, tc.shop)
			if !errors.Is(err, tc.expectedError ) {
				t.Fatalf("got\t\t%v\nwant\t%v", err, tc.expectedError )
			}
			if !reflect.DeepEqual(res, tc.expectedTotalPrice) {
				t.Fatalf("got\t\t%v\nwant\t%v", res, tc.expectedTotalPrice)
			}
		})
	}

}
//
func TestCalculateOrderWithCache(t *testing.T) {
	testCases := []struct {
		name               string
		shop               Products
		order              order
		expectedTotalPrice float64
		expectedError      error
	}{
		{
			"basic case. nil slice",
			Products{},
			nil,
			0,
			nil,
		},
		{
			"basic case. empty slice",
			Products{},
			order{},
			0,
			nil,
		},
		{
			"basic case. single element slice",
			Products{
				"a": 1,
				"b": 10,
			},
			order{"a"},
			1,
			nil,
		},
		{
			"basic case. two element slice",
			Products{
				"a": 1,
				"b": 10,
			},
			order{"b", "a"},
			11,
			nil,
		},
		{
			"basic case. single unknown item",
			Products{
				"a": 1,
				"b": 10,
			},
			order{"xxx"},
			0,
			errItemNoFound,
		},
		{
			"basic case. single unknown item inbetween ",
			Products{
				"a": 1,
				"b": 10,
			},
			order{"a", "xxx", "b"},
			0,
			errItemNoFound,
		},
		{
			"partial match",
			Products{
				"a": 1,
				"b": 10,
			},
			order{"aa"},
			0,
			errItemNoFound,
		},
	}
	dborders := make(dbOrders)
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for i := 0; i < 5; i++ {
				res, err := dborders.calculateOrdersCost(tc.order, tc.shop)

				if !errors.Is(err, tc.expectedError) {
					t.Fatalf("got\t\t%v\nwant\t%v", err, tc.expectedError)
				}
				if tc.expectedError == nil && i > 0 {
					sort.Strings(tc.order)
					key := strings.Join(tc.order[:],"")
					totalFromCache, ok := dborders[key]
					if !ok {
						t.Fatalf("cache was not used: %v", key)
					}
					if totalFromCache != tc.expectedTotalPrice {
						t.Fatalf("got in cache\t\t%v\nwant\t%v", totalFromCache, tc.expectedTotalPrice)
					}
				}
				if !reflect.DeepEqual(res, tc.expectedTotalPrice) {
					t.Fatalf("got\t\t%v\nwant\t%v", res, tc.expectedTotalPrice)
				}
			}
		})
	}
}
//
func TestAddItem(t *testing.T) {
	testCases := []struct {
		name          string
		shop          Products
		item          string
		price         float64
		expectedShop  Products
		expectedError error
	}{
		{
			"basic case. empty Item",
			Products{},
			"",
			0.0,
			Products{},
			errEmptyItem,

		},
		{
			"basic case. empty item name",
			Products{},
			"",
			10,
			Products{},
			errEmptyItem,
		},
		{
			"basic case. already exists",
			Products{
				"a": 1,
				"b": 10,
			},
			"a",
			1,
			Products{
				"a": 1,
				"b": 10,
			},
			errItemAlreadyExists,
		},
		{
			"basic case. correct item",
			Products{
				"a": 1,
				"b": 10},
			"xxx",
			100,
			Products{
				"a":   1,
				"b":   10,
				"xxx": 100,
			},
			nil,
		},
	}

	//products := make(Products)
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.shop.AppendNewProducts(tc.item, tc.price)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("got\t\t%v\nwant\t%v", err, tc.expectedError)
			}
			if !reflect.DeepEqual(tc.shop, tc.expectedShop) {
				t.Fatalf("got\t\t%v\nwant\t%v", tc.shop, tc.expectedShop)
			}
		})
	}
}
//
func TestChangePrice(t *testing.T) {
	testCases := []struct {
		name          string
		shop          Products
		item          string
		price         float64
		expectedShop  Products
		expectedError error
	}{
		{
			"basic case. empty Item",
			Products{},
			"",
			0,
			Products{},
			errEmptyItem,
		},
		{
			"basic case. empty item name",
			Products{},
			"",
			10,
			Products{},
			errEmptyItem,
		},
		{
			"basic case. correct price change",
			Products{
				"a": 1,
				"b": 10,
			},
			"a",
			10,
			Products{
				"a": 10,
				"b": 10,
			},
			nil,
		},
		{
			"basic case. correct price change",
			Products{
				"a": 1,
				"b": 10,
			},
			"xxx",
			101,
			Products{
				"a": 1,
				"b": 10,
			},
			errItemNoFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.shop.changePriceProduct(tc.item, tc.price)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("got\t\t%v\nwant\t%v", err, tc.expectedError)
			}
			if !reflect.DeepEqual(tc.shop, tc.expectedShop) {
				t.Fatalf("got\t\t%v\nwant\t%v", tc.shop, tc.expectedShop)
			}
		})
	}
}
//
func TestChangeName(t *testing.T) {
	testCases := []struct {
		name          string
		shop          Products
		itemName      string
		newItemName   string
		expectedShop  Products
		expectedError error
	}{
		{
			"basic case. empty Item",
			Products{},
			"a",
			"",
			Products{},
			errEmptyItem,
		},
		{
			"basic case. empty item name",
			Products{
				"a": 10,
			},
			"a",
			"",
			Products{
				"a": 10,
			},
			errEmptyItem,
		},
		{
			"basic case. already exists",
			Products{
				"a": 1,
				"b": 10,
			},
			"a",
			"aa",
			Products{
				"aa": 1,
				"b":  10,
			},
			nil,
		},
		{
			"basic case. already exists",
			Products{
				"a": 1,
				"b": 10,
			},
			"xxx",
			"aa",
			Products{
				"a": 1,
				"b": 10,
			},
			errItemNoFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.shop.changeNameproduct(tc.newItemName, tc.itemName)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("got\t\t%v\nwant\t%v", err, tc.expectedError)
			}
			if !reflect.DeepEqual(tc.shop, tc.expectedShop) {
				t.Fatalf("got\t\t%v\nwant\t%v", tc.shop, tc.expectedShop)
			}
		})
	}
}
//
func TestAddAccount(t *testing.T) {
	testCases := []struct {
		name             string
		accounts         accountList
		account          account
		expectedAccounts accountList
		expectedError    error
	}{
		{
			"basic case. empty Account",
			accountList{},
			account{},
			accountList{},
			errEmptyAccount,
		},
		{
			"basic case. empty Account name",
			accountList{},
			account{"", 10},
			accountList{},
			errEmptyAccount,
		},
		{
			"basic case. already exists",
			accountList{
				"a": {"a", 1},
				"b": {"b", 10},
			},
			account{
				uname: "a",
				money: 10,
			},
			accountList{
				"a": {"a", 1},
				"b": {"b", 10},
			},
			errAccountAlreadyExists,
		},
		{
			"basic case. correct item",
			accountList{
				"a": {"a", 1},
				"b": {"b", 10},
			},
			account{"xxx", 100},
			accountList{
				"a":   {"a", 1},
				"b":   {"b", 10},
				"xxx": {"xxx", 100},
			},
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.accounts.AddNewAccount(tc.account.uname, tc.account)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("got\t\t%v\nwant\t%v", err, tc.expectedError)
			}
			if !reflect.DeepEqual(tc.accounts, tc.expectedAccounts) {
				t.Fatalf("got\t\t%v\nwant\t%v", tc.account, tc.expectedAccounts)
			}
		})
	}
}
//
func TestSortAccounts(t *testing.T) {
	testCases := []struct {
		name             string
		accounts         accountList
		sortBy           int
		expectedAccounts []tempSortStruct
		expectedError    error
	}{
		// name
		{
			"SortByNameAsc. empty Accounts",
			accountList{},
			1,
			[]tempSortStruct{},
			nil,
		},
		{
			"SortByNameAsc. single Account ",
			accountList{
				"a": {"a", 10},
			},
			1,
			[]tempSortStruct{{"a", account{
				uname: "a",
				money: 10,
			}}},
			nil,
		},
		{
			"SortByNameAsc. already sorted",
			accountList{
				"a":     {"a", 1},
				"b":     {"b", 10},
				"d":     {"d", 12},
				"d10":   {"d10", 11},
				"xxx_1": {"xxx_1", 22},
			},
			1,
			[]tempSortStruct{
				{"a", account{
					uname: "a",
					money: 1,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
			},
			nil,
		},
		{
			"SortByNameAsc. already sorted in reverse order",
			accountList{
				"xxx_1": {"xxx_1", 22},
				"d10":   {"d10", 11},
				"d":     {"d", 12},
				"b":     {"b", 10},
				"a":     {"a", 1},
			},
			1,
			[]tempSortStruct{
				{"a", account{
					uname: "a",
					money: 1,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
			},
			nil,
		},
		{
			"SortByNameAsc. random order",
			accountList{
				"d10":   {"d10", 11},
				"a":     {"a", 1},
				"xxx_1": {"xxx_1", 22},
				"b":     {"b", 10},
				"d":     {"d", 12},
			},
			1,
			[]tempSortStruct{
				{"a", account{
					uname: "a",
					money: 1,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
			},
			nil,
		},

		{
			"SortByNameDesc. empty Accounts",
			accountList{},
			2,
			[]tempSortStruct{},
			nil,
		},
		{
			"SortByNameDesc. single Account ",
			accountList{
				"a": {"a", 10},
			},
			2,
			[]tempSortStruct{{"a", account{
				uname: "a",
				money: 10,
			}}},
			nil,
		},
		{
			"SortByNameDesc. already sorted",
			accountList{
				"xxx_1": {"xxx_1", 22},
				"d10":   {"d10", 11},
				"d":     {"d", 12},
				"b":     {"b", 10},
				"a":     {"a", 1},
			},
			2,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},
		{
			"SortByNameDesc. already sorted in reverse order",
			accountList{
				"a":     {"a", 1},
				"b":     {"b", 10},
				"d":     {"d", 12},
				"d10":   {"d10", 11},
				"xxx_1": {"xxx_1", 22},
			},
			2,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},
		{
			"SortByNameDesc. random order",
			accountList{
				"d10":   {"d10", 11},
				"a":     {"a", 1},
				"xxx_1": {"xxx_1", 22},
				"b":     {"b", 10},
				"d":     {"d", 12},
			},
			2,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},

		// balance
		{
			"SortByBalanceDesc. empty Accounts",
			accountList{},
			3,
			[]tempSortStruct{},
			nil,
		},
		{
			"SortByBalanceDesc. single Account",
			accountList{
				"a": {"a", 10},
			},
			3,
			[]tempSortStruct{{"a", account{
				uname: "a",
				money: 10,
			}}},
			nil,
		},
		{
			"SortByBalanceDesc. already sorted",
			accountList{
				"xxx_1": {"xxx_1", 22},
				"d":     {"d", 12},
				"d10":   {"d10", 11},
				"b":     {"b", 10},
				"a":     {"a", 1},
			},
			3,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},
		{
			"SortByBalanceDesc. already sorted with duplicates",
			accountList{
				"xxx_1": {"xxx_1", 22},
				"xxx_2": {"xxx_2", 22},
				"d":     {"d", 12},
				"d11":   {"d11", 11},
				"d10":   {"d10", 11},
				"b":     {"b", 10},
				"a":     {"a", 1},
			},
			3,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"xxx_2", account{
					uname: "xxx_2",
					money: 22,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"d11", account{
					uname: "d11",
					money: 11,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},
		{
			"SortByBalanceDesc. already sorted in reverse order",
			accountList{
				"a":     {"a", 1},
				"b":     {"b", 10},
				"d10":   {"d10", 11},
				"d":     {"d", 12},
				"xxx_1": {"xxx_1", 22},
			},
			3,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},
		{
			"SortByBalanceDesc. already sorted in reverse order with duplicated",
			accountList{
				"a":     {"a", 1},
				"b":     {"b", 10},
				"d10":   {"d10", 11},
				"d11":   {"d11", 11},
				"d":     {"d", 12},
				"xxx_1": {"xxx_1", 22},
				"xxx_2": {"xxx_2", 22},
			},
			3,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"xxx_2", account{
					uname: "xxx_2",
					money: 22,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"d11", account{
					uname: "d11",
					money: 11,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},
		{
			"SortByBalanceDesc. random order",
			accountList{
				"d10":   {"d10", 11},
				"a":     {"a", 1},
				"xxx_1": {"xxx_1", 22},
				"b":     {"b", 10},
				"d":     {"d", 12},
			},
			3,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},
			},
			nil,
		},
		{
			"SortByBalanceDesc. random order with duplicates",
			accountList{
				"d10":   {"d10", 11},
				"a":     {"a", 1},
				"a1":    {"a1", 1},
				"xxx_1": {"xxx_1", 22},
				"xxx_2": {"xxx_2", 22},
				"b":     {"b", 10},
				"d":     {"d", 12},
			},
			3,
			[]tempSortStruct{
				{"xxx_1", account{
					uname: "xxx_1",
					money: 22,
				}},
				{"xxx_2", account{
					uname: "xxx_2",
					money: 22,
				}},
				{"d", account{
					uname: "d",
					money: 12,
				}},
				{"d10", account{
					uname: "d10",
					money: 11,
				}},
				{"b", account{
					uname: "b",
					money: 10,
				}},
				{"a1", account{
					uname: "a1",
					money: 1,
				}},
				{"a", account{
					uname: "a",
					money: 1,
				}},

			},
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := tc.accounts.SortAcc(tc.sortBy)

			//if !errors.Is(err, tc.expectedError) {
			//	t.Errorf("got\t\t%v\nwant\t%v", err, tc.expectedError)
			//}
			if !reflect.DeepEqual(res, tc.expectedAccounts) {
				t.Fatalf("got\t\t%v\nwant\t%v", res, tc.expectedAccounts)
			}
		})
	}
}

