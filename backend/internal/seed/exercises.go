package seed

import (
	"encoding/json"

	"github.com/akaitigo/review-gym/internal/model"
)

func exercise01SQLInjection() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "User search endpoint with string concatenation",
			Description: "A PR adding a user search API endpoint. The implementation builds SQL queries by concatenating user input directly into the query string.",
			Difficulty:  model.DifficultyBeginner,
			Category:    model.CategorySecurity,
			CategoryTags: []model.Category{
				model.CategorySecurity,
				model.CategoryDesign,
			},
			Language: "Go",
			DiffContent: `--- a/internal/handler/user.go
+++ b/internal/handler/user.go
@@ -0,0 +1,35 @@
+package handler
+
+import (
+	"database/sql"
+	"fmt"
+	"net/http"
+	"encoding/json"
+)
+
+type UserHandler struct {
+	db *sql.DB
+}
+
+func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
+	query := r.URL.Query().Get("q")
+	if query == "" {
+		http.Error(w, "missing query parameter", http.StatusBadRequest)
+		return
+	}
+
+	sqlQuery := fmt.Sprintf("SELECT id, name, email FROM users WHERE name LIKE '%%%s%%'", query)
+	rows, err := h.db.Query(sqlQuery)
+	if err != nil {
+		http.Error(w, "internal error", http.StatusInternalServerError)
+		return
+	}
+	defer rows.Close()
+
+	var users []map[string]string
+	for rows.Next() {
+		var id, name, email string
+		rows.Scan(&id, &name, &email)
+		users = append(users, map[string]string{"id": id, "name": name, "email": email})
+	}
+
+	json.NewEncoder(w).Encode(users)
+}`,
			FilePaths:   []string{"internal/handler/user.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/handler/user.go",
				LineNumber:  21,
				Content:     "SQL injection vulnerability: user input is directly concatenated into the SQL query string.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "The query parameter 'q' is interpolated into the SQL string using fmt.Sprintf without any sanitization. An attacker can inject arbitrary SQL. Use parameterized queries (db.Query with $1 placeholders) instead.",
			},
			{
				FilePath:    "internal/handler/user.go",
				LineNumber:  32,
				Content:     "rows.Scan error is silently ignored, which may lead to incomplete or corrupt data.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "The return value of rows.Scan() is not checked. If scanning fails, the zero-value struct is appended to results, causing silent data corruption.",
			},
			{
				FilePath:    "internal/handler/user.go",
				LineNumber:  35,
				Content:     "Missing rows.Err() check after iteration loop.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "After iterating over sql.Rows, you must check rows.Err() to detect errors that occurred during iteration. Without this, partial results may be returned without any indication of failure.",
			},
			{
				FilePath:    "internal/handler/user.go",
				LineNumber:  29,
				Content:     "Using map[string]string instead of a typed struct reduces type safety and readability.",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityMinor,
				Explanation: "Define a User struct with proper fields instead of using a map. This provides compile-time type checking and makes the API contract explicit.",
			},
		},
	}
}

func exercise02UnboundedQuery() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Product listing API without pagination",
			Description: "A PR implementing a product listing endpoint that fetches all products from the database without any pagination or limit.",
			Difficulty:  model.DifficultyBeginner,
			Category:    model.CategoryPerformance,
			CategoryTags: []model.Category{
				model.CategoryPerformance,
				model.CategoryDesign,
			},
			Language: "Go",
			DiffContent: `--- a/internal/handler/product.go
+++ b/internal/handler/product.go
@@ -0,0 +1,42 @@
+package handler
+
+import (
+	"database/sql"
+	"encoding/json"
+	"net/http"
+)
+
+type Product struct {
+	ID          string  ` + "`json:\"id\"`" + `
+	Name        string  ` + "`json:\"name\"`" + `
+	Description string  ` + "`json:\"description\"`" + `
+	Price       float64 ` + "`json:\"price\"`" + `
+	ImageURL    string  ` + "`json:\"image_url\"`" + `
+}
+
+type ProductHandler struct {
+	db *sql.DB
+}
+
+func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
+	rows, err := h.db.Query("SELECT id, name, description, price, image_url FROM products")
+	if err != nil {
+		http.Error(w, "failed to query products", http.StatusInternalServerError)
+		return
+	}
+	defer rows.Close()
+
+	products := make([]Product, 0)
+	for rows.Next() {
+		var p Product
+		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL); err != nil {
+			http.Error(w, "failed to scan product", http.StatusInternalServerError)
+			return
+		}
+		products = append(products, p)
+	}
+
+	w.Header().Set("Content-Type", "application/json")
+	json.NewEncoder(w).Encode(products)
+}`,
			FilePaths:   []string{"internal/handler/product.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/handler/product.go",
				LineNumber:  22,
				Content:     "Unbounded query: SELECT without LIMIT can return millions of rows and exhaust memory.",
				Category:    model.CategoryPerformance,
				Severity:    model.SeverityCritical,
				Explanation: "The query fetches ALL products from the database. In production with thousands or millions of products, this will consume excessive memory and cause very slow response times. Add LIMIT/OFFSET pagination or cursor-based pagination.",
			},
			{
				FilePath:    "internal/handler/product.go",
				LineNumber:  22,
				Content:     "No filtering or search capability. Users cannot narrow down results.",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityMajor,
				Explanation: "A listing endpoint should support filtering (by category, price range, etc.) and sorting. Without these, the API is impractical for any real use case.",
			},
			{
				FilePath:    "internal/handler/product.go",
				LineNumber:  39,
				Content:     "Missing rows.Err() check after the iteration loop.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "Always check rows.Err() after the loop to catch errors that occurred during iteration (e.g., connection dropped mid-stream).",
			},
		},
	}
}

