package internal

import (
	"io/fs"
	"testing/fstest"
	"time"
)

var TestFS = fstest.MapFS{
	"locales":    &fstest.MapFile{Mode: 0777 | fs.ModeDir},
	"locales/ru": &fstest.MapFile{Mode: 0777 | fs.ModeDir},
	"locales/en": &fstest.MapFile{Mode: 0777 | fs.ModeDir},
	"locales/ru/active.json": &fstest.MapFile{
		Data: []byte(`[
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
  }
]`),
		Mode:    0555,
		ModTime: time.Now(),
		Sys:     1,
	},
	"locales/en/active.json": &fstest.MapFile{
		Data: []byte(`[
  {
    "key": "test",
    "description": "Для тестов, не трогать",
    "translation": "Test of the %s"
  },
  {
    "key": "test plural",
    "description": "Для тестов, не трогать",
    "plural": {
      "other": "exactly %d spiders",
      "one": "spider",
      "=0": "no spiders",
      "=2": "just pair of spiders"
    }
  }
]`),
		Mode:    0555,
		ModTime: time.Now(),
		Sys:     1,
	},
}
