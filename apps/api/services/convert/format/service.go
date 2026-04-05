package format

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v3"

	"requiems-api/platform/httpx"
)

const maxContentSize = 512 * 1024 // 512 KB

// Service converts content between supported data formats.
type Service struct{}

// NewService creates a new format conversion Service.
func NewService() *Service { return &Service{} }

// Convert converts content from one format to another.
func (s *Service) Convert(req Request) (Response, error) {
	if len(req.Content) > maxContentSize {
		return Response{}, &httpx.AppError{
			Status:  http.StatusRequestEntityTooLarge,
			Code:    "content_too_large",
			Message: fmt.Sprintf("content exceeds maximum allowed size of %d bytes", maxContentSize),
		}
	}

	if req.From == req.To {
		return Response{Result: req.Content}, nil
	}

	// Parse input to intermediate representation.
	intermediate, err := parseInput(req.From, req.Content)
	if err != nil {
		return Response{}, err
	}

	// Serialize intermediate to output format.
	result, err := serializeOutput(req.To, intermediate)
	if err != nil {
		return Response{}, err
	}

	return Response{Result: result}, nil
}

// parseInput parses the content string from the given format into a generic
// Go value (map, slice, or scalar).
func parseInput(format, content string) (any, error) {
	switch format {
	case "json":
		return parseJSON(content)
	case "yaml":
		return parseYAML(content)
	case "csv":
		return parseCSV(content)
	case "xml":
		return parseXML(content)
	case "toml":
		return parseTOML(content)
	default:
		return nil, &httpx.AppError{
			Status:  http.StatusBadRequest,
			Code:    "unsupported_format",
			Message: fmt.Sprintf("unsupported input format: %s", format),
		}
	}
}

// serializeOutput converts a generic Go value to the given format string.
func serializeOutput(format string, v any) (string, error) {
	switch format {
	case "json":
		return toJSON(v)
	case "yaml":
		return toYAML(v)
	case "csv":
		return toCSV(v)
	case "xml":
		return toXML(v)
	case "toml":
		return toTOML(v)
	default:
		return "", &httpx.AppError{
			Status:  http.StatusBadRequest,
			Code:    "unsupported_format",
			Message: fmt.Sprintf("unsupported output format: %s", format),
		}
	}
}

// --- JSON ---

func parseJSON(content string) (any, error) {
	dec := json.NewDecoder(strings.NewReader(content))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "invalid_json",
			Message: fmt.Sprintf("invalid JSON: %s", err.Error()),
		}
	}
	return normalizeNumbers(v), nil
}

func toJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", &httpx.AppError{
			Status:  http.StatusInternalServerError,
			Code:    "conversion_error",
			Message: "failed to serialize to JSON",
		}
	}
	return string(b), nil
}

// --- YAML ---

func parseYAML(content string) (any, error) {
	var v any
	if err := yaml.Unmarshal([]byte(content), &v); err != nil {
		return nil, &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "invalid_yaml",
			Message: fmt.Sprintf("invalid YAML: %s", err.Error()),
		}
	}
	// yaml.v3 unmarshals maps as map[string]interface{}, which is what we want.
	return v, nil
}

func toYAML(v any) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", &httpx.AppError{
			Status:  http.StatusInternalServerError,
			Code:    "conversion_error",
			Message: "failed to serialize to YAML",
		}
	}
	return string(b), nil
}

// --- CSV ---

