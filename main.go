package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"time"
)

type HarvestClient struct {
	AccessToken string
	AccountID   string
	BaseURL     string
	UserAgent   string
}

type TimeReport struct {
	Results      []TeamMember `json:"results"`
	TotalPages   int          `json:"total_pages"`
	TotalEntries int          `json:"total_entries"`
	PerPage      int          `json:"per_page"`
}

type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TeamMember struct {
	UserName      string  `json:"user_name"`
	TotalHours    float64 `json:"total_hours"`
	BillableHours float64 `json:"billable_hours"`
}

func NewHarvestClient(accessToken, accountID string) *HarvestClient {
	return &HarvestClient{
		AccessToken: accessToken,
		AccountID:   accountID,
		BaseURL:     "https://api.harvestapp.com/v2",
		UserAgent:   "HarvestCLI (self@travisfantina.com)",
	}
}

func (c *HarvestClient) GetTimeReport(from, to string) (*TeamMember, error) {
	url := fmt.Sprintf("%s/reports/time/team?from=%s&to=%s", c.BaseURL, from, to)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Harvest-Account-Id", c.AccountID)
	req.Header.Set("User-Agent", c.UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Read and log raw response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var report TimeReport
	if err := json.Unmarshal(body, &report); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	idx := slices.IndexFunc(report.Results, func(r TeamMember) bool { return r.UserName == "Travis Fantina" })

	if idx < 0 {
		return nil, nil
	}

	return &report.Results[idx], nil
}

func getCurrentMonthDateRange() (string, string) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	return startOfMonth.Format("20060102"), now.Format("20060102")
}

func getCurrentWeek() (string, string) {
	now := time.Now()
	startOfWeek := now

	for startOfWeek.Weekday() != 1 {
		startOfWeek = startOfWeek.AddDate(0, 0, -1)
	}

	return startOfWeek.Format("20060102"), now.Format("20060102")
}

func getCurrentDay() string {
	now := time.Now()
	return now.Format("20060102")
}

func main() {
	godotenv.Load()
	accessToken := os.Getenv("ACCESS_TOKEN")
	accountID := os.Getenv("ACCOUNT_ID")
	accessToken2 := os.Getenv("ACCESS_TOKEN2")
	accountID2 := os.Getenv("ACCOUNT_ID2")

	var credentials = []string{accessToken, accountID, accessToken2, accountID2}

	for _, c := range credentials {
		if c == "" {
			log.Fatalf("Missing credential %s", c)
		}
	}

	client1 := NewHarvestClient(accessToken, accountID)
	client2 := NewHarvestClient(accessToken2, accountID2)

	monthFrom, monthTo := getCurrentMonthDateRange()
	weekFrom, weekTo := getCurrentWeek()
	day := getCurrentDay()

	hours1, err1 := client1.GetTimeReport(monthFrom, monthTo)
	hours2, err2 := client2.GetTimeReport(monthFrom, monthTo)
	weekHours1, err3 := client1.GetTimeReport(weekFrom, weekTo)
	weekHours2, err4 := client2.GetTimeReport(weekFrom, weekTo)
	dayHours1, err5 := client1.GetTimeReport(day, day)
	dayHours2, err6 := client2.GetTimeReport(day, day)

	fmt.Printf("%s", weekFrom)
	fmt.Print(weekTo)

	var errors = []error{err1, err2, err3, err4, err5, err6}

	if !slices.Contains(errors, nil) {
		fmt.Print(errors)
		log.Fatalf("Error fetching expense report")
	}

	var results = []*TeamMember{hours1, hours2}
	var weekResults = []*TeamMember{weekHours1, weekHours2}
	var dayResults = []*TeamMember{dayHours1, dayHours2}

	var grandTotal float64
	var grandTotalWeek float64
	var grandTotalDay float64

	for _, result := range results {
		if result.TotalHours > 0 {
			grandTotal += result.TotalHours
		}
	}

	for _, result := range weekResults {
		if result != nil && result.TotalHours > 0 {
			grandTotalWeek += result.TotalHours
		}
	}

	for _, result := range dayResults {
		if result != nil && result.TotalHours > 0 {
			grandTotalDay += result.TotalHours
		}
	}

	fmt.Printf("Hours Today: %.2f\n", grandTotalDay)
	fmt.Printf("Hours Week: %.2f\n", grandTotalWeek)
	fmt.Printf("----------\n")
	fmt.Printf("Objective %.2f\n", hours1.TotalHours)
	fmt.Printf("Built %.2f\n", hours2.TotalHours)
	fmt.Printf("\nGrand Total: %.2f\n", grandTotal)
	fmt.Printf("Remaining To Goal: %.2f\n", 140-grandTotal)
}