func exercise03GodFunction() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Order processing function handling everything",
			Description: "A PR adding an order processing function that validates input, checks inventory, calculates pricing, processes payment, sends notifications, and updates analytics - all in a single function.",
			Difficulty:  model.DifficultyIntermediate,
			Category:    model.CategoryDesign,
			CategoryTags: []model.Category{
				model.CategoryDesign,
				model.CategoryReadability,
			},
			Language: "Go",
			DiffContent: `--- a/internal/service/order.go
+++ b/internal/service/order.go
@@ -0,0 +1,85 @@
+package service
+
+import (
+	"fmt"
+	"net/smtp"
+	"time"
+)
+
+func ProcessOrder(userID string, items []Item, paymentInfo PaymentInfo) error {
+	// Validate input
+	if userID == "" {
+		return fmt.Errorf("user ID is required")
+	}
+	if len(items) == 0 {
+		return fmt.Errorf("at least one item is required")
+	}
+	for _, item := range items {
+		if item.Quantity <= 0 {
+			return fmt.Errorf("invalid quantity for item %s", item.ID)
+		}
+	}
+
+	// Check inventory
+	for _, item := range items {
+		stock, _ := getInventory(item.ID)
+		if stock < item.Quantity {
+			return fmt.Errorf("insufficient stock for %s", item.ID)
+		}
+	}
+
+	// Calculate total
+	var total float64
+	for _, item := range items {
+		price, _ := getPrice(item.ID)
+		total += price * float64(item.Quantity)
+	}
+	if total > 10000 {
+		total = total * 0.9 // 10% discount for large orders
+	}
+	tax := total * 0.1
+	total = total + tax
+
+	// Process payment
+	err := chargeCard(paymentInfo.CardNumber, paymentInfo.CVV, total)
+	if err != nil {
+		return fmt.Errorf("payment failed: %w", err)
+	}
+
+	// Update inventory
+	for _, item := range items {
+		decrementInventory(item.ID, item.Quantity)
+	}
+
+	// Create order record
+	orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
+	saveOrder(orderID, userID, items, total)
+
+	// Send confirmation email
+	body := fmt.Sprintf("Your order %s has been placed. Total: $%.2f", orderID, total)
+	smtp.SendMail("smtp.example.com:587", nil, "noreply@example.com",
+		[]string{getUserEmail(userID)}, []byte(body))
+
+	// Update analytics
+	incrementSalesCount(len(items))
+	updateRevenue(total)
+	trackUserPurchase(userID, orderID)
+
+	return nil
+}`,
			FilePaths:   []string{"internal/service/order.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/service/order.go",
				LineNumber:  9,
				Content:     "God function: ProcessOrder handles validation, inventory, pricing, payment, notifications, and analytics. Split into smaller, single-responsibility functions.",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityCritical,
				Explanation: "This function violates the Single Responsibility Principle. Each concern (validation, inventory management, pricing, payment processing, notification, analytics) should be a separate function or service. This makes the code hard to test, maintain, and extend.",
			},
			{
				FilePath:    "internal/service/order.go",
				LineNumber:  25,
				Content:     "Error from getInventory() is silently ignored with blank identifier.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "If the inventory check fails (e.g., database down), the error is discarded and execution continues with a zero value for stock. This can lead to overselling.",
			},
			{
				FilePath:    "internal/service/order.go",
				LineNumber:  38,
				Content:     "Business rule (10% discount for orders > $10,000) is hardcoded as magic numbers.",
				Category:    model.CategoryReadability,
				Severity:    model.SeverityMinor,
				Explanation: "The discount threshold (10000) and rate (0.9) should be named constants or configurable values. Magic numbers make the code harder to understand and maintain.",
			},
			{
				FilePath:    "internal/service/order.go",
				LineNumber:  50,
				Content:     "Inventory decrement happens after payment but has no error handling. If it fails, inventory becomes inconsistent.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityCritical,
				Explanation: "This is not transactional. If decrementInventory fails after payment succeeds, the user is charged but inventory is not updated. Use database transactions or implement a saga pattern for distributed operations.",
			},
		},
	}
}

