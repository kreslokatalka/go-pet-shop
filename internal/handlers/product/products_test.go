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
}

func TestCreateProduct_BadRequest(t *testing.T) {
	mock := mocks.NewProducts(t)

	product := &models.Product{}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPost, "/products", bodyReader)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("exspected 400, got %d", w.Code)
	}
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
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	mock := mocks.NewProducts(t)

	product := &models.Product{}
	data, _ := json.Marshal(product)
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
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	mock := mocks.NewProducts(t)

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("exspected 400, got %d", w.Code)
	}
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
}
