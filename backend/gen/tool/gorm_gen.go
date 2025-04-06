package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	LoadConfig()

	g := gen.NewGenerator(gen.Config{
		OutPath:      "../gorm",
		ModelPkgPath: "gorm",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	db, err := gorm.Open(postgres.Open(value.Config.PostgresSQL), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	g.UseDB(db)
	g.GenerateAllTable()
	// 生成代码
	g.Execute()
}
