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
}
