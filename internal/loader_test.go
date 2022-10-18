package internal_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/derfenix/goi18n/internal"
)

var vvv = `[
  {
    "key": "test",
    "description": "Для тестов, не трогать",
    "translation": "Тест %s"
  },
  {
    "key": "test plural",
    "description": "Для тестов, не трогать",
    "plural": {
      "other": "всего %d пауков",
      "one": "паучок",
      "=0": "нет пауков",
      "=2": "всего пара пауков"
    }
  },
  {
    "key": "transition_not_allowed",
    "description": "Ошибка при попытке произвести запрещённое изменение состояния",
    "translation": "Переход запрещён"
  },
  {
    "key": "Should be shorter than %d symbols",
    "translation": "Должно быть короче %d символов"
  },
  {
    "key": "Should be longer than %d symbols",
    "translation": "Должно быть длиннее %d символов"
  }
]`

func Benchmark_plurals_UnmarshalJSON(b *testing.B) {
	var trans []Translation

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		require.NoError(b, json.Unmarshal([]byte(vvv), &trans))
	}
}
