// Package model описывает доменные сущности каталога деталей.
package model

import "time"

// Category — категория детали.
type Category int32

const (
	// CategoryUnspecified — категория не указана.
	CategoryUnspecified Category = iota
	// CategoryEngine — двигатель.
	CategoryEngine
	// CategoryFuel — топливо.
	CategoryFuel
	// CategoryPorthole — иллюминатор.
	CategoryPorthole
	// CategoryWing — крыло.
	CategoryWing
)

// Dimensions — габариты и вес детали.
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

// Manufacturer — производитель детали.
type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// Value — значение произвольного типа в метаданных детали.
// Заполнено ровно одно поле, остальные — nil.
type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
	StringList  []string
}

// Part — деталь космического корабля.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

// PartsFilter — фильтр выборки деталей. Пустые поля не сужают выборку,
// заполненные складываются по «И», а значения внутри одного поля — по «ИЛИ».
type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
