package main

import (
	"log"  // Пакет для логирования
	"time" // Пакет для работы со временем

	"product-holder-project/internal/api"     // Импорт API
	"product-holder-project/internal/storage" // Импорт хранилища
)

func main() {
	log.Println("Запуск Product Holder...")

	// 1. Создаем хранилище
	storage := storage.NewStorage()

	// 2. Загружаем бэкап (если есть)
	backupFile := "backup/products.json"
	if err := storage.LoadBackup(backupFile); err != nil {
		log.Printf("Ошибка загрузки бэкапа: %v", err)
	} else {
		log.Println("Бэкап успешно загружен")
	}

	// 3. Запускаем периодическое сохранение бэкапа (раз в час)
	storage.StartBackupScheduler(backupFile, time.Hour)

	// 4. Добавляем тестовые данные (для демонстрации)
	initTestData(storage)

	// 5. Создаем и запускаем HTTP сервер
	server := api.NewServer(storage, "8082")

	log.Println("Сервис готов к работе")
	if err := server.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

// initTestData добавляет тестовые данные в хранилище
func initTestData(s *storage.Storage) {
	// Функция для безопасного обновления продукта
	updateProduct := func(id uint32, name string, price uint32) {
		s.UpdateProduct(id, func(p *storage.Product) {
			p.Name = name
			p.Price = price
		})
	}

	// Добавляем несколько тестовых продуктов
	updateProduct(1, "Ноутбук Lenovo", 75000)
	updateProduct(2, "Смартфон iPhone", 90000)
	updateProduct(3, "Наушники Sony", 15000)
	updateProduct(4, "Клавиатура Logitech", 5000)
	updateProduct(5, "Монитор Samsung", 30000)
	updateProduct(6, "Игровая мышь Razer", 8000)
	updateProduct(7, "Планшет iPad", 60000)
	updateProduct(8, "Умные часы Apple Watch", 40000)
	updateProduct(9, "Фитнес-браслет Xiaomi", 3000)
	updateProduct(10, "Внешний жесткий диск", 7000)

	log.Println("Добавлено 10 тестовых продуктов")
}
