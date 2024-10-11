package main

import "snippetbox.haonguyen.tech/internal/models"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}
