package model

import "testing"

func TestSaveData(t *testing.T) {
	// Vytvoření mock dat pro uložení
	data := map[string]interface{}{
		"temperature": 22.5,
		"humidity":    60.1,
	}

	// Vytvoření mock InfluxDB klienta
	i := &InfluxDB{
		Host:     "localhost",
		Port:     8086,
		Username: "",
		Password: "",
		Database: "test",
	}

	// Volání funkce SaveData
	if err := i.SaveData("sensors", data); err != nil {
		t.Errorf("Error saving data: %s", err)
	}
}
