package storage

import (
	"sync"
)

// Storage - основное хранилище продуктов
// Структура хранит состояние нашего приложения
type Storage struct {
	mu       sync.RWMutex        // Мьютекс для синхронизации доступа (подробнее в Этапе 3)
	products [1_000_000]*Product // Статический массив указателей на Product
	// [размер]тип - массив фиксированного размера
	// *Product - указатель на Product (экономит память)
}

// NewStorage создает и инициализирует новое хранилище
func NewStorage() *Storage {
	// Создаем хранилище
	storage := &Storage{}

	// Инициализируем массив продуктов (пока пустыми значениями)
	// В Go при создании массива все элементы инициализируются "нулевыми значениями":
	// - для указателей это nil
	// - для чисел это 0
	// - для строк это ""

	return storage
}

// GetProduct возвращает продукт по ID
func (s *Storage) GetProduct(id uint32) *Product {
	// Проверяем, что ID в допустимом диапазоне
	if id >= uint32(len(s.products)) {
		return nil // Возвращаем nil если ID вне диапазона
	}

	s.mu.RLock()         // Блокировка для чтения (множество горутин может читать одновременно)
	defer s.mu.RUnlock() // defer - отложенное выполнение, гарантирует разблокировку при выходе из функции

	return s.products[id]
}

// GetProductCount возвращает количество существующих продуктов
func (s *Storage) GetProductCount() int {
	s.mu.RLock()         // Блокируем для чтения
	defer s.mu.RUnlock() // Не забываем разблокировать

	count := 0
	// Проходим по всем продуктам
	for _, product := range s.products {
		if product != nil && product.Name != "" {
			// Если продукт существует и у него есть имя
			count++
		}
	}
	return count
}

// UpdateProduct обновляет информацию о продукте
func (s *Storage) UpdateProduct(id uint32, updateFunc func(*Product)) bool {
	// Проверяем диапазон
	if id >= uint32(len(s.products)) {
		return false
	}

	s.mu.Lock()         // Блокировка для записи (только одна горутина может писать)
	defer s.mu.Unlock() // Гарантированная разблокировка

	// Получаем продукт
	product := s.products[id]

	// Если продукт не существует, создаем новый
	if product == nil {
		product = NewProduct(id, "", 0)
		s.products[id] = product
	}

	// Вызываем функцию обновления
	updateFunc(product)

	return true
}
