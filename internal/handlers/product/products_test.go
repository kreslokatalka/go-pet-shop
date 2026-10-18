package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-pet-shop/internal/handlers/product/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	testifymock "github.com/stretchr/testify/mock"
)

var (
	ErrInternalError = errors.New("internal error")
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	mock := mocks.NewProducts(t)
	products := []models.Product{
		{ID: 1, Name: "Dog Food"},
	}
	mock.On("GetAllProducts", testifymock.Anything).
		Return(products, nil).
		Once()

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)
	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	assertJSONBody(t, w, products)
}
func TestGetAllProducts_Error(t *testing.T) {
	mock := mocks.NewProducts(t)

	mock.On("GetAllProducts", testifymock.Anything).
		Return(nil, ErrInternalError).
		Once()

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Internal server error",
		"message": "Failed to retrieve products",
	})
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	mock := mocks.NewProducts(t)

	mock.On("CreateProduct", testifymock.Anything, testifymock.Anything).
		Return(10, nil).
		Once()

	product := &models.Product{ID: 40, Name: "name", Price: 40, Stock: 40}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPost, "/products", bodyReader)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("exspected 200, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]interface{}{
		"status": "Product created successfully",
		"id":     10,
		"product": models.Product{
			ID: 10, Name: "name", Price: 40, Stock: 40,
		},
	})
}

func TestCreateProduct_BadRequest(t *testing.T) {
	mock := mocks.NewProducts(t)

	data, _ := json.Marshal(`ID= 40, Name "name Price: 40, Stock: 40}`)
	bodyReader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPost, "/products", bodyReader)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("exspected 400, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Bad request",
		"message": "Invalid JSON payload",
	})
}

func TestCreateProduct_Fail(t *testing.T) {
	mock := mocks.NewProducts(t)
	mock.On("CreateProduct", testifymock.Anything, testifymock.Anything).
		Return(0, ErrInternalError).
		Once()

	product := &models.Product{ID: 40, Name: "name", Price: 40, Stock: 40}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPost, "/products", bodyReader)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("exspected 500, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Internal server error",
		"message": "Failed to create product",
	})
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {

	mock := mocks.NewProducts(t)

	mock.On("UpdateProduct", testifymock.Anything, testifymock.Anything).
		Return(nil).
		Once()

	product := &models.Product{ID: 40, Name: "name", Price: 40, Stock: 40}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPut, "/products/40", bodyReader)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	router := chi.NewRouter()
	router.Put("/products/{id}", handler.UpdateProduct)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("exspected 200, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]interface{}{
		"status": "Product updated successfully",
		"id":     40,
		"product": models.Product{
			ID: 40, Name: "name", Price: 40, Stock: 40,
		},
	})
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	mock := mocks.NewProducts(t)

	data, _ := json.Marshal(`ID= 40, Name "name Price: 40, Stock: 40}`)
	bodyReader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPut, "/products/40", bodyReader)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	router := chi.NewRouter()
	router.Put("/products/{id}", handler.UpdateProduct)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("exspected 400, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Bad request",
		"message": "Invalid JSON payload",
	})
}

func TestUpdateProduct_Fail(t *testing.T) {
	mock := mocks.NewProducts(t)

	mock.On("UpdateProduct", testifymock.Anything, testifymock.Anything).
		Return(ErrInternalError).
		Once()

	product := &models.Product{ID: 40, Name: "name", Price: 40, Stock: 40}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Put("/products/{id}", handler.UpdateProduct)

	req := httptest.NewRequest(http.MethodPut, "/products/40", bodyReader)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("exspected 500, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Internal server error",
		"message": "Failed to update product",
	})
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	mock := mocks.NewProducts(t)

	mock.On("DeleteProduct", testifymock.Anything, testifymock.AnythingOfType("int")).
		Return(nil).
		Once()

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/40", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("exspected 200, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]interface{}{
		"status":  "Product deleted successfully",
		"id":      40,
		"message": "Product with ID 40 has been deleted",
	})
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	mock := mocks.NewProducts(t)

	req := httptest.NewRequest(http.MethodDelete, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("exspected 400, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Bad request",
		"message": "Product ID is required",
	})
}

func TestDeleteProduct_Fail(t *testing.T) {
	mock := mocks.NewProducts(t)

	mock.On("DeleteProduct", testifymock.Anything, testifymock.AnythingOfType("int")).
		Return(ErrInternalError).
		Once()

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/40", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Internal server error",
		"message": "Failed to delete product",
	})
}

func assertJSONBody(t *testing.T, w *httptest.ResponseRecorder, want interface{}) {
	t.Helper()

	expected, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("failed to marshal expected response body: %v", err)
	}

	var actual bytes.Buffer
	if err := json.Compact(&actual, w.Body.Bytes()); err != nil {
		t.Fatalf("failed to compact response body: %v", err)
	}
	if !bytes.Equal(actual.Bytes(), expected) {
		t.Fatalf("unexpected response body: got %s, want %s", actual.Bytes(), expected)
	}
}
