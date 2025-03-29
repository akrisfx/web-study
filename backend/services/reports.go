package services

import (
	"context"
	"net/http"
	"strconv"

	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

// Обработчики HTTP запросов

func (h *Handler) NetGetRecyclingReports(c *gin.Context) {
	reports, err := h.GetRecyclingReports(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, reports)
}

func (h *Handler) NetGetRecyclingReportsByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	report, err := h.GetRecyclingReportByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "report not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, report)
}

func (h *Handler) NetAddRecyclingReport(c *gin.Context) {
	var newReport models.RecyclingReport
	if err := c.BindJSON(&newReport); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	createdReport, err := h.AddRecyclingReport(c.Request.Context(), newReport)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, createdReport)
}

func (h *Handler) NetUpdateRecyclingReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	var updatedReport models.RecyclingReport
	if err := c.BindJSON(&updatedReport); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.UpdateRecyclingReport(c.Request.Context(), id, updatedReport)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "report not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedReport)
}

func (h *Handler) NetDeleteRecyclingReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.DeleteRecyclingReport(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "report not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "report deleted"})
}

// Методы работы с базой данных
// Получаем все отчеты о переработке
func (h *Handler) GetRecyclingReports(ctx context.Context) ([]models.RecyclingReport, error) {
	query := `
		SELECT id, user_id, collection_point_id, waste_type_id, quantity, date 
		FROM recycling_report
	`

	rows, err := h.DB.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reports")
	}
	defer rows.Close()

	var reports []models.RecyclingReport
	for rows.Next() {
		var r models.RecyclingReport
		err := rows.Scan(
			&r.ID,
			&r.UserID,
			&r.CollectionPointID,
			&r.WasteTypeID,
			&r.Quantity,
			&r.Date, // здесь используем time.Time
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan report")
		}
		reports = append(reports, r)
	}

	return reports, nil
}

// Получаем отчет по ID
func (h *Handler) GetRecyclingReportByID(ctx context.Context, id int) (*models.RecyclingReport, error) {
	query := `
		SELECT id, user_id, collection_point_id, waste_type_id, quantity, date 
		FROM recycling_report 
		WHERE id = $1
	`

	row := h.DB.QueryRow(ctx, query, id)
	var report models.RecyclingReport
	err := row.Scan(
		&report.ID,
		&report.UserID,
		&report.CollectionPointID,
		&report.WasteTypeID,
		&report.Quantity,
		&report.Date, // здесь используем time.Time
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get report by id")
	}

	return &report, nil
}

// Добавляем новый отчет
func (h *Handler) AddRecyclingReport(ctx context.Context, report models.RecyclingReport) (*models.RecyclingReport, error) {
	query := `
		INSERT INTO recycling_report 
			(user_id, collection_point_id, waste_type_id, quantity, date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var newID int
	err := h.DB.QueryRow(ctx, query,
		report.UserID,
		report.CollectionPointID,
		report.WasteTypeID,
		report.Quantity,
		report.Date, // передаем time.Time
	).Scan(&newID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert report")
	}

	report.ID = newID
	return &report, nil
}

// Обновляем отчет
func (h *Handler) UpdateRecyclingReport(ctx context.Context, id int, report models.RecyclingReport) error {
	query := `
		UPDATE recycling_report 
		SET 
			user_id = $1,
			collection_point_id = $2,
			waste_type_id = $3,
			quantity = $4,
			date = $5
		WHERE id = $6
	`

	tag, err := h.DB.Exec(ctx, query,
		report.UserID,
		report.CollectionPointID,
		report.WasteTypeID,
		report.Quantity,
		report.Date, // передаем time.Time
		id,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update report")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// Удаляем отчет
func (h *Handler) DeleteRecyclingReport(ctx context.Context, id int) error {
	query := `DELETE FROM recycling_report WHERE id = $1`
	tag, err := h.DB.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete report")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}


// func (h *Handler) GetRecyclingReports(ctx context.Context) ([]models.RecyclingReport, error) {
// 	query := `
// 		SELECT id, user_id, collection_point_id, waste_type_id, quantity, date 
// 		FROM recycling_report
// 	`

// 	rows, err := h.DB.Query(ctx, query)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to get reports")
// 	}
// 	defer rows.Close()

// 	var reports []models.RecyclingReport
// 	for rows.Next() {
// 		var r models.RecyclingReport
// 		err := rows.Scan(
// 			&r.ID,
// 			&r.UserID,
// 			&r.CollectionPointID,
// 			&r.WasteTypeID,
// 			&r.Quantity,
// 			&r.Date,
// 		)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "failed to scan report")
// 		}
// 		reports = append(reports, r)
// 	}

// 	return reports, nil
// }

// func (h *Handler) GetRecyclingReportByID(ctx context.Context, id int) (*models.RecyclingReport, error) {
// 	query := `
// 		SELECT id, user_id, collection_point_id, waste_type_id, quantity, date 
// 		FROM recycling_report 
// 		WHERE id = $1
// 	`

// 	row := h.DB.QueryRow(ctx, query, id)
// 	var report models.RecyclingReport
// 	err := row.Scan(
// 		&report.ID,
// 		&report.UserID,
// 		&report.CollectionPointID,
// 		&report.WasteTypeID,
// 		&report.Quantity,
// 		&report.Date,
// 	)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to get report by id")
// 	}

// 	return &report, nil
// }

// func (h *Handler) AddRecyclingReport(ctx context.Context, report models.RecyclingReport) (*models.RecyclingReport, error) {
// 	query := `
// 		INSERT INTO recycling_report 
// 			(user_id, collection_point_id, waste_type_id, quantity, date)
// 		VALUES ($1, $2, $3, $4, $5)
// 		RETURNING id
// 	`

// 	var newID int
// 	err := h.DB.QueryRow(ctx, query,
// 		report.UserID,
// 		report.CollectionPointID,
// 		report.WasteTypeID,
// 		report.Quantity,
// 		report.Date,
// 	).Scan(&newID)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to insert report")
// 	}

// 	report.ID = newID
// 	return &report, nil
// }

// func (h *Handler) UpdateRecyclingReport(ctx context.Context, id int, report models.RecyclingReport) error {
// 	query := `
// 		UPDATE recycling_report 
// 		SET 
// 			user_id = $1,
// 			collection_point_id = $2,
// 			waste_type_id = $3,
// 			quantity = $4,
// 			date = $5
// 		WHERE id = $6
// 	`

// 	tag, err := h.DB.Exec(ctx, query,
// 		report.UserID,
// 		report.CollectionPointID,
// 		report.WasteTypeID,
// 		report.Quantity,
// 		report.Date,
// 		id,
// 	)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to update report")
// 	}

// 	if tag.RowsAffected() == 0 {
// 		return pgx.ErrNoRows
// 	}

// 	return nil
// }

// func (h *Handler) DeleteRecyclingReport(ctx context.Context, id int) error {
// 	query := `DELETE FROM recycling_report WHERE id = $1`
// 	tag, err := h.DB.Exec(ctx, query, id)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to delete report")
// 	}

// 	if tag.RowsAffected() == 0 {
// 		return pgx.ErrNoRows
// 	}

// 	return nil
// }