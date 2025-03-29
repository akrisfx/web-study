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

func (h *Handler) NetGetUsers(c *gin.Context) {
	users, err := h.GetUsers(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func (h *Handler) NetGetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	user, err := h.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (h *Handler) NetAddUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	createdUser, err := h.AddUser(c.Request.Context(), newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, createdUser)
}

func (h *Handler) NetUpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.UpdateUser(c.Request.Context(), id, updatedUser)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}

func (h *Handler) NetDelUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	err = h.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// Database operations

func (h *Handler) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email, password FROM users`
	rows, err := h.DB.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan user")
		}
		users = append(users, u)
	}

	return users, nil
}

func (h *Handler) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE id = $1`
	row := h.DB.QueryRow(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by id")
	}

	return &user, nil
}

func (h *Handler) AddUser(ctx context.Context, user models.User) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, password) 
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var newID int
	err := h.DB.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&newID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}

	user.ID = newID
	return &user, nil
}

func (h *Handler) UpdateUser(ctx context.Context, id int, user models.User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, password = $3 
		WHERE id = $4
	`

	tag, err := h.DB.Exec(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		id,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (h *Handler) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	tag, err := h.DB.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}