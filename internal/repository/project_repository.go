package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"prosamik-backend/internal/database"
	"prosamik-backend/pkg/models"
	"strings"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		db: database.DB,
	}
}

var closeErrProject error

// Helper function to normalize strings
func normalizeProjectString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// GetProjectByTitle retrieves a project by its title
func (r *ProjectRepository) GetProjectByTitle(title string) (*models.Project, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM projects
        WHERE LOWER(title) = LOWER($1)
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	var closeErrProject error
	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	project := &models.Project{}
	err = stmt.QueryRow(normalizeProjectString(title)).Scan(
		&project.ID,
		&project.Title,
		&project.Path,
		&project.Description,
		&project.Tags,
		&project.ViewsCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return project, nil
}

// GetAllProjects retrieves all project posts
func (r *ProjectRepository) GetAllProjects() ([]*models.Project, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM projects
        ORDER BY id DESC
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID,
			&project.Title,
			&project.Path,
			&project.Description,
			&project.Tags,
			&project.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// CreateProject adds a new project post
func (r *ProjectRepository) CreateProject(project *models.Project) error {
	query := `
        INSERT INTO projects (title, path, description, tags)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	err = stmt.QueryRow(
		strings.TrimSpace(project.Title),
		strings.TrimSpace(project.Path),
		strings.TrimSpace(project.Description),
		strings.TrimSpace(project.Tags),
	).Scan(&project.ID)

	if err != nil {
		return fmt.Errorf("create project error: %w", err)
	}

	return nil
}

// UpdateProject updates an existing project post
func (r *ProjectRepository) UpdateProject(project *models.Project) error {
	query := `
        UPDATE projects
        SET title = $1, path = $2, description = $3, tags = $4
        WHERE id = $5
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	result, err := stmt.Exec(
		strings.TrimSpace(project.Title),
		strings.TrimSpace(project.Path),
		strings.TrimSpace(project.Description),
		strings.TrimSpace(project.Tags),
		project.ID,
	)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no project found with id: %d", project.ID)
	}

	return nil
}

// DeleteProject removes a project post
func (r *ProjectRepository) DeleteProject(id int64) error {
	query := `
        DELETE FROM projects
        WHERE id = $1
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no project found with id: %d", id)
	}

	return nil
}

// SearchProjects searches for projects by title, path, tags, or description
func (r *ProjectRepository) SearchProjects(query string) ([]*models.Project, error) {
	searchQuery := `
        SELECT id, title, path, description, tags, views_count
        FROM projects
        WHERE LOWER(title) LIKE LOWER($1) 
           OR LOWER(path) LIKE LOWER($1)
           OR LOWER(tags) LIKE LOWER($1)
           OR LOWER(description) LIKE LOWER($1)
        ORDER BY id DESC
    `

	stmt, err := r.db.Prepare(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	normalizedQuery := "%" + normalizeProjectString(query) + "%"
	rows, err := stmt.Query(normalizedQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID,
			&project.Title,
			&project.Path,
			&project.Description,
			&project.Tags,
			&project.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return projects, nil
}

// GetProject retrieves a single project by ID
func (r *ProjectRepository) GetProject(id int64) (*models.Project, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM projects
        WHERE id = $1
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	project := &models.Project{}
	err = stmt.QueryRow(id).Scan(
		&project.ID,
		&project.Title,
		&project.Path,
		&project.Description,
		&project.Tags,
		&project.ViewsCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return project, nil
}

// GetProjectByPath retrieves a project by its path
func (r *ProjectRepository) GetProjectByPath(path string) (*models.Project, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM projects
        WHERE path = $1
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmtProject(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErrProject == nil {
				closeErrProject = cerr
			}
		}
	}()

	project := &models.Project{}
	err = stmt.QueryRow(path).Scan(
		&project.ID,
		&project.Title,
		&project.Path,
		&project.Description,
		&project.Tags,
		&project.ViewsCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return project, nil
}

// Helper function to handle statement closing
func closeStmtProject(stmt *sql.Stmt) error {
	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error closing statement: %w", err)
	}
	return nil
}