// parseCSV parses CSV content into a slice of maps ([]map[string]string).
// The first row is treated as the header row.
func parseCSV(content string) (any, error) {
	r := csv.NewReader(strings.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		return nil, &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "invalid_csv",
			Message: fmt.Sprintf("invalid CSV: %s", err.Error()),
		}
	}
	if len(records) == 0 {
		return []any{}, nil
	}

	headers := records[0]
	rows := make([]any, 0, len(records)-1)
	for rowIdx, record := range records[1:] {
		if len(record) > len(headers) {
			return nil, &httpx.AppError{
				Status:  http.StatusUnprocessableEntity,
				Code:    "invalid_csv",
				Message: fmt.Sprintf("row %d has %d columns but header defines %d", rowIdx+2, len(record), len(headers)),
			}
		}
		row := make(map[string]any, len(headers))
		for i, h := range headers {
			if i < len(record) {
				row[h] = record[i]
			} else {
				row[h] = ""
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// toCSV serializes a generic value to CSV. The value must be a slice of maps
// with consistent string keys. All keys from the first element are used as headers.
func toCSV(v any) (string, error) {
	rows, ok := v.([]any)
	if !ok {
		return "", &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "conversion_error",
			Message: "CSV output requires a JSON array of objects",
		}
	}
	if len(rows) == 0 {
		return "", nil
	}

	firstRow, ok := rows[0].(map[string]any)
	if !ok {
		return "", &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "conversion_error",
			Message: "CSV output requires a JSON array of objects",
		}
	}

	// Collect headers from first row in deterministic order.
	headers := make([]string, 0, len(firstRow))
	for k := range firstRow {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var buf strings.Builder
	w := csv.NewWriter(&buf)

	if err := w.Write(headers); err != nil {
		return "", &httpx.AppError{
			Status:  http.StatusInternalServerError,
			Code:    "conversion_error",
			Message: "failed to write CSV headers",
		}
	}

	for _, row := range rows {
		m, ok := row.(map[string]any)
		if !ok {
			return "", &httpx.AppError{
				Status:  http.StatusUnprocessableEntity,
				Code:    "conversion_error",
				Message: "CSV output requires all array elements to be objects",
			}
		}
		record := make([]string, len(headers))
		for i, h := range headers {
			if val, exists := m[h]; exists {
				record[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := w.Write(record); err != nil {
			return "", &httpx.AppError{
				Status:  http.StatusInternalServerError,
				Code:    "conversion_error",
				Message: "failed to write CSV row",
			}
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", &httpx.AppError{
			Status:  http.StatusInternalServerError,
			Code:    "conversion_error",
			Message: "failed to flush CSV writer",
		}
	}

	return buf.String(), nil
}

// --- XML ---

// xmlNode is used to build and traverse a generic XML tree.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr
	Children []*xmlNode
	Content  string
}

func (n *xmlNode) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: n.XMLName, Attr: n.Attrs}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if len(n.Children) > 0 {
		for _, child := range n.Children {
			if err := e.EncodeElement(child, xml.StartElement{Name: child.XMLName}); err != nil {
				return err
			}
		}
	} else if n.Content != "" {
		if err := e.EncodeToken(xml.CharData(n.Content)); err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: n.XMLName})
}

// parseXML parses XML content into a generic map structure.
func parseXML(content string) (any, error) {
	dec := xml.NewDecoder(strings.NewReader(content))
	result, err := xmlDecodeElement(dec)
	if err != nil {
		return nil, &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "invalid_xml",
			Message: fmt.Sprintf("invalid XML: %s", err.Error()),
		}
	}
	return result, nil
}

// xmlDecodeElement reads the next XML element from the decoder and converts it
// to a map[string]any, recursively.
func xmlDecodeElement(dec *xml.Decoder) (any, error) {
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			return xmlReadElement(dec, t)
		case xml.CharData:
			s := strings.TrimSpace(string(t))
			if s != "" {
				return s, nil
			}
		}
	}
}

func xmlReadElement(dec *xml.Decoder, start xml.StartElement) (map[string]any, error) {
	result := make(map[string]any)
	var textContent strings.Builder

	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			child, err := xmlReadElement(dec, t)
			if err != nil {
				return nil, err
			}
			key := t.Name.Local
			if existing, ok := result[key]; ok {
				// Multiple children with same tag → convert to slice
				switch ev := existing.(type) {
				case []any:
					result[key] = append(ev, child)
				default:
					result[key] = []any{existing, child}
				}
			} else {
				result[key] = child
			}
		case xml.CharData:
			textContent.Write(t)
		case xml.EndElement:
			text := strings.TrimSpace(textContent.String())
			if len(result) == 0 && text != "" {
				return map[string]any{"#text": text}, nil
			}
			if text != "" {
				result["#text"] = text
			}
			return result, nil
		}
	}
}

