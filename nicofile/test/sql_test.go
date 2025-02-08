package test

import (
	config2 "main/config"
	"main/model"
	"testing"
)

func TestSqlQuery(t *testing.T) {
	DB := config2.InitDB()
	DB.FirstOrCreate(&model.File{FileName: "118268382_p0_0"})
	num := int64(0)
	DB.Model(&model.File{}).Where("file_name = ? && is_chunk = 1", "118268382_p0_0").Count(&num)
	if num == 0 {
		t.Errorf("Query failed num: %d\n", num)
	}
}
