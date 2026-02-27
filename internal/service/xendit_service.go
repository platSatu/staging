// service/xendit_service.go

package service

import (
	"backend_go/internal/model"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type XenditService struct {
	SecretKey        string
	CallbackToken    string
	TransaksiService *TransaksiService
}

func NewXenditService(transaksiService *TransaksiService) *XenditService {
	secretKey := os.Getenv("XENDIT_SECRET_KEY")
	callbackToken := os.Getenv("XENDIT_CALLBACK_TOKEN")

	// 🔥 DEBUG
	fmt.Println("=== XENDIT DEBUG ===")
	fmt.Println("Secret Key from env:", secretKey)
	fmt.Println("Callback Token from env:", callbackToken)
	// ======================

	if secretKey == "" {
		fmt.Println("WARNING: XENDIT_SECRET_KEY is empty!")
	}

	return &XenditService{
		SecretKey:        secretKey,
		CallbackToken:    callbackToken,
		TransaksiService: transaksiService,
	}
}

// CreateInvoice - Buat invoice Xendit menggunakan HTTP request langsung
func (s *XenditService) CreateInvoice(transaksiID string, amount float64, email, description string) (map[string]interface{}, error) {
	externalID := fmt.Sprintf("transaksi-%s-%d", transaksiID, time.Now().Unix())

	// 🔥 DEBUG
	fmt.Println("Creating invoice with key:", s.SecretKey)
	// ==========

	// Payload untuk Xendit
	payload := map[string]interface{}{
		"external_id": externalID,
		"amount":      amount,
		"payer_email": email,
		"description": description,
	}

	// Convert ke JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	// HTTP Request ke Xendit
	url := "https://api.xendit.co/v2/invoices"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 🔥 PAKE BASIC AUTH - Sesuai требование Xendit
	auth := s.SecretKey + ":"
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encoded)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Baca response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Printf("Xendit Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Xendit Response Body: %s\n", string(body))

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("Xendit API error (status %d): %s", resp.StatusCode, string(body))
	}

	var xenditResp map[string]interface{}
	if err := json.Unmarshal(body, &xenditResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Update transaksi dengan reference_id dan status
	refID, _ := xenditResp["invoice_url"].(string)
	updateData := &model.Transaksi{
		ID:            transaksiID,
		ReferenceID:   &refID,
		StatusGateway: stringToPtr("pending"),
	}

	if err := s.TransaksiService.UpdateTransaksi(updateData); err != nil {
		return nil, fmt.Errorf("failed to update transaksi: %v", err)
	}

	return map[string]interface{}{
		"invoice_id":  xenditResp["id"],
		"external_id": xenditResp["external_id"],
		"invoice_url": xenditResp["invoice_url"],
		"amount":      xenditResp["amount"],
		"status":      xenditResp["status"],
		"expiry_date": xenditResp["expiry_date"],
	}, nil
}

// HandleCallback - Handle callback dari Xendit
func (s *XenditService) HandleCallback(c *gin.Context) {
	callbackToken := c.GetHeader("x-callback-token")
	if callbackToken != s.CallbackToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid callback token",
		})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to read request body"})
		return
	}

	var callbackData map[string]interface{}
	if err := json.Unmarshal(body, &callbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to parse callback data"})
		return
	}

	externalID, _ := callbackData["external_id"].(string)
	status, _ := callbackData["status"].(string)

	transaksiID := strings.TrimPrefix(externalID, "transaksi-")
	parts := strings.Split(transaksiID, "-")
	if len(parts) > 0 {
		transaksiID = parts[0]
	}

	var newStatus string
	switch status {
	case "PAID":
		newStatus = "success"
	case "EXPIRED":
		newStatus = "failed"
	default:
		newStatus = "pending"
	}

	updateData := &model.Transaksi{
		ID:            transaksiID,
		StatusGateway: &newStatus,
	}

	if err := s.TransaksiService.UpdateTransaksi(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Callback processed successfully"})
}

func stringToPtr(s string) *string {
	return &s
}
