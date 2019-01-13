package sunshinemotion

func PasswordHash(passwordHash string) string {
	return md5String(passwordHash)
}
