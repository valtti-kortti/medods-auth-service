package repository

const (
	insertSessionQuery = `INSERT INTO sessions (user_guid, user_agent, token_hash, ip_address, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING session_id`
	deleteSessionQuery = `DELETE FROM sessions WHERE session_id = $1`
	selectSessionQuery = `SELECT session_id, user_guid, user_agent, token_hash, ip_address, expires_at FROM sessions WHERE session_id = $1`
)
