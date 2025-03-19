package db

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func Init() {
	// Connect to DB
	var err error

	DB, err = Connect()
	if err != nil {
		log.Fatalf("Error connecting to db :%v", err)
	}
}

// Insert a new uploaded image into DB
func SaveImage(imageKey string) error {
	query := `INSERT INTO uploaded_images (image_key) VALUES (?);`

	_, err := DB.Exec(query, imageKey)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetAllImages(cdnUrl string) ([]string, error) {
	// query all users
	res, err := DB.Query("SELECT * FROM uploaded_images")
	if err != nil {
		log.Printf("error querying db: %v", err)
		return nil, err
	}
	defer res.Close()

	var images []string

	for res.Next() {
		var id int
		var imageKey string

		if err := res.Scan(&id, &imageKey); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Construct CDN URL for the image
		imageCDNUrl := fmt.Sprintf("%s/%s", cdnUrl, imageKey)

		images = append(images, imageCDNUrl)
	}

	return images, nil
}