func exercise04ErrorSwallowing() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "File upload handler with ignored errors",
			Description: "A PR implementing a file upload endpoint where multiple error returns are silently discarded, leading to potential data loss and silent failures.",
			Difficulty:  model.DifficultyBeginner,
			Category:    model.CategoryErrorHandling,
			CategoryTags: []model.Category{
				model.CategoryErrorHandling,
				model.CategorySecurity,
			},
			Language: "Go",
			DiffContent: `--- a/internal/handler/upload.go
+++ b/internal/handler/upload.go
@@ -0,0 +1,45 @@
+package handler
+
+import (
+	"io"
+	"net/http"
+	"os"
+	"path/filepath"
+)
+
+func UploadFile(w http.ResponseWriter, r *http.Request) {
+	r.ParseMultipartForm(32 << 20)
+
+	file, header, _ := r.FormFile("document")
+	defer file.Close()
+
+	uploadDir := "/var/uploads/" + r.FormValue("user_id")
+	os.MkdirAll(uploadDir, 0777)
+
+	destPath := filepath.Join(uploadDir, header.Filename)
+	dst, _ := os.Create(destPath)
+	defer dst.Close()
+
+	io.Copy(dst, file)
+
+	w.WriteHeader(http.StatusOK)
+	w.Write([]byte("uploaded"))
+}`,
			FilePaths:   []string{"internal/handler/upload.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/handler/upload.go",
				LineNumber:  13,
				Content:     "Error from r.FormFile is silently ignored. If the file is missing, 'file' will be nil and the next line will panic.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityCritical,
				Explanation: "If r.FormFile returns an error (missing file, wrong field name), the file variable is nil. Calling file.Close() on nil will panic. Always check errors from FormFile.",
			},
			{
				FilePath:    "internal/handler/upload.go",
				LineNumber:  16,
				Content:     "Path traversal vulnerability: user_id is used directly in the file path without sanitization.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "An attacker can set user_id to '../../etc' to write files outside the intended directory. Sanitize user_id to prevent directory traversal attacks.",
			},
			{
				FilePath:    "internal/handler/upload.go",
				LineNumber:  17,
				Content:     "Directory created with 0777 permissions allows any user on the system to read/write/execute.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityMajor,
				Explanation: "Use restrictive permissions like 0750 or 0700. World-writable directories are a security risk, especially for uploaded content.",
			},
			{
				FilePath:    "internal/handler/upload.go",
				LineNumber:  19,
				Content:     "Original filename is used directly, enabling path traversal via filename (e.g., '../../../etc/passwd').",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "The uploaded filename should be sanitized or replaced with a generated name. filepath.Base() at minimum, but a UUID-based name is safer.",
			},
		},
	}
}

func exercise05HardcodedSecret() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Authentication middleware with hardcoded JWT secret",
			Description: "A PR implementing JWT authentication middleware where the signing secret is hardcoded in the source code and token validation has several issues.",
			Difficulty:  model.DifficultyBeginner,
			Category:    model.CategorySecurity,
			CategoryTags: []model.Category{
				model.CategorySecurity,
				model.CategoryErrorHandling,
			},
			Language: "Go",
			DiffContent: `--- a/internal/middleware/auth.go
+++ b/internal/middleware/auth.go
@@ -0,0 +1,48 @@
+package middleware
+
+import (
+	"context"
+	"net/http"
+	"strings"
+
+	"github.com/golang-jwt/jwt/v5"
+)
+
+const jwtSecret = "super-secret-key-2024"
+
+type contextKey string
+const userIDKey contextKey = "userID"
+
+func AuthMiddleware(next http.Handler) http.Handler {
+	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		authHeader := r.Header.Get("Authorization")
+		if authHeader == "" {
+			http.Error(w, "missing authorization header", http.StatusUnauthorized)
+			return
+		}
+
+		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
+
+		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
+			return []byte(jwtSecret), nil
+		})
+
+		claims, _ := token.Claims.(jwt.MapClaims)
+		userID, _ := claims["sub"].(string)
+
+		ctx := context.WithValue(r.Context(), userIDKey, userID)
+		next.ServeHTTP(w, r.WithContext(ctx))
+	})
+}`,
			FilePaths:   []string{"internal/middleware/auth.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/middleware/auth.go",
				LineNumber:  11,
				Content:     "JWT secret is hardcoded in source code. This will be committed to version control.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "Secrets must never be hardcoded. Use environment variables or a secrets manager (e.g., AWS Secrets Manager, HashiCorp Vault). If this is committed, the secret is compromised and must be rotated.",
			},
			{
				FilePath:    "internal/middleware/auth.go",
				LineNumber:  26,
				Content:     "JWT parsing error is silently ignored. Invalid/expired tokens will pass through.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "The error from jwt.Parse is discarded with _. If the token is invalid, expired, or tampered with, the code continues to extract claims from a nil or invalid token. This bypasses authentication entirely.",
			},
			{
				FilePath:    "internal/middleware/auth.go",
				LineNumber:  26,
				Content:     "No algorithm validation in jwt.Parse. Vulnerable to algorithm confusion attacks.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityMajor,
				Explanation: "The keyFunc does not validate token.Method. An attacker could use the 'none' algorithm or switch from RS256 to HS256 to forge tokens. Verify that token.Method matches the expected algorithm.",
			},
		},
	}
}

