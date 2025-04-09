package sl

import (
	"log/slog"
)

func ErisErr(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
