package resultswriters

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/USACE/go-consequences/consequences"
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/parquet/file"
)

type geoParquetResultsWriter struct {
	filepath string
	w        file.Writer
}

func InitGeoParquetResultsWriterFromFile(filepath string) (*geoParquetResultsWriter, error) {

	return &geoParquetResultsWriter{filepath: filepath}, nil
}

func (srw *geoParquetResultsWriter) Write(r consequences.Result) {
	/*
	   //properties

	   var pqWriterProps *parquet.WriterProperties
	   var writerOptions []parquet.WriterProperty
	   writerOptions = append(writerOptions, parquet.WithCompression(compress.Codecs.Zstd))
	   writerOptions = append(writerOptions, parquet.WithMaxRowGroupLength(int64(1000))) //??

	   	if len(writerOptions) > 0 {
	   		pqWriterProps = parquet.NewWriterProperties(writerOptions...)
	   	}

	   builder := pqutil.NewArrowSchemaBuilder()
	   var featureWriter *geoparquet.FeatureWriter

	   	writeBuffered := func() error {
	   		if !builder.Ready() {
	   			return fmt.Errorf("failed to create schema after reading %d features", len(buffer))
	   		}
	   		if err := builder.AddGeometry(geoparquet.DefaultGeometryColumn, geoparquet.DefaultGeometryEncoding); err != nil {
	   			return err
	   		}
	   		sc, scErr := builder.Schema()
	   		if scErr != nil {
	   			return scErr
	   		}
	   		fw, fwErr := geoparquet.NewFeatureWriter(&geoparquet.WriterConfig{
	   			Writer:             output,
	   			ArrowSchema:        sc,
	   			ParquetWriterProps: pqWriterProps,
	   		})
	   		if fwErr != nil {
	   			return fwErr
	   		}

	   		for _, buffered := range buffer {
	   			if err := fw.Write(buffered); err != nil {
	   				return err
	   			}
	   		}
	   		featureWriter = fw
	   		return nil
	   	}

	   sp := "\"properties\": {\""
	   //get the properties from the result
	   x := 0.0
	   y := 0.0
	   result := r.Result

	   	for i, val := range r.Headers {
	   		value, _ := json.Marshal(result[i])
	   		sp += val + "\":" + string(value) + ",\""
	   		if val == "x" {
	   			x = result[i].(float64)
	   		}
	   		if val == "y" {
	   			y = result[i].(float64)
	   		}
	   	}

	   atype := reflect.TypeOf(result[len(result)-1])

	   	if atype.Kind() == reflect.String {
	   		sp = strings.TrimRight(sp, ",\"")
	   		sp += "\""
	   	} else {

	   		sp = strings.TrimRight(sp, ",\"")
	   	}

	   //write out a feature
	   s := "{\"type\": \"Feature\",\n\"geometry\": {\n\"type\": \"Point\",\n\"coordinates\": ["
	   //get the x and y
	   s += fmt.Sprintf("%g, ", x)  //x value
	   s += fmt.Sprintf("%g]\n", y) //y value
	   //close out geometry
	   s += "},\n"
	   s += sp + "}},\n" //this comma might be bad news...

	   //srw.S = s
	*/
}

