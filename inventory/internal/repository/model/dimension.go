package model

type Dimensions struct {
	Length float64 `bson:"length"` // Длина в см
	Width  float64 `bson:"width"`  // Ширина в см
	Height float64 `bson:"height"` // Высота в см
	Weight float64 `bson:"weight"` // Вес в кг
}