func exercise06NPlus1Query() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Blog posts API with N+1 query problem",
			Description: "A PR implementing a blog posts listing endpoint that loads author information in a loop, resulting in N+1 database queries.",
			Difficulty:  model.DifficultyIntermediate,
			Category:    model.CategoryPerformance,
			CategoryTags: []model.Category{
				model.CategoryPerformance,
				model.CategoryDesign,
			},
			Language: "Go",
			DiffContent: `--- a/internal/handler/blog.go
+++ b/internal/handler/blog.go
@@ -0,0 +1,55 @@
+package handler
+
+import (
+	"database/sql"
+	"encoding/json"
+	"net/http"
+)
+
+type BlogPost struct {
+	ID         string ` + "`json:\"id\"`" + `
+	Title      string ` + "`json:\"title\"`" + `
+	Content    string ` + "`json:\"content\"`" + `
+	AuthorID   string ` + "`json:\"author_id\"`" + `
+	AuthorName string ` + "`json:\"author_name\"`" + `
+}
+
+type BlogHandler struct {
+	db *sql.DB
+}
+
+func (h *BlogHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
+	rows, err := h.db.Query("SELECT id, title, content, author_id FROM posts ORDER BY created_at DESC LIMIT 20")
+	if err != nil {
+		http.Error(w, "internal error", http.StatusInternalServerError)
+		return
+	}
+	defer rows.Close()
+
+	var posts []BlogPost
+	for rows.Next() {
+		var p BlogPost
+		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID); err != nil {
+			http.Error(w, "internal error", http.StatusInternalServerError)
+			return
+		}
+
+		// Get author name for each post
+		var authorName string
+		err := h.db.QueryRow("SELECT name FROM users WHERE id = $1", p.AuthorID).Scan(&authorName)
+		if err != nil {
+			authorName = "Unknown"
+		}
+		p.AuthorName = authorName
+
+		posts = append(posts, p)
+	}
+
+	w.Header().Set("Content-Type", "application/json")
+	json.NewEncoder(w).Encode(posts)
+}`,
			FilePaths:   []string{"internal/handler/blog.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/handler/blog.go",
				LineNumber:  39,
				Content:     "N+1 query problem: executing a separate query for each post's author inside the loop.",
				Category:    model.CategoryPerformance,
				Severity:    model.SeverityCritical,
				Explanation: "For 20 posts, this executes 21 queries (1 for posts + 20 for authors). Use a JOIN in the original query (SELECT p.*, u.name FROM posts p JOIN users u ON p.author_id = u.id) or batch-load authors with IN clause.",
			},
			{
				FilePath:    "internal/handler/blog.go",
				LineNumber:  22,
				Content:     "Loading full content for a listing endpoint. Only title/summary should be fetched for list views.",
				Category:    model.CategoryPerformance,
				Severity:    model.SeverityMinor,
				Explanation: "Blog post content can be very large. For a listing endpoint, fetch only the title and a truncated preview. Full content should only be loaded on the detail view.",
			},
			{
				FilePath:    "internal/handler/blog.go",
				LineNumber:  48,
				Content:     "Missing rows.Err() check after the iteration loop.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "Always check rows.Err() after the loop. Errors during iteration (e.g., connection lost) will cause silent data truncation.",
			},
		},
	}
}

