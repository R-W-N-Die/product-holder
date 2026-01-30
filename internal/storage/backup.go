package storage

import (
	"encoding/json" // Пакет для работы с JSON
	"os"            // Пакет для работы с файловой системой
	"time"          // Пакет для работы со временем
)

// BackupInfo содержит информацию о бэкапе
type BackupInfo struct {
	Timestamp time.Time      `json:"timestamp"`
	Products  []*ProductInfo `json:"products"`
}

// ProductInfo упрощенная структура для бэкапа
type ProductInfo struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	Price      uint32 `json:"price"`
	SoldAmount uint32 `json:"sold_amount"`
}

// SaveBackup сохраняет бэкап в файл
func (s *Storage) SaveBackup(filename string) error {
	s.mu.RLock() // Блокируем для чтения
	defer s.mu.RUnlock()

	// Подготавливаем данные для бэкапа
	backup := BackupInfo{
		Timestamp: time.Now(),
		Products:  make([]*ProductInfo, 0),
	}

	// Собираем только существующие продукты
	for id, product := range s.products {
		if product != nil {
			backup.Products = append(backup.Products, &ProductInfo{
				ID:         uint32(id),
				Name:       product.Name,
				Price:      product.Price,
				SoldAmount: product.SoldAmount,
			})
		}
	}

	// Создаем или перезаписываем файл
	// os.O_CREATE - создать если не существует
	// os.O_WRONLY - только запись
	// os.O_TRUNC - очистить файл если существует
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // Закрываем файл при выходе из функции

	// Создаем JSON encoder
	encoder := json.NewEncoder(file)

	// Кодируем данные в JSON и записываем в файл
	return encoder.Encode(backup)
}

// LoadBackup загружает бэкап из файла
func (s *Storage) LoadBackup(filename string) error {
	// Открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		// Если файла нет - это нормально при первом запуске
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	// Декодируем JSON
	var backup BackupInfo
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&backup); err != nil {
		return err
	}

	s.mu.Lock() // Блокируем для записи
	defer s.mu.Unlock()

	// Восстанавливаем продукты
	for _, productInfo := range backup.Products {
		if productInfo.ID < uint32(len(s.products)) {
			s.products[productInfo.ID] = &Product{
				ID:         productInfo.ID,
				Name:       productInfo.Name,
				Price:      productInfo.Price,
				SoldAmount: productInfo.SoldAmount,
			}
		}
	}

	return nil
}