// toXML serializes a generic value to XML. Maps are converted to elements;
// the root element is named <root>.
func toXML(v any) (string, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	root := xml.StartElement{Name: xml.Name{Local: "root"}}
	if err := enc.EncodeToken(root); err != nil {
		return "", conversionError("XML")
	}
	if err := encodeXMLValue(enc, v); err != nil {
		if ae, ok := err.(*httpx.AppError); ok {
			return "", ae
		}
		return "", conversionError("XML")
	}
	if err := enc.EncodeToken(root.End()); err != nil {
		return "", conversionError("XML")
	}
	if err := enc.Flush(); err != nil {
		return "", conversionError("XML")
	}

	return buf.String(), nil
}

func encodeXMLValue(enc *xml.Encoder, v any) error {
	switch val := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			elem := xml.StartElement{Name: xml.Name{Local: sanitizeXMLName(k)}}
			if err := enc.EncodeToken(elem); err != nil {
				return err
			}
			if err := encodeXMLValue(enc, val[k]); err != nil {
				return err
			}
			if err := enc.EncodeToken(elem.End()); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range val {
			elem := xml.StartElement{Name: xml.Name{Local: "item"}}
			if err := enc.EncodeToken(elem); err != nil {
				return err
			}
			if err := encodeXMLValue(enc, item); err != nil {
				return err
			}
			if err := enc.EncodeToken(elem.End()); err != nil {
				return err
			}
		}
	case nil:
		// emit nothing
	default:
		if err := enc.EncodeToken(xml.CharData(fmt.Sprintf("%v", val))); err != nil {
			return err
		}
	}
	return nil
}

// sanitizeXMLName replaces characters that are invalid in XML element names
// with underscores. The result is always non-empty.
func sanitizeXMLName(name string) string {
	if name == "" {
		return "_"
	}
	var b strings.Builder
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
			b.WriteRune(r)
		case i > 0 && (r >= '0' && r <= '9' || r == '-' || r == '.'):
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func conversionError(target string) *httpx.AppError {
	return &httpx.AppError{
		Status:  http.StatusInternalServerError,
		Code:    "conversion_error",
		Message: fmt.Sprintf("failed to serialize to %s", target),
	}
}

// --- TOML ---

func parseTOML(content string) (any, error) {
	var v map[string]any
	if _, err := toml.Decode(content, &v); err != nil {
		return nil, &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "invalid_toml",
			Message: fmt.Sprintf("invalid TOML: %s", err.Error()),
		}
	}
	return v, nil
}

func toTOML(v any) (string, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return "", &httpx.AppError{
			Status:  http.StatusUnprocessableEntity,
			Code:    "conversion_error",
			Message: "TOML output requires a JSON object (not an array or scalar)",
		}
	}
	var buf strings.Builder
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		return "", &httpx.AppError{
			Status:  http.StatusInternalServerError,
			Code:    "conversion_error",
			Message: "failed to serialize to TOML",
		}
	}
	return buf.String(), nil
}

// --- helpers ---

// normalizeNumbers converts json.Number values (produced by json.Decoder with
// UseNumber()) to int64 (when the number has no fractional part) or float64,
// so that downstream serializers (YAML, TOML, XML) emit proper numeric values
// rather than quoted strings.
//
// Numbers with no fractional part are converted to int64. Numbers that exceed
// int64 range fall through to float64, which may lose precision for very large
// integers (e.g. values > 2^53). In the rare case both conversions fail the
// original string representation is preserved.
func normalizeNumbers(v any) any {
	switch val := v.(type) {
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return i
		}
		if f, err := val.Float64(); err == nil {
			return f
		}
		return val.String()
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, child := range val {
			out[k] = normalizeNumbers(child)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = normalizeNumbers(item)
		}
		return out
	default:
		return val
	}
}
