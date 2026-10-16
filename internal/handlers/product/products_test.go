package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

var (
	ErrInternalError = errors.New("internal error")
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: 1, Name: "Dog Food"},
			}, nil
		},
	}

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
	assertJSONBody(t, w, []models.Product{{ID: 1, Name: "Dog Food"}})
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("DB error")
		},
	}

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
	// TODO: Написать Unit-тест для создания продукта (200 OK)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 10, nil
		},
	}

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
	// TODO: Написать Unit-тест для создания продукта с невалидным JSON (400 Bad Request)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 10, nil
		},
	}

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
	assertJSONBody(t, w, map[string]string{
		"error":   "Bad request",
		"message": "Product name is required",
	})
}

func TestCreateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта при ошибке сервиса (500 Internal Server Error)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 0, ErrInternalError
		},
	}

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
	// TODO: Написать Unit-тест для обновления продукта (200 OK)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}

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
	// TODO: Написать Unit-тест для обновления продукта с невалидным JSON (400 Bad Request)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}

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
	assertJSONBody(t, w, map[string]string{
		"error":   "Bad request",
		"message": "Product name is required",
	})
}

func TestUpdateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта при ошибке сервиса (500 Internal Server Error)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return ErrInternalError
		},
	}

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
	// TODO: Написать Unit-тест для удаления продукта (200 OK)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	product := &models.Product{ID: 40, Name: "name", Price: 40, Stock: 40}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/40", bodyReader)
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
	// TODO: Написать Unit-тест для удаления продукта с пустым id (400 Bad Request)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	product := &models.Product{Name: "name", Price: 40, Stock: 40}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/abc", bodyReader)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("exspected 400, got %d", w.Code)
	}
	assertJSONBody(t, w, map[string]string{
		"error":   "Bad request",
		"message": "Product ID must be a number",
	})
}

func TestDeleteProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта при ошибке сервиса (500 Internal Server Error)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return ErrInternalError
		},
	}

	product := &models.Product{ID: 40, Name: "name", Price: 40, Stock: 40}
	data, _ := json.Marshal(product)
	bodyReader := bytes.NewReader(data)

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)

	req := httptest.NewRequest(http.MethodDelete, "/products/40", bodyReader)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("exspected 200, got %d", w.Code)
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
