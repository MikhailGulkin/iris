package scylla

const (
	queryGetMessageByID = `
		SELECT
			id,
			chat_id,
			user_id,
			text,
			content_type,
			send_at,
			created_at,
			deleted_at,
			edit_at
		FROM
		chat.messages
		WHERE id = ?
`
	queryGetMessagesCursor = `
		SELECT
			id,
			chat_id,
			user_id,
			text,
			content_type,
			send_at,
			created_at,
			deleted_at,
			edit_at
		FROM
		chat.messages
		WHERE chat_id = ? AND send_at <= ?
		LIMIT ?
`
	queryGetMessagesByChatID = `
		SELECT
			id,
			chat_id,
			user_id,
			text,
			content_type,
			send_at,
			created_at,
			deleted_at,
			edit_at
		FROM
		chat.messages
		WHERE chat_id = ?
		LIMIT ?
`
)