func exercise07RaceCondition() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "In-memory cache with concurrent access and no synchronization",
			Description: "A PR implementing a simple in-memory cache using a map that is accessed from multiple goroutines without any synchronization mechanism.",
			Difficulty:  model.DifficultyIntermediate,
			Category:    model.CategorySecurity,
			CategoryTags: []model.Category{
				model.CategorySecurity,
				model.CategoryDesign,
				model.CategoryPerformance,
			},
			Language: "Go",
			DiffContent: `--- a/internal/cache/cache.go
+++ b/internal/cache/cache.go
@@ -0,0 +1,50 @@
+package cache
+
+import (
+	"time"
+)
+
+type entry struct {
+	value     interface{}
+	expiresAt time.Time
+}
+
+type Cache struct {
+	items map[string]entry
+}
+
+func New() *Cache {
+	c := &Cache{
+		items: make(map[string]entry),
+	}
+	go c.cleanup()
+	return c
+}
+
+func (c *Cache) Get(key string) (interface{}, bool) {
+	item, ok := c.items[key]
+	if !ok {
+		return nil, false
+	}
+	if time.Now().After(item.expiresAt) {
+		delete(c.items, key)
+		return nil, false
+	}
+	return item.value, true
+}
+
+func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
+	c.items[key] = entry{
+		value:     value,
+		expiresAt: time.Now().Add(ttl),
+	}
+}
+
+func (c *Cache) cleanup() {
+	for {
+		time.Sleep(1 * time.Minute)
+		for key, item := range c.items {
+			if time.Now().After(item.expiresAt) {
+				delete(c.items, key)
+			}
+		}
+	}
+}`,
			FilePaths:   []string{"internal/cache/cache.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/cache/cache.go",
				LineNumber:  13,
				Content:     "Data race: map is accessed concurrently from multiple goroutines without sync.Mutex or sync.RWMutex.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "Go maps are not safe for concurrent use. The cleanup goroutine reads/deletes from the map while other goroutines may be calling Get/Set. This will cause a runtime panic ('concurrent map read and map write'). Use sync.RWMutex to protect access.",
			},
			{
				FilePath:    "internal/cache/cache.go",
				LineNumber:  43,
				Content:     "Cleanup goroutine runs forever with no way to stop it, causing a goroutine leak.",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityMajor,
				Explanation: "The goroutine started in New() has no shutdown mechanism. Add a context.Context or done channel and a Close() method so the goroutine can be stopped when the cache is no longer needed.",
			},
			{
				FilePath:    "internal/cache/cache.go",
				LineNumber:  8,
				Content:     "Using interface{} (empty interface) loses type safety. Consider using generics (Go 1.18+).",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityMinor,
				Explanation: "With Go generics, the cache can be parameterized: Cache[T any] to provide type safety at compile time, eliminating the need for type assertions by callers.",
			},
		},
	}
}

func exercise08XSSVulnerability() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Comment rendering with unsanitized HTML output",
			Description: "A PR implementing a comment display component in a web application that renders user-submitted content without proper HTML escaping, enabling cross-site scripting.",
			Difficulty:  model.DifficultyBeginner,
			Category:    model.CategorySecurity,
			CategoryTags: []model.Category{
				model.CategorySecurity,
			},
			Language: "TypeScript",
			DiffContent: `--- a/src/components/CommentList.tsx
+++ b/src/components/CommentList.tsx
@@ -0,0 +1,35 @@
+import { useEffect, useState } from 'react';
+
+interface Comment {
+  id: string;
+  author: string;
+  body: string;
+  created_at: string;
+}
+
+export function CommentList({ postId }: { postId: string }) {
+  const [comments, setComments] = useState<Comment[]>([]);
+
+  useEffect(() => {
+    fetch('/api/posts/' + postId + '/comments')
+      .then(res => res.json())
+      .then(data => setComments(data));
+  }, [postId]);
+
+  return (
+    <div className="comments">
+      {comments.map(comment => (
+        <div key={comment.id} className="comment">
+          <strong>{comment.author}</strong>
+          <div dangerouslySetInnerHTML={{ __html: comment.body }} />
+          <time>{comment.created_at}</time>
+        </div>
+      ))}
+    </div>
+  );
+}`,
			FilePaths:   []string{"src/components/CommentList.tsx"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"TypeScript"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "src/components/CommentList.tsx",
				LineNumber:  24,
				Content:     "XSS vulnerability: dangerouslySetInnerHTML renders raw user input as HTML.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "User-submitted comment body is rendered as raw HTML. An attacker can inject <script> tags or event handlers. Use plain text rendering ({comment.body}) or sanitize with a library like DOMPurify before rendering.",
			},
			{
				FilePath:    "src/components/CommentList.tsx",
				LineNumber:  14,
				Content:     "postId is concatenated directly into the URL without encoding, allowing URL injection.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityMajor,
				Explanation: "Use encodeURIComponent(postId) or template literals with proper URL construction. A malicious postId could manipulate the request URL.",
			},
			{
				FilePath:    "src/components/CommentList.tsx",
				LineNumber:  15,
				Content:     "No error handling for fetch failure. If the API returns an error, the UI fails silently.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "Add .catch() to handle network errors and check res.ok before parsing JSON. Display an error state to the user instead of silently failing.",
			},
		},
	}
}

