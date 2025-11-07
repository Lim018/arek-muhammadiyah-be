package helper

import (
	"encoding/csv"
	"io"
	"arek-muhammadiyah-be/app/model"
	"strings"
	"time"
)

type CSVUserData struct {
	ID        string `csv:"id"`
	Name      string `csv:"name"`
	BirthDate string `csv:"birth_date"`
	Telp      string `csv:"telp"`
	Gender    string `csv:"gender"`
	Job       string `csv:"job"`
	NIK       string `csv:"nik"`
	Address   string `csv:"address"`
	VillageID string `csv:"village_id"`
	IsMobile  string `csv:"is_mobile"`
}

func ParseUsersFromCSV(reader io.Reader) ([]model.CreateUserRequest, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var users []model.CreateUserRequest
	
	for i, record := range records {
		if i == 0 {
			continue
		}
		
		if len(record) < 3 {
			continue
		}

		user := model.CreateUserRequest{
			ID:       strings.TrimSpace(record[0]),
			Name:     strings.TrimSpace(record[1]),
			Password: GenerateRandomString(8),
		}

		if len(record) > 2 && strings.TrimSpace(record[2]) != "" {
			nik := strings.TrimSpace(record[2])
			user.NIK = &nik
		}

		// Address (index 3)
		if len(record) > 3 && strings.TrimSpace(record[3]) != "" {
			address := strings.TrimSpace(record[3])
			user.Address = &address
		}

		// VillageID (index 4) - ID dari JSON wilayah (string 10 digit)
		if len(record) > 4 && strings.TrimSpace(record[4]) != "" {
			villageID := strings.TrimSpace(record[4])
			// Validasi format (harus 10 digit)
			if len(villageID) == 10 {
				user.VillageID = &villageID
			}
		}

		// Telp (index 5)
		if len(record) > 5 && strings.TrimSpace(record[5]) != "" {
			telp := strings.TrimSpace(record[5])
			user.Telp = &telp
		}

		// Gender (index 6)
		if len(record) > 6 && strings.TrimSpace(record[6]) != "" {
			gender := strings.TrimSpace(strings.ToLower(record[6]))
			user.Gender = &gender
		}

		// Job (index 7)
		if len(record) > 7 && strings.TrimSpace(record[7]) != "" {
			job := strings.TrimSpace(record[7])
			user.Job = &job
		}

		// BirthDate (index 8) - format: YYYY-MM-DD
		if len(record) > 8 && strings.TrimSpace(record[8]) != "" {
			if birthDate, err := time.Parse("2006-01-02", strings.TrimSpace(record[8])); err == nil {
				user.BirthDate = &birthDate
			}
		}

		// IsMobile (index 9)
		if len(record) > 9 && strings.TrimSpace(record[9]) != "" {
			isMobile := strings.TrimSpace(strings.ToLower(record[9])) == "true" || record[9] == "1"
			user.IsMobile = &isMobile
		}

		users = append(users, user)
	}

	return users, nil
}