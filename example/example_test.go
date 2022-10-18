package example_test

import (
	"github.com/derfenix/goi18n/example"
)

func ExampleBasic() {
	example.Basic()

	// Output:
	// Тест теста
	// Test of the beer
	// всего пара пауков
	// Использован русский язык, кириллица
	// just pair of spiders
	// Использован английский язык, латиница
}
