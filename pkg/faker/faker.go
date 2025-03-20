package faker

import (
	"fmt"
	"math/rand"
	"time"
)

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Boolean returns a random boolean value
func Boolean() bool {
	return rand.Intn(2) == 1
}

// City returns a random city name
func City() string {
	cities := []string{
		"New York", "Los Angeles", "Chicago", "Houston", "Phoenix",
		"Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose",
		"Austin", "Jacksonville", "Fort Worth", "Columbus", "San Francisco",
		"Charlotte", "Indianapolis", "Seattle", "Denver", "Washington",
		"Boston", "El Paso", "Nashville", "Detroit", "Portland",
		"Memphis", "Oklahoma City", "Las Vegas", "Louisville", "Baltimore",
	}
	return cities[rand.Intn(len(cities))]
}

// State returns a random US state
func State() string {
	states := []string{
		"Alabama", "Alaska", "Arizona", "Arkansas", "California",
		"Colorado", "Connecticut", "Delaware", "Florida", "Georgia",
		"Hawaii", "Idaho", "Illinois", "Indiana", "Iowa",
		"Kansas", "Kentucky", "Louisiana", "Maine", "Maryland",
		"Massachusetts", "Michigan", "Minnesota", "Mississippi", "Missouri",
		"Montana", "Nebraska", "Nevada", "New Hampshire", "New Jersey",
		"New Mexico", "New York", "North Carolina", "North Dakota", "Ohio",
		"Oklahoma", "Oregon", "Pennsylvania", "Rhode Island", "South Carolina",
		"South Dakota", "Tennessee", "Texas", "Utah", "Vermont",
		"Virginia", "Washington", "West Virginia", "Wisconsin", "Wyoming",
	}
	return states[rand.Intn(len(states))]
}

// Zip returns a random ZIP code
func Zip() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

// CountryAbbr returns a random country abbreviation
func CountryAbbr() string {
	countries := []string{
		"US", "CA", "MX", "UK", "FR", "DE", "IT", "ES", "JP", "CN",
		"AU", "NZ", "BR", "AR", "CL", "RU", "IN", "ZA", "NG", "EG",
	}
	return countries[rand.Intn(len(countries))]
}

// ProductName returns a random product name
func ProductName() string {
	adjectives := []string{
		"Premium", "Deluxe", "Luxury", "Basic", "Essential", "Professional",
		"Advanced", "Smart", "Ultra", "Super", "Mega", "Compact", "Portable",
		"Wireless", "Digital", "Analog", "Classic", "Modern", "Vintage", "Retro",
	}

	nouns := []string{
		"Laptop", "Smartphone", "Tablet", "Headphones", "Speaker", "Camera",
		"Watch", "TV", "Monitor", "Keyboard", "Mouse", "Printer", "Scanner",
		"Router", "Charger", "Cable", "Adapter", "Case", "Stand", "Holder",
	}

	return adjectives[rand.Intn(len(adjectives))] + " " + nouns[rand.Intn(len(nouns))]
}

// ProductDescription returns a random product description
func ProductDescription() string {
	descriptions := []string{
		"This high-quality product is designed to meet all your needs.",
		"Experience the ultimate performance with this innovative product.",
		"A reliable solution for everyday use with exceptional durability.",
		"Combining style and functionality in a compact design.",
		"The perfect balance of quality, performance, and value.",
		"Engineered for maximum efficiency and user satisfaction.",
		"A versatile product suitable for various applications.",
		"Featuring cutting-edge technology for superior results.",
		"Designed with user comfort and convenience in mind.",
		"A must-have addition to your collection of premium products.",
	}
	return descriptions[rand.Intn(len(descriptions))]
}

// SKU generates a random SKU (Stock Keeping Unit)
func SKU() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"

	sku := make([]byte, 8)

	// First 3 characters are letters
	for i := 0; i < 3; i++ {
		sku[i] = letters[rand.Intn(len(letters))]
	}

	// Last 5 characters are numbers
	for i := 3; i < 8; i++ {
		sku[i] = numbers[rand.Intn(len(numbers))]
	}

	return string(sku)
}

