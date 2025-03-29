package services

import (
	"backend/models"
	"context"

	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	"github.com/pkg/errors"
)

// GET all collection points
func (h *Handler) NetGetCollectionPoints(c *gin.Context) {
	points, err := h.GetCollectionPoints(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, points)
}

// GET collection point by ID
func (h *Handler) NetGetCollectionPointByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	cp, err := h.GetCollectionPointByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "collection point not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, cp)
}

// POST a new collection point
func (h *Handler) NetAddCollectionPoint(c *gin.Context) {
	var newCollectionPoint models.CollectionPoint
	if err := c.BindJSON(&newCollectionPoint); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	createdPoint, err := h.AddCollectionPoint(c.Request.Context(), newCollectionPoint)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusCreated, createdPoint)
}

// PUT/UPDATE a collection point by ID
func (h *Handler) NetUpdateCollectionPoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	var updatedCollectionPoint models.CollectionPoint
	if err := c.BindJSON(&updatedCollectionPoint); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.UpdateCollectionPoint(c.Request.Context(), id, updatedCollectionPoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "collection point not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedCollectionPoint)
}

// DELETE a collection point by ID
func (h *Handler) NetDelCollectionPoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.DeleteCollectionPoint(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "collection point not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "collection point deleted"})
}

// Database operations

func (h *Handler) GetCollectionPoints(ctx context.Context) ([]models.CollectionPoint, error) {
	query := `SELECT id, name, address, lat, long FROM collection_point`
	rows, err := h.DB.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get collection points")
	}
	defer rows.Close()

	var points []models.CollectionPoint
	for rows.Next() {
		var cp models.CollectionPoint
		err := rows.Scan(&cp.ID, &cp.Name, &cp.Address, &cp.Lat, &cp.Long)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan collection point")
		}
		points = append(points, cp)
	}

	return points, nil
}

func (h *Handler) GetCollectionPointByID(ctx context.Context, id int) (*models.CollectionPoint, error) {
	query := `SELECT id, name, address, lat, long FROM collection_point WHERE id = $1`
	row := h.DB.QueryRow(ctx, query, id)

	var cp models.CollectionPoint
	err := row.Scan(&cp.ID, &cp.Name, &cp.Address, &cp.Lat, &cp.Long)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get collection point by id")
	}

	return &cp, nil
}

func (h *Handler) AddCollectionPoint(ctx context.Context, cp models.CollectionPoint) (*models.CollectionPoint, error) {
	query := `
		INSERT INTO collection_point (name, address, lat, long) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := h.DB.QueryRow(
		ctx,
		query,
		cp.Name,
		cp.Address,
		cp.Lat,
		cp.Long,
	).Scan(&id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert collection point")
	}

	cp.ID = id
	return &cp, nil
}

func (h *Handler) UpdateCollectionPoint(ctx context.Context, id int, cp models.CollectionPoint) error {
	query := `
		UPDATE collection_point 
		SET name = $1, address = $2, lat = $3, long = $4 
		WHERE id = $5
	`

	tag, err := h.DB.Exec(
		ctx,
		query,
		cp.Name,
		cp.Address,
		cp.Lat,
		cp.Long,
		id,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update collection point")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (h *Handler) DeleteCollectionPoint(ctx context.Context, id int) error {
	query := `DELETE FROM collection_point WHERE id = $1`
	tag, err := h.DB.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete collection point")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}