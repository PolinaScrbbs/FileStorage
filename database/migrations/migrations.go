package migrations

import (
	"gorm.io/gorm"
	"log"
	"reflect"
)

func Migrate(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		err := db.AutoMigrate(model)

		if err != nil {
			log.Fatalf("Error applying migration for model %v: %v", model, err)
		}
		modelName := reflect.TypeOf(model).Elem().Name()
		log.Printf("Migration applied for model: %s", modelName)
	}

	log.Printf("All migrations applied")
}