// Price returns a random price between min and max
func Price(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// OrderStatus returns a random order status
func OrderStatus() string {
	statuses := []string{
		"Pending", "Processing", "Shipped", "Delivered", "Cancelled", "Refunded",
	}
	return statuses[rand.Intn(len(statuses))]
}

// PaymentMethod returns a random payment method
func PaymentMethod() string {
	methods := []string{
		"Credit Card", "Debit Card", "PayPal", "Apple Pay", "Google Pay",
		"Bank Transfer", "Cash on Delivery",
	}
	return methods[rand.Intn(len(methods))]
}

// ShippingMethod returns a random shipping method
func ShippingMethod() string {
	methods := []string{
		"Standard Shipping", "Express Shipping", "Next Day Delivery",
		"Two-Day Shipping", "Free Shipping", "International Shipping",
	}
	return methods[rand.Intn(len(methods))]
}

// TrackingNumber returns a random tracking number
func TrackingNumber() string {
	prefix := "TRK"
	number := ""
	for i := 0; i < 10; i++ {
		number += fmt.Sprintf("%d", rand.Intn(10))
	}
	return prefix + number
}

// ReviewTitle returns a random review title
func ReviewTitle() string {
	titles := []string{
		"Great product!", "Highly recommended", "Excellent value",
		"Not what I expected", "Could be better", "Amazing quality",
		"Disappointed", "Perfect for my needs", "Good but overpriced",
		"Exceeded expectations", "Just okay", "Very satisfied",
	}
	return titles[rand.Intn(len(titles))]
}

// ReviewContent returns a random review content
func ReviewContent() string {
	contents := []string{
		"I've been using this product for a few weeks now and I'm very satisfied with its performance and quality.",
		"This product exceeded my expectations in every way. The build quality is excellent and it works perfectly.",
		"While the product is good overall, I think it's a bit overpriced for what you get.",
		"I was disappointed with this purchase. The quality is not what I expected and it doesn't work as advertised.",
		"This is exactly what I was looking for. It's well-made, easy to use, and does the job perfectly.",
		"The product is okay, but there are better options available at this price point.",
		"I've tried many similar products, but this one is by far the best. Highly recommended!",
		"Great value for money. It's not perfect, but it gets the job done and is very affordable.",
		"I bought this as a gift and the recipient loved it. Great quality and nice packaging.",
		"The product arrived damaged, but customer service was excellent and sent a replacement right away.",
	}
	return contents[rand.Intn(len(contents))]
}

// CategoryName returns a random category name
func CategoryName() string {
	categories := []string{
		"Electronics", "Clothing", "Home & Kitchen", "Books", "Sports & Outdoors",
		"Beauty & Personal Care", "Toys & Games", "Automotive", "Health & Wellness",
		"Jewelry", "Office Supplies", "Pet Supplies", "Food & Grocery", "Garden & Outdoor",
		"Baby Products", "Tools & Home Improvement", "Musical Instruments", "Arts & Crafts",
	}
	return categories[rand.Intn(len(categories))]
}

// CategoryDescription returns a random category description
func CategoryDescription() string {
	descriptions := []string{
		"Find everything you need for your home and daily life.",
		"Quality products at affordable prices.",
		"The latest trends and innovations in this category.",
		"Essential items for every household.",
		"Premium selection of top-rated products.",
		"Discover new and exciting products in this category.",
		"Handpicked items to meet your specific needs.",
		"A wide range of products for every budget.",
		"Specialized products for enthusiasts and professionals.",
		"Everything you need in one convenient category.",
	}
	return descriptions[rand.Intn(len(descriptions))]
}

// Dimensions returns random product dimensions
func Dimensions() string {
	width := 1 + rand.Float64()*50
	height := 1 + rand.Float64()*50
	depth := 1 + rand.Float64()*50
	return fmt.Sprintf("%.1f x %.1f x %.1f cm", width, height, depth)
}

// Weight returns a random weight between min and max
func Weight(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// ImageURL returns a random product image URL
func ImageURL(productID int) string {
	baseURLs := []string{
		"https://example.com/images/products/",
		"https://store.example.org/product-images/",
		"https://cdn.example.net/shop/items/",
		"https://images.example.io/catalog/",
	}

	extensions := []string{".jpg", ".png", ".webp"}

	baseURL := baseURLs[rand.Intn(len(baseURLs))]
	extension := extensions[rand.Intn(len(extensions))]

	return fmt.Sprintf("%s%d-%d%s", baseURL, productID, rand.Intn(5)+1, extension)
}
