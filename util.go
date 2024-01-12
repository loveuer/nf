package nf

import (
	"fmt"
	"github.com/loveuer/nf/internal/schema"
	"strings"
)

const (
	MIMETextXML         = "text/xml"
	MIMETextHTML        = "text/html"
	MIMETextPlain       = "text/plain"
	MIMETextJavaScript  = "text/javascript"
	MIMEApplicationXML  = "application/xml"
	MIMEApplicationJSON = "application/json"
	// Deprecated: use MIMETextJavaScript instead
	MIMEApplicationJavaScript = "application/javascript"
	MIMEApplicationForm       = "application/x-www-form-urlencoded"
	MIMEOctetStream           = "application/octet-stream"
	MIMEMultipartForm         = "multipart/form-data"

	MIMETextXMLCharsetUTF8         = "text/xml; charset=utf-8"
	MIMETextHTMLCharsetUTF8        = "text/html; charset=utf-8"
	MIMETextPlainCharsetUTF8       = "text/plain; charset=utf-8"
	MIMETextJavaScriptCharsetUTF8  = "text/javascript; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8  = "application/xml; charset=utf-8"
	MIMEApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"
	// Deprecated: use MIMETextJavaScriptCharsetUTF8 instead
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
)

func verifyHandlers(path string, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic(fmt.Sprintf("missing handler in route: %s", path))
	}

	for _, handler := range handlers {
		if handler == nil {
			panic(fmt.Sprintf("nil handler found in route: %s", path))
		}
	}
}

// parseVendorSpecificContentType check if content type is vendor specific and
// if it is parsable to any known types. If it's not vendor specific then returns
// the original content type.
func parseVendorSpecificContentType(cType string) string {
	plusIndex := strings.Index(cType, "+")

	if plusIndex == -1 {
		return cType
	}

	var parsableType string
	if semiColonIndex := strings.Index(cType, ";"); semiColonIndex == -1 {
		parsableType = cType[plusIndex+1:]
	} else if plusIndex < semiColonIndex {
		parsableType = cType[plusIndex+1 : semiColonIndex]
	} else {
		return cType[:semiColonIndex]
	}

	slashIndex := strings.Index(cType, "/")

	if slashIndex == -1 {
		return cType
	}

	return cType[0:slashIndex+1] + parsableType
}

func parseToStruct(aliasTag string, out interface{}, data map[string][]string) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.SetAliasTag(aliasTag)

	if err := schemaDecoder.Decode(out, data); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	return nil
}