func (srw *geoParquetResultsWriter) Close() {
	srw.w.Close()
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//largely cloned from planetlabs/gpq as a work around to avoid requirement of pushing to geojson first and then into geoparquet.

type ArrowSchemaBuilder struct {
	fields map[string]*arrow.Field
}

func NewArrowSchemaBuilder() *ArrowSchemaBuilder {
	return &ArrowSchemaBuilder{
		fields: map[string]*arrow.Field{},
	}
}

func (b *ArrowSchemaBuilder) Has(name string) bool {
	_, has := b.fields[name]
	return has
}

func (b *ArrowSchemaBuilder) AddGeometry(name string, encoding string) error {
	var dataType arrow.DataType
	switch encoding {
	case "WKB":
		dataType = arrow.BinaryTypes.Binary
	case "WKT":
		dataType = arrow.BinaryTypes.String
	default:
		return fmt.Errorf("unsupported geometry encoding: %s", encoding)
	}
	b.fields[name] = &arrow.Field{Name: name, Type: dataType, Nullable: true}
	return nil
}

func (b *ArrowSchemaBuilder) Add(record map[string]any) error {
	for name, value := range record {
		if b.fields[name] != nil {
			continue
		}
		if value == nil {
			b.fields[name] = nil
			continue
		}
		if values, ok := value.([]any); ok {
			if len(values) == 0 {
				b.fields[name] = nil
				continue

			}
		}
		field, err := fieldFromValue(name, value, true)
		if err != nil {
			return fmt.Errorf("error converting value for %s: %w", name, err)
		}
		b.fields[name] = field
	}
	return nil
}

func fieldFromValue(name string, value any, nullable bool) (*arrow.Field, error) {
	switch v := value.(type) {
	case bool:
		return &arrow.Field{Name: name, Type: arrow.FixedWidthTypes.Boolean, Nullable: nullable}, nil
	case int, int64:
		return &arrow.Field{Name: name, Type: arrow.PrimitiveTypes.Int64, Nullable: nullable}, nil
	case int32:
		return &arrow.Field{Name: name, Type: arrow.PrimitiveTypes.Int32, Nullable: nullable}, nil
	case float32:
		return &arrow.Field{Name: name, Type: arrow.PrimitiveTypes.Float32, Nullable: nullable}, nil
	case float64:
		return &arrow.Field{Name: name, Type: arrow.PrimitiveTypes.Float64, Nullable: nullable}, nil
	case []byte:
		return &arrow.Field{Name: name, Type: arrow.BinaryTypes.Binary, Nullable: nullable}, nil
	case string:
		return &arrow.Field{Name: name, Type: arrow.BinaryTypes.String, Nullable: nullable}, nil
	case []any:
		if len(v) == 0 {
			return nil, nil
		}
		if err := assertUniformType(v); err != nil {
			return nil, err
		}
		field, err := fieldFromValue(name, v[0], nullable)
		if err != nil {
			return nil, err
		}
		return &arrow.Field{Name: name, Type: arrow.ListOf(field.Type), Nullable: nullable}, nil
	case map[string]any:
		if len(v) == 0 {
			return nil, nil
		}
		return fieldFromMap(name, v, nullable)
	default:
		return nil, fmt.Errorf("cannot convert value: %v", v)
	}
}

func fieldFromMap(name string, value map[string]any, nullable bool) (*arrow.Field, error) {
	keys := sortedKeys(value)
	length := len(keys)
	fields := make([]arrow.Field, length)
	for i, key := range keys {
		field, err := fieldFromValue(key, value[key], nullable)
		if err != nil {
			return nil, fmt.Errorf("trouble generating schema for field %q: %w", key, err)
		}
		if field == nil {
			return nil, nil
		}
		fields[i] = *field
	}
	return &arrow.Field{Name: name, Type: arrow.StructOf(fields...), Nullable: nullable}, nil
}

func assertUniformType(values []any) error {
	length := len(values)
	if length == 0 {
		return errors.New("cannot determine type from zero length slice")
	}
	mixedTypeErr := errors.New("slices must be of all the same type")
	switch v := values[0].(type) {
	case bool:
		for i := 1; i < length; i += 1 {
			if _, ok := values[i].(bool); !ok {
				return mixedTypeErr
			}
		}
	case float64:
		for i := 1; i < length; i += 1 {
			if _, ok := values[i].(float64); !ok {
				return mixedTypeErr
			}
		}
	case string:
		for i := 1; i < length; i += 1 {
			if _, ok := values[i].(string); !ok {
				return mixedTypeErr
			}
		}
	default:
		t := reflect.TypeOf(v)
		for i := 1; i < length; i += 1 {
			if reflect.TypeOf(values[i]) != t {
				return mixedTypeErr
			}
		}
	}
	return nil
}

func (b *ArrowSchemaBuilder) Ready() bool {
	for _, field := range b.fields {
		if field == nil {
			return false
		}
	}
	return true
}

func (b *ArrowSchemaBuilder) Schema() (*arrow.Schema, error) {
	fields := make([]arrow.Field, len(b.fields))
	for i, name := range sortedKeys(b.fields) {
		field := b.fields[name]
		if field == nil {
			return nil, fmt.Errorf("could not derive type for field: %s", name)
		}
		fields[i] = *field
	}
	return arrow.NewSchema(fields, nil), nil
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i += 1
	}
	sort.Strings(keys)
	return keys
}