func exercise09MemoryLeak() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "WebSocket handler that never removes disconnected clients",
			Description: "A PR implementing a WebSocket broadcast server that adds connections to a slice but never removes them when clients disconnect, causing an ever-growing slice and memory leak.",
			Difficulty:  model.DifficultyAdvanced,
			Category:    model.CategoryPerformance,
			CategoryTags: []model.Category{
				model.CategoryPerformance,
				model.CategoryDesign,
				model.CategoryErrorHandling,
			},
			Language: "Go",
			DiffContent: `--- a/internal/ws/hub.go
+++ b/internal/ws/hub.go
@@ -0,0 +1,60 @@
+package ws
+
+import (
+	"log"
+	"net/http"
+
+	"github.com/gorilla/websocket"
+)
+
+var upgrader = websocket.Upgrader{
+	CheckOrigin: func(r *http.Request) bool {
+		return true
+	},
+}
+
+type Hub struct {
+	clients []*websocket.Conn
+}
+
+func NewHub() *Hub {
+	return &Hub{
+		clients: make([]*websocket.Conn, 0),
+	}
+}
+
+func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
+	conn, err := upgrader.Upgrade(w, r, nil)
+	if err != nil {
+		log.Printf("upgrade error: %v", err)
+		return
+	}
+
+	h.clients = append(h.clients, conn)
+
+	go func() {
+		for {
+			_, msg, err := conn.ReadMessage()
+			if err != nil {
+				log.Printf("read error: %v", err)
+				return
+			}
+			h.Broadcast(msg)
+		}
+	}()
+}
+
+func (h *Hub) Broadcast(msg []byte) {
+	for _, client := range h.clients {
+		err := client.WriteMessage(websocket.TextMessage, msg)
+		if err != nil {
+			log.Printf("write error: %v", err)
+		}
+	}
+}`,
			FilePaths:   []string{"internal/ws/hub.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/ws/hub.go",
				LineNumber:  33,
				Content:     "Memory leak: clients are added but never removed when they disconnect.",
				Category:    model.CategoryPerformance,
				Severity:    model.SeverityCritical,
				Explanation: "When a client disconnects (ReadMessage returns error), the goroutine exits but the connection remains in h.clients slice forever. Implement a removal mechanism in the error path of ReadMessage.",
			},
			{
				FilePath:    "internal/ws/hub.go",
				LineNumber:  11,
				Content:     "CheckOrigin always returns true, allowing connections from any origin (CSRF risk).",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityMajor,
				Explanation: "Accepting WebSocket connections from any origin allows cross-site WebSocket hijacking. Validate the origin against an allowlist of trusted domains.",
			},
			{
				FilePath:    "internal/ws/hub.go",
				LineNumber:  17,
				Content:     "Concurrent access to clients slice: append in HandleWS and range in Broadcast can race.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "Multiple goroutines access h.clients without synchronization. Use a mutex to protect the slice, or use a channel-based approach for client registration/deregistration.",
			},
			{
				FilePath:    "internal/ws/hub.go",
				LineNumber:  52,
				Content:     "Failed writes to disconnected clients are logged but the client is not removed.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "When WriteMessage fails, the client is likely disconnected. Close the connection and remove it from the clients list to prevent repeated failed write attempts.",
			},
		},
	}
}

