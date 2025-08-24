package decoder

import "github.com/google/uuid"

// parseUUID парсит строку в UUID
func parseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}
