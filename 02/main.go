package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID           string  `json:"id"`
	Name         string  `json:"productName"`
	Price        float64 `json:"price"`
	Rating       float64 `json:"rating"`
	Discount     int     `json:"discount"`
	Availability string  `json:"availability"`
	Company      string  `json:"company"`
	Category     string  `json:"category"`
}

type ECommerceClient struct {
	baseURL string
}

func NewECommerceClient(baseURL string) *ECommerceClient {
	return &ECommerceClient{
		baseURL: baseURL,
	}
}

func (c *ECommerceClient) GetTopProducts(category string, minPrice, maxPrice float64, top int) ([]Product, error) {
	url := fmt.Sprintf("%s/categories/%s/products?top=%d&minPrice=%.2f&maxPrice=%.2f", c.baseURL, category, top, minPrice, maxPrice)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error fetching top products")
	}

	var products []Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c *ECommerceClient) GetProductDetails(category, productId string) (Product, error) {
	url := fmt.Sprintf("%s/categories/%s/products/%s", c.baseURL, category, productId)

	resp, err := http.Get(url)
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Product{}, errors.New("Error fetching product details")
	}

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func TopProductsHandler(eCommerceClient ECommerceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Param("category")

		top, _ := strconv.Atoi(c.Query("top"))
		minPrice, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
		maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)

		var products []Product
		var err error
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			products, err = eCommerceClient.GetTopProducts(category, minPrice, maxPrice, top)
		}()
		wg.Wait()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func ProductDetailsHandler(eCommerceClient ECommerceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Param("category")
		productId := c.Param("productid")

		var product Product
		var err error
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			product, err = eCommerceClient.GetProductDetails(category, productId)
		}()
		wg.Wait()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func main() {
	router := gin.Default()
	router.GET("/categories/:category/products", TopProductsHandler(ECommerceClient{}))
	router.GET("/categories/:category/products/:productid", ProductDetailsHandler(ECommerceClient{}))
	router.Run(":8080")
}