func exercise10MagicNumbers() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Rate limiter with unexplained numeric constants",
			Description: "A PR implementing a token bucket rate limiter where all configuration values are hardcoded as unexplained numeric literals throughout the code.",
			Difficulty:  model.DifficultyBeginner,
			Category:    model.CategoryReadability,
			CategoryTags: []model.Category{
				model.CategoryReadability,
				model.CategoryDesign,
			},
			Language: "Go",
			DiffContent: `--- a/internal/middleware/ratelimit.go
+++ b/internal/middleware/ratelimit.go
@@ -0,0 +1,45 @@
+package middleware
+
+import (
+	"net/http"
+	"sync"
+	"time"
+)
+
+type limiter struct {
+	tokens    float64
+	lastCheck time.Time
+	mu        sync.Mutex
+}
+
+var limiters = make(map[string]*limiter)
+var mu sync.Mutex
+
+func RateLimit(next http.Handler) http.Handler {
+	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		ip := r.RemoteAddr
+
+		mu.Lock()
+		l, exists := limiters[ip]
+		if !exists {
+			l = &limiter{tokens: 100, lastCheck: time.Now()}
+			limiters[ip] = l
+		}
+		mu.Unlock()
+
+		l.mu.Lock()
+		defer l.mu.Unlock()
+
+		elapsed := time.Since(l.lastCheck).Seconds()
+		l.tokens += elapsed * 10
+		if l.tokens > 100 {
+			l.tokens = 100
+		}
+		l.lastCheck = time.Now()
+
+		if l.tokens < 1 {
+			w.Header().Set("Retry-After", "10")
+			http.Error(w, "rate limit exceeded", 429)
+			return
+		}
+
+		l.tokens -= 1
+		next.ServeHTTP(w, r)
+	})
+}`,
			FilePaths:   []string{"internal/middleware/ratelimit.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/middleware/ratelimit.go",
				LineNumber:  25,
				Content:     "Magic numbers: 100 (bucket size), 10 (refill rate), 429, '10' (retry-after) are unexplained constants.",
				Category:    model.CategoryReadability,
				Severity:    model.SeverityMajor,
				Explanation: "Define named constants (e.g., maxTokens, refillRate, retryAfterSeconds) or make these configurable. Anyone reading this code must guess what these numbers mean.",
			},
			{
				FilePath:    "internal/middleware/ratelimit.go",
				LineNumber:  15,
				Content:     "Global limiters map grows indefinitely. Old entries for inactive IPs are never cleaned up.",
				Category:    model.CategoryPerformance,
				Severity:    model.SeverityMajor,
				Explanation: "The limiters map only grows. Add periodic cleanup of expired entries or use a cache with TTL. In production, this will consume unbounded memory.",
			},
			{
				FilePath:    "internal/middleware/ratelimit.go",
				LineNumber:  20,
				Content:     "r.RemoteAddr may include port number and is unreliable behind proxies. Use X-Forwarded-For or X-Real-IP.",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityMajor,
				Explanation: "Behind a load balancer or reverse proxy, RemoteAddr will be the proxy's address, making rate limiting ineffective. Parse X-Forwarded-For header to get the real client IP.",
			},
		},
	}
}

func exercise11InsecureDeserialization() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Configuration loader accepting arbitrary YAML from user input",
			Description: "A PR implementing a configuration update endpoint that deserializes user-submitted YAML directly into internal config structs without validation.",
			Difficulty:  model.DifficultyAdvanced,
			Category:    model.CategorySecurity,
			CategoryTags: []model.Category{
				model.CategorySecurity,
				model.CategoryDesign,
				model.CategoryErrorHandling,
			},
			Language: "Go",
			DiffContent: `--- a/internal/handler/config.go
+++ b/internal/handler/config.go
@@ -0,0 +1,45 @@
+package handler
+
+import (
+	"io"
+	"net/http"
+	"os"
+
+	"gopkg.in/yaml.v3"
+)
+
+type AppConfig struct {
+	Database struct {
+		Host     string ` + "`yaml:\"host\"`" + `
+		Port     int    ` + "`yaml:\"port\"`" + `
+		Password string ` + "`yaml:\"password\"`" + `
+	} ` + "`yaml:\"database\"`" + `
+	Server struct {
+		Port    int    ` + "`yaml:\"port\"`" + `
+		Debug   bool   ` + "`yaml:\"debug\"`" + `
+		DataDir string ` + "`yaml:\"data_dir\"`" + `
+	} ` + "`yaml:\"server\"`" + `
+}
+
+var currentConfig AppConfig
+
+func UpdateConfig(w http.ResponseWriter, r *http.Request) {
+	body, _ := io.ReadAll(r.Body)
+
+	var cfg AppConfig
+	yaml.Unmarshal(body, &cfg)
+
+	currentConfig = cfg
+
+	// Write config to disk
+	out, _ := yaml.Marshal(cfg)
+	os.WriteFile("/etc/app/config.yaml", out, 0644)
+
+	w.WriteHeader(http.StatusOK)
+	w.Write([]byte("config updated"))
+}`,
			FilePaths:   []string{"internal/handler/config.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/handler/config.go",
				LineNumber:  26,
				Content:     "No authentication or authorization. Any user can modify the application configuration.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityCritical,
				Explanation: "Configuration endpoints must be protected with strong authentication and restricted to admin users only. Without auth, any HTTP client can change database credentials, enable debug mode, or modify the data directory.",
			},
			{
				FilePath:    "internal/handler/config.go",
				LineNumber:  27,
				Content:     "Request body is read without size limit. An attacker can send a massive payload to exhaust memory.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityMajor,
				Explanation: "Use io.LimitReader(r.Body, maxBytes) to cap the request size. Without this, a denial-of-service attack can be mounted by sending gigabytes of data.",
			},
			{
				FilePath:    "internal/handler/config.go",
				LineNumber:  30,
				Content:     "No validation of deserialized config values. Arbitrary data is accepted and applied.",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityCritical,
				Explanation: "After unmarshaling, validate all fields: port ranges, non-empty host, data_dir existence, etc. Without validation, invalid config could crash the application or open security holes (e.g., debug=true in production).",
			},
			{
				FilePath:    "internal/handler/config.go",
				LineNumber:  36,
				Content:     "Writing to /etc/app/config.yaml with 0644 means world-readable. Config contains database password.",
				Category:    model.CategorySecurity,
				Severity:    model.SeverityMajor,
				Explanation: "The config file contains the database password in plaintext. Use 0600 permissions to restrict read access to the file owner only. Better yet, don't store passwords in config files at all.",
			},
		},
	}
}

