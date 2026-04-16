package handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"url-shortener/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type URLHandler struct {
	DB      *gorm.DB
	BaseURL string
}

type shortenRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

type shortenResponse struct {
	Code        string `json:"code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type statsResponse struct {
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int       `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
	ShortURL    string    `json:"short_url"`
}

func NewURLHandler(db *gorm.DB, baseURL string) *URLHandler {
	return &URLHandler{DB: db, BaseURL: strings.TrimSuffix(baseURL, "/")}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[API] shorten failed: invalid payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	log.Printf("[API] shorten requested for URL=%s", req.OriginalURL)

	if _, err := url.ParseRequestURI(req.OriginalURL); err != nil {
		log.Printf("[API] shorten failed: invalid URL=%s", req.OriginalURL)
		c.JSON(http.StatusBadRequest, gin.H{"error": "original_url must be a valid URL"})
		return
	}

	code, err := h.generateUniqueCode(6)
	if err != nil {
		log.Printf("[API] shorten failed: could not generate code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate short code"})
		return
	}

	record := models.URL{OriginalURL: req.OriginalURL, ShortCode: code}
	if err := h.DB.Create(&record).Error; err != nil {
		log.Printf("[API] shorten failed: database insert error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save URL"})
		return
	}

	log.Printf("[API] shorten success: code=%s original=%s", record.ShortCode, record.OriginalURL)

	c.JSON(http.StatusCreated, shortenResponse{
		Code:        record.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", h.BaseURL, record.ShortCode),
		OriginalURL: record.OriginalURL,
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	log.Printf("[API] redirect requested for code=%s", code)
	var record models.URL

	if err := h.DB.Where("short_code = ?", code).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[API] redirect failed: code not found=%s", code)
			c.JSON(http.StatusNotFound, gin.H{"error": "short code not found"})
			return
		}
		log.Printf("[API] redirect failed: database lookup error for code=%s err=%v", code, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch URL"})
		return
	}

	if err := h.DB.Model(&record).UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error; err != nil {
		log.Printf("[API] redirect failed: click count update error for code=%s err=%v", code, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update click count"})
		return
	}

	log.Printf("[API] redirect success: code=%s -> %s", code, record.OriginalURL)

	c.Redirect(http.StatusFound, record.OriginalURL)
}

func (h *URLHandler) GetStats(c *gin.Context) {
	code := c.Param("code")
	log.Printf("[API] stats requested for code=%s", code)
	var record models.URL

	if err := h.DB.Where("short_code = ?", code).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[API] stats failed: code not found=%s", code)
			c.JSON(http.StatusNotFound, gin.H{"error": "short code not found"})
			return
		}
		log.Printf("[API] stats failed: database lookup error for code=%s err=%v", code, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch URL"})
		return
	}

	log.Printf("[API] stats success: code=%s clicks=%d", code, record.ClickCount)

	c.JSON(http.StatusOK, statsResponse{
		Code:        record.ShortCode,
		OriginalURL: record.OriginalURL,
		ClickCount:  record.ClickCount,
		CreatedAt:   record.CreatedAt,
		ShortURL:    fmt.Sprintf("%s/%s", h.BaseURL, record.ShortCode),
	})
}

func (h *URLHandler) generateUniqueCode(length int) (string, error) {
	for range 10 {
		code, err := generateCode(length)
		if err != nil {
			return "", err
		}

		var count int64
		if err := h.DB.Model(&models.URL{}).Where("short_code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique code")
}

func generateCode(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}

	return string(result), nil
}
