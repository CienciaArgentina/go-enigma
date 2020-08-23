package domain

const (
	// Request.
	ErrInvalidBody     = "El cuerpo del mensaje que intentás enviar no es válido"
	ErrInvalidBodyCode = "invalid_body"

	ErrInternalCode = "internal_error"

	// Empty.
	ErrEmptyField              = "Hay algún campo vacío y no puede estarlo"
	ErrEmptyFieldCode          = "empty_field"
	ErrEmptyUsername           = "El nombre de usuario no puede estar vacío"
	ErrEmptyPassword           = "La contraseña no puede estar vacía"
	ErrEmptyEmail              = "El email no puede estar vacío"
	ErrEmptyEmailCode          = "empty_email"
	ErrEmptyFieldUserCodeLogin = "invalid_user_login"

	// General.
	ErrUnexpectedError = "Ocurrió un error en el sistema, por favor, ponete en contacto con sistemas"
)
