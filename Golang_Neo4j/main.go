package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	r := gin.Default()

	// Inisialisasi koneksi ke Neo4j
	dbUri := "neo4j://192.168.18.191:7687"
	dbUser := "neo4j"
	dbPassword := "ddi123"
	dbName := "kemenag"
	dbFullUri := dbUri + "/" + dbName

	driver, err := neo4j.NewDriver(
		dbFullUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""),
	)
	if err != nil {
		panic(err)
	}
	defer driver.Close()

	// Endpoint untuk halaman web
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Endpoint /api/jemaah untuk mengambil data dari Neo4j
	r.GET("/api/jemaah", func(c *gin.Context) {
		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: dbName})
		defer session.Close()

		query := "MATCH (n:Jemaah) RETURN n.ktp, n.nama, n.gender, n.passport, n.umur, n.latitude, n.longitude"
		result, err := session.Run(query, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Neo4j"})
			return
		}

		var data []map[string]interface{}
		for result.Next() {
			record := result.Record()
			ktp, found := record.Get("n.ktp")
			nama, _ := record.Get("n.nama")
			gender, _ := record.Get("n.gender")
			passport, _ := record.Get("n.passport")
			umur, _ := record.Get("n.umur")
			latitude, _ := record.Get("n.latitude")
			longitude, _ := record.Get("n.longitude")

			if found {
				data = append(data, map[string]interface{}{
					"ktp":       ktp,
					"nama":      nama,
					"gender":    gender,
					"passport":  passport,
					"umur":      umur,
					"latitude":  latitude,
					"longitude": longitude,
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{"data": data})
	})

	// Menjalankan server web
	port := ":8080"
	r.LoadHTMLGlob("mywebapp/templates/*") // Set folder templates sebagai lokasi templat HTML
	fmt.Printf("Server is running on port %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
