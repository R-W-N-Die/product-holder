package storage

// Product - структура для хранения информации о продукте
// Структура в Go - это тип данных, который группирует поля под одним именем
type Product struct {
	ID         uint32 `json:"id"`          // Id продукта
	Name       string `json:"name"`        // Название продукта
	Price      uint32 `json:"price"`       // Цена продукта в рублях
	SoldAmount uint32 `json:"sold_amount"` // На сколько рублей продано
}

// Конструктор для Product (функция, которая создает новый экземпляр)
// NewProduct создает новый продукт с заданными параметрами
func NewProduct(id uint32, name string, price uint32) *Product {
	// &Product{...} - создает структуру и возвращает указатель на нее
	// Указатель (*Product) позволяет работать с оригинальной структурой, а не с копией
	return &Product{
		ID:         id,
		Name:       name,
		Price:      price,
		SoldAmount: 0, // Изначально продано на 0 рублей
	}
}
