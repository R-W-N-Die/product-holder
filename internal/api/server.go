package api

import (
	"encoding/json" // Пакет для работы с JSON
	"log"           // Пакет для логирования
	"net/http"      // Пакет для работы с HTTP
	"strconv"       // Пакет для конвертации строк
	"time"

	"product-holder-project/internal/storage" // Импорт нашего хранилища
)

// Server представляет HTTP сервер
type Server struct {
	storage *storage.Storage
	port    string
}

// NewServer создает новый HTTP сервер
func NewServer(storage *storage.Storage, port string) *Server {
	return &Server{
		storage: storage,
		port:    port,
	}
}

// productHandler обрабатывает запросы к /api/v1/product
func (s *Server) productHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS (для тестирования)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Получаем параметр id из query string
	query := r.URL.Query()
	idStr := query.Get("id")

	if idStr == "" {
		// Если id не указан, возвращаем все продукты
		s.getAllProducts(w)
		return
	}

	// Конвертируем строку в число
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "Invalid product ID"}`, http.StatusBadRequest)
		return
	}

	// Получаем продукт из хранилища
	product := s.storage.GetProduct(uint32(id))
	if product == nil {
		http.Error(w, `{"error": "Product not found"}`, http.StatusNotFound)
		return
	}

	// Кодируем продукт в JSON и отправляем
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}

// productCountHandler обрабатывает запросы к /api/v1/product/count
func (s *Server) productCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Теперь используем новый метод
	count := s.storage.GetProductCount()

	response := map[string]interface{}{
		"count":     count,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// getAllProducts возвращает все существующие продукты
func (s *Server) getAllProducts(w http.ResponseWriter) {
	// В реальном приложении здесь была бы пагинация
	// Для простоты возвращаем только первые 1000 продуктов

	products := make([]*storage.Product, 0)

	for i := 0; i < 1000; i++ {
		if product := s.storage.GetProduct(uint32(i)); product != nil && product.Name != "" {
			products = append(products, product)
		}
	}

	// Кодируем в JSON
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
	}
}

// healthHandler для проверки здоровья сервиса
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	// Регистрируем обработчики
	http.HandleFunc("/api/v1/product", s.productHandler)
	http.HandleFunc("/api/v1/product/count", s.productCountHandler)
	http.HandleFunc("/health", s.healthHandler)

	// Запускаем сервер
	log.Printf("HTTP сервер запущен на порту %s", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
