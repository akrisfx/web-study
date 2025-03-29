package main

// import (
// 	"backend/models"
// 	"backend/services"
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/pashagolub/pgxmock"
// 	"github.com/stretchr/testify/assert"
// )

// func setupRouter() *gin.Engine {
// 	pool := ConnectToDB()
// 	defer pool.Close()
// 	router := gin.Default()
// 	h := &services.Handler{DB: pool}
// 	router.GET("/api/reports", h.NetGetRecyclingReports)
// 	router.GET("/api/reports/:id", h.NetGetRecyclingReportsByID)
// 	router.POST("/api/reports", h.NetAddRecyclingReport)
// 	router.PUT("/api/reports/:id", h.NetUpdateRecyclingReport)
// 	router.DELETE("/api/reports/:id", h.NetDeleteRecyclingReport)

// 	// Waste Types
// 	router.GET("/api/waste-types", h.NetGetWasteTypes)
// 	router.GET("/api/waste-types/:id", h.NetGetWasteTypeByID)
// 	router.POST("/api/waste-types", h.NetAddWasteType)
// 	router.PUT("/api/waste-types/:id", h.NetUpdateWasteType)
// 	router.DELETE("/api/waste-types/:id", h.NetDelWasteType)

// 	// Collection Points
// 	router.GET("/api/collection-points", h.NetGetCollectionPoints)
// 	router.GET("/api/collection-points/:id", h.NetGetCollectionPointByID)
// 	router.POST("/api/collection-points", h.NetAddCollectionPoint)
// 	router.PUT("/api/collection-points/:id", h.NetUpdateCollectionPoint)
// 	router.DELETE("/api/collection-points/:id", h.NetDelCollectionPoint)

// 	// Users
// 	router.GET("/api/users", h.NetGetUsers)
// 	router.GET("/api/users/:id", h.NetGetUserByID)
// 	router.POST("/api/users", h.NetAddUser)
// 	router.PUT("/api/users/:id", h.NetUpdateUser)
// 	router.DELETE("/api/users/:id", h.NetDelUser)

// 	return router
// }

// func TestNetGetCollectionPoints_Success(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer mock.Close()

// 	expected := []models.CollectionPoint{
// 		{ID: 1, Name: "Point A", Address: "123 Main St", Lat: 40.7128, Long: -74.0060},
// 	}

// 	rows := pgxmock.NewRows([]string{"id", "name", "address", "lat", "long"}).
// 		AddRow(1, "Point A", "123 Main St", 40.7128, -74.0060)

// 	mock.ExpectQuery("SELECT id, name, address, lat, long FROM collection_points").
// 		WillReturnRows(rows)

// 	handler := &services.Handler{DB: mock}
// 	router := setupRouter(handler)

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/collection-points", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response []models.CollectionPoint
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, response)
// }

// func TestNetGetCollectionPoints_DBError(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer mock.Close()

// 	mock.ExpectQuery("SELECT id, name, address, lat, long FROM collection_points").
// 		WillReturnError(pgx.ErrTxClosed)

// 	handler := &Handler{DB: mock}
// 	router := setupRouter(handler)

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/collection-points", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// }

// func TestNetGetCollectionPointByID_Success(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer mock.Close()

// 	expected := models.CollectionPoint{ID: 1, Name: "Point A", Address: "123 Main St", Lat: 40.7128, Long: -74.0060}

// 	row := pgxmock.NewRows([]string{"id", "name", "address", "lat", "long"}).
// 		AddRow(1, "Point A", "123 Main St", 40.7128, -74.0060)

// 	mock.ExpectQuery("SELECT id, name, address, lat, long FROM collection_points WHERE id = ?").
// 		WithArgs(1).
// 		WillReturnRows(row)

// 	handler := &Handler{DB: mock}
// 	router := setupRouter(handler)

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/collection-points/1", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response models.CollectionPoint
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, response)
// }

// func TestNetAddCollectionPoint_Success(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer mock.Close()

// 	newPoint := models.CollectionPoint{
// 		Name:    "New Point",
// 		Address: "456 Elm St",
// 		Lat:     34.0522,
// 		Long:    -118.2437,
// 	}

// 	mock.ExpectQuery("INSERT INTO collection_points").
// 		WithArgs(newPoint.Name, newPoint.Address, newPoint.Lat, newPoint.Long).
// 		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

