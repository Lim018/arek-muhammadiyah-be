package helper

import (
	"encoding/csv"
	"io"
	"github.com/Lim018/arek-muhammadiyah-be/app/model"
	"strconv"
	"strings"
)

type CSVUserData struct {
	ID        string `csv:"id"`
	Name      string `csv:"name"`
	NIK       string `csv:"nik"`
	Address   string `csv:"address"`
	VillageID string `csv:"village_id"`
}

func ParseUsersFromCSV(reader io.Reader) ([]model.CreateUserRequest, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var users []model.CreateUserRequest
	
	// Skip header row
	for i, record := range records {
		if i == 0 {
			continue
		}
		
		if len(record) < 5 {
			continue
		}

		user := model.CreateUserRequest{
			ID:       strings.TrimSpace(record[0]),
			Name:     strings.TrimSpace(record[1]),
			Password: GenerateRandomString(8), // Generate random password
			NIK:      &record[2],
			Address:  &record[3],
		}

		// Parse village ID
		if villageIDStr := strings.TrimSpace(record[4]); villageIDStr != "" {
			if villageID, err := strconv.ParseUint(villageIDStr, 10, 32); err == nil {
				villageIDUint := uint(villageID)
				user.VillageID = &villageIDUint
			}
		}

		users = append(users, user)
	}

	return users, nil
}