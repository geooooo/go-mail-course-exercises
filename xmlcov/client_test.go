// TODO: По-хорошему нужно еще больше тестов, но для учебы пойдет

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"testing"
	"time"
)

var testClient = SearchClient{}

func formatError(message string, expected any, actual any) string {
	return fmt.Sprintf("%s\nExpected: %v\nActual: %v\n", message, expected, actual)
}

func TestLimitNegativeFindUsers(t *testing.T) {
	request := SearchRequest{Limit: -1}

	_, err := testClient.FindUsers(request)

	if !errors.Is(err, LimitError) {
		t.Error(formatError("Should be an error", LimitError, err))
	}
}

func TestLimitGreat25FindUsers(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		actualLimit := req.URL.Query().Get("limit")
		if actualLimit != "26" {
			t.Error(formatError("Limit param", "26", actualLimit))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := SearchClient{URL: server.URL}
	request := SearchRequest{Limit: 100}
	client.FindUsers(request)
}

func TestOffsetNegativeFindUsers(t *testing.T) {
	request := SearchRequest{Offset: -1}

	_, err := testClient.FindUsers(request)

	if !errors.Is(err, OffsetError) {
		t.Error(formatError("Should be an error", OffsetError, err))
	}
}

func TestRequestFullStructFindUsers(t *testing.T) {
	expectedUrl := "/url/to"
	expectedAccessToken := "token"
	expectedQuery := "field"
	limit := 1
	offset := 2
	order := OrderByAsc

	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		actualMethod := req.Method
		if actualMethod != "GET" {
			t.Error(formatError("Http method", "GET", actualMethod))
		}

		actualUrl := req.URL.Path
		if actualUrl != expectedUrl {
			t.Error(formatError("Url", expectedUrl, actualUrl))
		}

		actualAccessToken := req.Header.Get("AccessToken")
		if actualAccessToken != expectedAccessToken {
			t.Error(formatError("Access token header", expectedAccessToken, actualAccessToken))
		}

		actualLimit := req.URL.Query().Get("limit")
		expectedLimit := strconv.FormatInt(int64(limit+1), 10)
		if actualLimit != expectedLimit {
			t.Error(formatError("Limit param", expectedLimit, actualLimit))
		}

		actualOffset := req.URL.Query().Get("offset")
		expectedOffset := strconv.FormatInt(int64(offset), 10)
		if actualOffset != expectedOffset {
			t.Error(formatError("Offset param", expectedOffset, actualOffset))
		}

		actualQuery := req.URL.Query().Get("query")
		if actualQuery != expectedQuery {
			t.Error(formatError("Query param", expectedQuery, actualQuery))
		}

		actualOrderField := req.URL.Query().Get("order_field")
		if actualOrderField != expectedQuery {
			t.Error(formatError("Order field param", expectedQuery, actualOrderField))
		}

		actualOrderBy := req.URL.Query().Get("order_by")
		expectedOrderBy := strconv.FormatInt(int64(order), 10)
		if actualOrderBy != expectedOrderBy {
			t.Error(formatError("Order by param", expectedOrderBy, actualOrderBy))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := SearchClient{
		URL:         server.URL + expectedUrl,
		AccessToken: expectedAccessToken,
	}
	request := SearchRequest{
		Limit:      limit,
		Offset:     offset,
		OrderBy:    order,
		Query:      expectedQuery,
		OrderField: expectedQuery,
	}
	client.FindUsers(request)
}

func TestTimeoutErrorFindUsers(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// dummy impl
		time.Sleep(time.Duration(2) * time.Second)
		res.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := SearchClient{URL: server.URL}
	if _, err := client.FindUsers(SearchRequest{}); err == nil {
		t.Error(formatError("Should be timeout", "it does not matter", err))
	}
}

func TestServerErrorFindUsers(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusInternalServerError)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := SearchClient{URL: server.URL}
	if _, err := client.FindUsers(SearchRequest{}); err == nil {
		t.Error(formatError("Should catch server internal error", "it does not matter", err))
	}
}

func TestAuthErrorFindUsers(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusUnauthorized)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := SearchClient{
		URL:         server.URL,
		AccessToken: "token",
	}
	if _, err := client.FindUsers(SearchRequest{}); err == nil {
		t.Error(formatError("Should catch auth error", "it does not matter", err))
	}
}

func TestBadRequestErrorFindUsers(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := SearchClient{
		URL: server.URL,
	}
	if _, err := client.FindUsers(SearchRequest{}); err == nil {
		t.Error(formatError("Should catch bad request error", "it does not matter", err))
	}
}

func TestBaseResponseFindUsers(t *testing.T) {
	users := []User{
		{
			Id:     0,
			Name:   "John Peach",
			Age:    35,
			About:  "Top 10 streamer of the decade and skuff",
			Gender: "Beer elemental",
		},
		{
			Id:     1,
			Name:   "Hanz Lavanda",
			Age:    20,
			About:  "Student of MIT",
			Gender: "Apache helicopter",
		},
	}
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		data, _ := json.Marshal(users)
		res.Write(data)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := SearchClient{
		URL: server.URL,
	}

	response, _ := client.FindUsers(SearchRequest{Limit: 1})
	if response.NextPage != true {
		t.Error(formatError("exists next page", true, response.NextPage))
	}
	if !slices.Equal(users[0:1], response.Users) {
		t.Error(formatError("list of users", users, response.Users))
	}

	response, _ = client.FindUsers(SearchRequest{Limit: 2})
	if response.NextPage != false {
		t.Error(formatError("not exists next page", false, response.NextPage))
	}
	if !slices.Equal(users, response.Users) {
		t.Error(formatError("list of users", users, response.Users))
	}
}
