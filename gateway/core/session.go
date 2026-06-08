package core

type AuthSession struct {
	AccessToken  string
	OriginalPath string
	Principal    *Principal
}
