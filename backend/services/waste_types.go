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

// HTTP Handlers

func (h *Handler) NetGetWasteTypes(c *gin.Context) {
	wasteTypes, err := h.GetWasteTypes(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, wasteTypes)
}

func (h *Handler) NetGetWasteTypeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	wt, err := h.GetWasteTypeByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "waste type not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, wt)
}

func (h *Handler) NetAddWasteType(c *gin.Context) {
	var newWasteType models.WasteType
	if err := c.BindJSON(&newWasteType); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	createdWasteType, err := h.AddWasteType(c.Request.Context(), newWasteType)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		
		return
	}

	c.IndentedJSON(http.StatusCreated, createdWasteType)
}

func (h *Handler) NetUpdateWasteType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	var updatedWasteType models.WasteType
	if err := c.BindJSON(&updatedWasteType); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.UpdateWasteType(c.Request.Context(), id, updatedWasteType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "waste type not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedWasteType)
}

func (h *Handler) NetDelWasteType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.DeleteWasteType(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "waste type not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "waste type deleted"})
}

// Database operations

func (h *Handler) GetWasteTypes(ctx context.Context) ([]models.WasteType, error) {
	query := `SELECT id, name FROM waste_type`
	rows, err := h.DB.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get waste types")
	}
	defer rows.Close()

	var wasteTypes []models.WasteType
	for rows.Next() {
		var wt models.WasteType
		err := rows.Scan(&wt.ID, &wt.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan waste type")
		}
		wasteTypes = append(wasteTypes, wt)
	}

	return wasteTypes, nil
}

func (h *Handler) GetWasteTypeByID(ctx context.Context, id int) (*models.WasteType, error) {
	query := `SELECT id, name FROM waste_type WHERE id = $1`
	row := h.DB.QueryRow(ctx, query, id)

	var wt models.WasteType
	err := row.Scan(&wt.ID, &wt.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get waste type by id")
	}

	return &wt, nil
}

func (h *Handler) AddWasteType(ctx context.Context, wt models.WasteType) (*models.WasteType, error) {
	query := `
		INSERT INTO waste_type (name) 
		VALUES ($1)
		RETURNING id
	`

	var newID int
	err := h.DB.QueryRow(
		ctx,
		query,
		wt.Name,
	).Scan(&newID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert waste type")
	}

	wt.ID = newID
	return &wt, nil
}

func (h *Handler) UpdateWasteType(ctx context.Context, id int, wt models.WasteType) error {
	query := `
		UPDATE waste_type 
		SET name = $1
		WHERE id = $2
	`

	tag, err := h.DB.Exec(
		ctx,
		query,
		wt.Name,
		id,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update waste type")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (h *Handler) DeleteWasteType(ctx context.Context, id int) error {
	query := `DELETE FROM waste_type WHERE id = $1`
	tag, err := h.DB.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete waste type")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}