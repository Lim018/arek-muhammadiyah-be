package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"arek-muhammadiyah-be/app/model"
	"regexp"
	"strings"
	"time"
)

func CreatePagination(page, limit, total int64) model.Pagination {
	totalPages := (total + limit - 1) / limit
	return model.Pagination{
		CurrentPage: int(page),
		PerPage:     int(limit),
		TotalPages:  int(totalPages),
		TotalItems:  total,
	}
}

func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")
	
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	return slug
}

func GenerateUniqueSlug(title string) string {
	baseSlug := GenerateSlug(title)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", baseSlug, timestamp)
}

func GetStringValue(newVal *string, existing string) string {
	if newVal != nil {
		return *newVal
	}
	return existing
}

func GetStringPointer(newVal *string, existing *string) *string {
	if newVal != nil {
		return newVal
	}
	return existing
}

func GetUintPointer(newVal *uint, existing *uint) *uint {
	if newVal != nil {
		return newVal
	}
	return existing
}

func GetBoolValue(newVal *bool, defaultVal bool) bool {
	if newVal != nil {
		return *newVal
	}
	return defaultVal
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}