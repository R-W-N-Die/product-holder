package storage

import (
	"log"  // Пакет для логирования
	"time" // Пакет для работы со временем
)

// StartBackupScheduler запускает периодическое сохранение бэкапов
func (s *Storage) StartBackupScheduler(backupFile string, interval time.Duration) {
	// Ticker - тикер отправляет значение в канал через заданные интервалы
	ticker := time.NewTicker(interval)

	// Запускаем горутину для обработки тиков
	go func() {
		for range ticker.C { // Читаем из канала тикера
			if err := s.SaveBackup(backupFile); err != nil {
				log.Printf("Ошибка сохранения бэкапа: %v", err)
			} else {
				log.Printf("Бэкап успешно сохранен в %s", backupFile)
			}
		}
	}()
}