// 	handler := &Handler{DB: mock}
// 	router := setupRouter(handler)

// 	body, _ := json.Marshal(newPoint)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/collection-points", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusCreated, w.Code)

// 	var response models.CollectionPoint
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, response.ID)
// 	assert.Equal(t, newPoint.Name, response.Name)
// }

// func TestNetUpdateCollectionPoint_NotFound(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer mock.Close()

// 	updatedPoint := models.CollectionPoint{
// 		Name:    "Updated Point",
// 		Address: "789 Oak St",
// 		Lat:     41.8781,
// 		Long:    -87.6298,
// 	}

// 	mock.ExpectExec("UPDATE collection_points").
// 		WithArgs(updatedPoint.Name, updatedPoint.Address, updatedPoint.Lat, updatedPoint.Long, 1).
// 		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

// 	handler := &Handler{DB: mock}
// 	router := setupRouter(handler)

// 	body, _ := json.Marshal(updatedPoint)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("PUT", "/collection-points/1", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// }

// func TestNetDelCollectionPoint_Success(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer mock.Close()

// 	mock.ExpectExec("DELETE FROM collection_points WHERE id = ?").
// 		WithArgs(1).
// 		WillReturnResult(pgxmock.NewResult("DELETE", 1))

// 	handler := &Handler{DB: mock}
// 	router := setupRouter(handler)

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("DELETE", "/collection-points/1", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]string
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "collection point deleted", response["message"])
// }

// // func TestGetRecyclingReports(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/reports", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestGetRecyclingReportByID(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/reports/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestAddRecyclingReport(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":5,"user_id":1,"collection_point_id":1,"waste_type_id":1,"quantity":5.0,"date":"2023-01-05"}`)
// // 	req, _ := http.NewRequest("POST", "/api/reports", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusCreated, w.Code)
// // }

// // func TestUpdateRecyclingReport(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":1,"user_id":1,"collection_point_id":1,"waste_type_id":1,"quantity":10.0,"date":"2023-01-01"}`)
// // 	req, _ := http.NewRequest("PUT", "/api/reports/1", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestDeleteRecyclingReport(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("DELETE", "/api/reports/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestGetWasteTypes(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/waste-types", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestGetWasteTypeByID(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/waste-types/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestAddWasteType(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":5,"name":"Organic"}`)
// // 	req, _ := http.NewRequest("POST", "/api/waste-types", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusCreated, w.Code)
// // }

// // func TestUpdateWasteType(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":1,"name":"Updated Plastic"}`)
// // 	req, _ := http.NewRequest("PUT", "/api/waste-types/1", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestDeleteWasteType(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("DELETE", "/api/waste-types/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestGetCollectionPoints(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/collection-points", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestGetCollectionPointByID(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/collection-points/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestAddCollectionPoint(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":5,"name":"Point E","address":"202 Birch St","lat":40.7128,"long":-74.0060}`)
// // 	req, _ := http.NewRequest("POST", "/api/collection-points", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusCreated, w.Code)
// // }

// // func TestUpdateCollectionPoint(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":1,"name":"Updated Point A","address":"123 Main St","lat":40.7128,"long":-74.0060}`)
// // 	req, _ := http.NewRequest("PUT", "/api/collection-points/1", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestDeleteCollectionPoint(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("DELETE", "/api/collection-points/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestGetUsers(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/users", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestGetUserByID(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("GET", "/api/users/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestAddUser(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":5,"name":"Eve","email":"eve@example.com","password":"password123"}`)
// // 	req, _ := http.NewRequest("POST", "/api/users", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusCreated, w.Code)
// // }

// // func TestUpdateUser(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	body := bytes.NewBufferString(`{"id":1,"name":"Updated Alice","email":"alice@example.com","password":"password123"}`)
// // 	req, _ := http.NewRequest("PUT", "/api/users/1", body)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }

// // func TestDeleteUser(t *testing.T) {
// // 	router := setupRouter()
// // 	w := httptest.NewRecorder()
// // 	req, _ := http.NewRequest("DELETE", "/api/users/1", nil)
// // 	router.ServeHTTP(w, req)

// // 	assert.Equal(t, http.StatusOK, w.Code)
// // }
