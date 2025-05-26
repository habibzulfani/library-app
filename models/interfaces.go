package models

type UserInterface interface {
	GetID() uint
	GetName() string
	GetNIM() string
	// Add other methods as needed
}

type PaperInterface interface {
	GetID() uint
	// Add other methods as needed
}

type BookInterface interface {
	GetID() uint
	// Add other methods as needed
}