func exercise12UnbufferedChannel() ExerciseWithReviews {
	return ExerciseWithReviews{
		Exercise: model.Exercise{
			Title:       "Job queue with unbuffered channel causing goroutine blocking",
			Description: "A PR implementing a background job processing system using unbuffered channels that blocks producers when the consumer is slow.",
			Difficulty:  model.DifficultyAdvanced,
			Category:    model.CategoryPerformance,
			CategoryTags: []model.Category{
				model.CategoryPerformance,
				model.CategoryDesign,
				model.CategoryErrorHandling,
			},
			Language: "Go",
			DiffContent: `--- a/internal/worker/queue.go
+++ b/internal/worker/queue.go
@@ -0,0 +1,55 @@
+package worker
+
+import (
+	"fmt"
+	"log"
+	"time"
+)
+
+type Job struct {
+	ID      string
+	Payload string
+}
+
+type Queue struct {
+	jobs chan Job
+}
+
+func NewQueue() *Queue {
+	return &Queue{
+		jobs: make(chan Job),
+	}
+}
+
+func (q *Queue) Submit(job Job) {
+	q.jobs <- job
+}
+
+func (q *Queue) Start() {
+	for job := range q.jobs {
+		processJob(job)
+	}
+}
+
+func processJob(job Job) {
+	log.Printf("processing job %s", job.ID)
+
+	// Simulate heavy processing
+	time.Sleep(5 * time.Second)
+
+	result := fmt.Sprintf("completed: %s", job.Payload)
+	log.Println(result)
+}`,
			FilePaths:   []string{"internal/worker/queue.go"},
			Metadata:    json.RawMessage(`{"source":"anonymized-oss","original_language":"Go"}`),
			IsPublished: true,
		},
		Reviews: []model.ReferenceReview{
			{
				FilePath:    "internal/worker/queue.go",
				LineNumber:  20,
				Content:     "Unbuffered channel: Submit() will block the caller until the consumer reads the job.",
				Category:    model.CategoryPerformance,
				Severity:    model.SeverityCritical,
				Explanation: "make(chan Job) creates an unbuffered channel. When processJob takes 5 seconds, any goroutine calling Submit() will block for that duration. Use a buffered channel (make(chan Job, bufferSize)) to decouple producer and consumer.",
			},
			{
				FilePath:    "internal/worker/queue.go",
				LineNumber:  29,
				Content:     "Single worker: only one goroutine processes jobs, creating a bottleneck.",
				Category:    model.CategoryPerformance,
				Severity:    model.SeverityMajor,
				Explanation: "Start() runs a single consumer loop. With 5-second processing time per job, throughput is 0.2 jobs/second. Launch multiple worker goroutines to parallelize processing.",
			},
			{
				FilePath:    "internal/worker/queue.go",
				LineNumber:  34,
				Content:     "No error handling or retry mechanism for failed jobs.",
				Category:    model.CategoryErrorHandling,
				Severity:    model.SeverityMajor,
				Explanation: "processJob doesn't return an error, and the queue has no retry logic, dead-letter queue, or error reporting. Failed jobs are silently lost. Add error handling and a configurable retry policy.",
			},
			{
				FilePath:    "internal/worker/queue.go",
				LineNumber:  18,
				Content:     "No graceful shutdown mechanism. In-flight jobs may be lost on application exit.",
				Category:    model.CategoryDesign,
				Severity:    model.SeverityMajor,
				Explanation: "Add context.Context support for graceful shutdown. When the application is stopping, drain the channel, wait for in-flight jobs to complete, then exit cleanly.",
			},
		},
	}
}
