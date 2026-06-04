package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/abrarr21/notes-in-golang/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type stringFields struct {
	value *string // nil means client did not sent this
	key   string  // MongoDB dot-Notation key e.g "title"
	label string  // used in error message
}

func buildUpdateDoc(input models.UpdateNoteRequest) (bson.M, error) {
	update := bson.M{}

	stringFields := []stringFields{
		{input.Title, "title", "Title"},
		{input.Content, "content", "Content"},
	}

	for _, f := range stringFields {
		if f.value == nil {
			continue // field not sent by client, skip entirely
		}

		val := strings.TrimSpace(*f.value)
		if val == "" {
			return nil, fmt.Errorf("%s cannot be empty", f.label)
		}

		update[f.key] = val
	}

	if len(update) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	update["updated_at"] = time.Now().UTC()
	return update, nil
}
