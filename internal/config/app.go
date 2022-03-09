package config

//Так как приложение подразумевает в будущем API,
//здесь могут появиться другие важные настройки конфиуграции.
type Application struct {
	Postgres Postgres
}
