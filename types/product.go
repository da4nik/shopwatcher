package types

import (
	"reflect"
)

// Size - size availability
type Size struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

// User - user to sent back to
type User struct {
	ChatType string `json:"chat_type"`
	ChatID   string `json:"chat_id"`
}

// Product represents parsed product
type Product struct {
	URL       string `json:"url"`
	Shop      string `json:"shop"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Currency  string `json:"currency"`
	Available bool   `json:"available"`
	Sizes     []Size `json:"sizes"`
	Users     []User `json:"users"`
}

// Equal - checks what was changed
func (p Product) Equal(product Product) bool {
	return reflect.DeepEqual(p, product)
}
