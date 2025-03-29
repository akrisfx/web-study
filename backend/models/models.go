package models

import "time"

// // WasteType represents data about a type of waste.
// type WasteType struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

// // CollectionPoint represents data about a collection point.
// type CollectionPoint struct {
// 	ID      string  `json:"id"`
// 	Name    string  `json:"name"`
// 	Address string  `json:"address"`
// 	Lat     float64 `json:"lat"`  // Latitude
// 	Long    float64 `json:"long"` // Longitude
// }

// // User represents data about a User.
// type User struct {
// 	ID       string `json:"id"`
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// // RecyclingReport represents data about a recycling report.
// type RecyclingReport struct {
// 	ID                string  `json:"id"`
// 	UserID            string  `json:"user_id"`             // user
// 	CollectionPointid int  `json:"collection_point_id"` // collectionPoint
// 	WasteTypeID       string  `json:"waste_type_id"`       // wasteType
// 	Quantity          float64 `json:"quantity"`
// 	Date              string  `json:"date"`
// }

type WasteType struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type CollectionPoint struct {
	ID      int     `db:"id" json:"id"`
	Name    string  `db:"name" json:"name"`
	Address string  `db:"address" json:"address"`
	Lat     float64 `db:"lat" json:"lat"`
	Long    float64 `db:"long" json:"long"`
}

type User struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
}

type RecyclingReport struct {
	ID                int     `db:"id" json:"id"`
	UserID            int     `db:"user_id" json:"user_id"`
	CollectionPointID int     `db:"collection_point_id" json:"collection_point_id"`
	WasteTypeID       int     `db:"waste_type_id" json:"waste_type_id"`
	Quantity          float64 `db:"quantity" json:"quantity"`
	// Date              string  `db:"date" json:"date"`
	Date              time.Time  `db:"date" json:"date"`
}
