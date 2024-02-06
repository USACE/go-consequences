package resultswriters

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type spatialResultsWriter struct {
	FilePath      string
	LayerName     string
	Layer         *gdal.Layer
	ds            *gdal.DataSource
	FieldsCreated bool
	index         int
}

func InitSpatialResultsWriter_Projected(filepath string, layerName string, driver string, ESPG int) (*spatialResultsWriter, error) {
	driverOut := gdal.OGRDriverByName(driver)
	dsOut, okOut := driverOut.Create(filepath, []string{})
	if !okOut {
		//error out?
		return &spatialResultsWriter{}, errors.New("spatial writer at path " + filepath + " of driver type " + driver + " not created")
	}
	//defer dsOut.Destroy() -> probably should destroy on close?
	//set spatial reference?
	sr := gdal.CreateSpatialReference("")
	sr.FromEPSG(ESPG)
	newLayer := dsOut.CreateLayer(layerName, sr, gdal.GeometryType(gdal.GT_Point), []string{"GEOMETRY_NAME=shape"}) //forcing point data type.  source type (using lyaer.type()) from postgis was a generic geometry

	return &spatialResultsWriter{FilePath: filepath, LayerName: layerName, ds: &dsOut, Layer: &newLayer, index: 0}, nil
}
func InitSpatialResultsWriter(filepath string, layerName string, driver string) (*spatialResultsWriter, error) {
	return InitSpatialResultsWriter_Projected(filepath, layerName, driver, 4326)
}
func (srw *spatialResultsWriter) Write(r consequences.Result) {
	//if header has not been built:
	result := r.Result
	if !srw.FieldsCreated {
		func() {
			fieldDef := gdal.CreateFieldDefinition("objectid", gdal.FieldType(gdal.FT_Integer))
			defer fieldDef.Destroy()
			srw.Layer.CreateField(fieldDef, true)
		}()
		for i, val := range r.Headers {
			//need to identify value type
			func() {
				if val == "hazard" { //not a huge fan of this, because it is specific to that kind of hazard.
					fieldDef := gdal.CreateFieldDefinition("depth", gdal.FieldType(gdal.FT_Real))
					defer fieldDef.Destroy()
					srw.Layer.CreateField(fieldDef, true) //approxOk.
				} else {
					atype := reflect.TypeOf(result[i]) //.Elem()
					gotype := atype.Kind()
					fieldName := val
					if len(val) > 10 {
						fieldName = val[0:10]
						fieldName = strings.TrimSpace(fieldName)
					}
					gdaltype := gdalTypes[gotype]
					fieldDef := gdal.CreateFieldDefinition(fieldName, gdaltype)
					defer fieldDef.Destroy()
					srw.Layer.CreateField(fieldDef, true) //approxOk.
				}
			}()
		}
		srw.FieldsCreated = true
		srw.Layer.StartTransaction()
	}

	//add a feature to a layer?
	layerDef := srw.Layer.Definition()
	//if header has been built, add the feature, and the attributes.

	feature := layerDef.Create()
	//defer feature.Destroy()
	feature.SetFieldInteger(0, srw.index)
	//create a point geometry - not sure the best way to do that.
	x := 0.0
	y := 0.0
	g := gdal.Create(gdal.GeometryType(gdal.GT_Point))
	defer g.Destroy()
	for i, val := range r.Headers {
		if val == "x" {
			x = result[i].(float64)
		}
		if val == "y" {
			y = result[i].(float64)
		}
		fieldName := val
		if len(val) > 10 {
			fieldName = val[0:10]
			fieldName = strings.TrimSpace(fieldName)
		}
		value := result[i]
		att := reflect.TypeOf(result[i])
		valType := att.Kind()
		if val == "hazard" {
			fieldName = "depth"
			de, dok := value.(hazards.HazardEvent)
			if dok {
				valType = reflect.Float64
				if de.Has(hazards.Depth) {
					fieldName = "depth"
					value = de.Depth()
				}
			} else {
				//must be an array - bummer.
				//get at the elements of the slice, add all depths to the table?
				fieldName = "multidepths"
				valType = reflect.Float64
				value = 123.456
			}

		}
		if val == "hazards" {
			fieldName = "hazards"
			de, dok := value.(string)
			if dok {
				valType = reflect.String
				value = de
			} else {
				//must be an array - bummer.
				//get at the elements of the slice, add all depths to the table?
				fieldName = "multidepths"
				valType = reflect.Float64
				value = 123.456
			}

		}
		idx := layerDef.FieldIndex(fieldName)
		switch valType {
		case reflect.String:
			feature.SetFieldString(idx, value.(string))
		case reflect.Float32:
			gval := float64(value.(float32))
			feature.SetFieldFloat64(idx, gval)
		case reflect.Float64:
			gval := value.(float64)
			feature.SetFieldFloat64(idx, gval)
		case reflect.Int32:
			gval := int(value.(int32))
			feature.SetFieldInteger(idx, gval)
		case reflect.Uint8:
			gval := int(value.(uint8))
			feature.SetFieldInteger(idx, gval)
		}

	}
	g.SetPoint(0, x, y, 0)
	feature.SetGeometryDirectly(g)
	err := srw.Layer.Create(feature)
	if err != nil {
		fmt.Println(err)
	}
	if srw.index%100000 == 0 {
		err2 := srw.Layer.CommitTransaction()
		if err2 != nil {
			fmt.Println(err2)
		}
		srw.Layer.StartTransaction()
	}

	srw.index++ //incriment.
	//feature.Destroy() //testing an explicit call.//causes seg fault error, probably not calling causes a memory leak... oy vey.
}
func (srw *spatialResultsWriter) Close() {
	//not sure what this should do - Destroy should close resource connections.
	err2 := srw.Layer.CommitTransaction()
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Printf("Closing, wrote %v features\n", srw.index)
	srw.ds.Destroy()
}

/*
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/USACE/go-consequences/consequences"
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet"
	"github.com/apache/arrow/go/v14/parquet/compress"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/metadata"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/apache/arrow/go/v14/parquet/schema"
	pqschema "github.com/apache/arrow/go/v14/parquet/schema"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
	orbjson "github.com/paulmach/orb/geojson"
)

type geoParquetResultsWriter struct {
	filepath string
	w        file.Writer
}

func InitGeoParquetResultsWriterFromFile(filepath string) (*geoParquetResultsWriter, error) {

	return &geoParquetResultsWriter{filepath: filepath}, nil
}

func (srw *geoParquetResultsWriter) Write(r consequences.Result) {

	   //properties
	if convertOptions == nil {
		convertOptions = defaultOptions
	}
	reader := NewFeatureReader(input)
	buffer := []*Feature{}
	builder := NewArrowSchemaBuilder()
	featuresRead := 0

	var pqWriterProps *parquet.WriterProperties
	var writerOptions []parquet.WriterProperty
	if convertOptions.Compression != "" {
		compression, err := pqutil.GetCompression(convertOptions.Compression)
		if err != nil {
			return err
		}
		writerOptions = append(writerOptions, parquet.WithCompression(compression))
	}
	if convertOptions.RowGroupLength > 0 {
		writerOptions = append(writerOptions, parquet.WithMaxRowGroupLength(int64(convertOptions.RowGroupLength)))
	}
	if len(writerOptions) > 0 {
		pqWriterProps = parquet.NewWriterProperties(writerOptions...)
	}

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

	for {
		feature, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		featuresRead += 1
		if featureWriter == nil {
			if err := builder.Add(feature.Properties); err != nil {
				return err
			}

			if !builder.Ready() {
				buffer = append(buffer, feature)
				if len(buffer) > convertOptions.MaxFeatures {
					return fmt.Errorf("failed to create parquet schema after reading %d features", convertOptions.MaxFeatures)
				}
				continue
			}

			if len(buffer) < convertOptions.MinFeatures-1 {
				buffer = append(buffer, feature)
				continue
			}

			if err := writeBuffered(); err != nil {
				panic(err)
			}
		}
		if err := featureWriter.Write(feature); err != nil {
			panic(err)
		}
	}
	if featuresRead > 0 {
		if featureWriter == nil {
			if err := writeBuffered(); err != nil {
				panic(err)
			}
		}
		return featureWriter.Close()
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

var ParquetStringType = pqschema.StringLogicalType{}

func LookupNode(schema *pqschema.Schema, name string) (pqschema.Node, bool) {
	root := schema.Root()
	index := root.FieldIndexByName(name)
	if index < 0 {
		return nil, false
	}

	return root.Field(index), true
}

func LookupPrimitiveNode(schema *pqschema.Schema, name string) (*pqschema.PrimitiveNode, bool) {
	node, ok := LookupNode(schema, name)
	if !ok {
		return nil, false
	}

	primitive, ok := node.(*pqschema.PrimitiveNode)
	return primitive, ok
}

func LookupGroupNode(schema *pqschema.Schema, name string) (*pqschema.GroupNode, bool) {
	node, ok := LookupNode(schema, name)
	if !ok {
		return nil, false
	}

	group, ok := node.(*pqschema.GroupNode)
	return group, ok
}

func LookupListElementNode(sc *pqschema.Schema, name string) (*pqschema.PrimitiveNode, bool) {
	node, ok := LookupGroupNode(sc, name)
	if !ok {
		return nil, false
	}

	if node.NumFields() != 1 {
		return nil, false
	}

	group, ok := node.Field(0).(*pqschema.GroupNode)
	if !ok {
		return nil, false
	}

	if group.NumFields() != 1 {
		return nil, false
	}

	element, ok := group.Field(0).(*pqschema.PrimitiveNode)
	return element, ok
}

// ParquetSchemaString generates a string representation of the schema as documented
// in https://pkg.go.dev/github.com/fraugster/parquet-go/parquetschema
func ParquetSchemaString(schema *pqschema.Schema) string {
	w := &parquetWriter{}
	return w.String(schema)
}

type parquetWriter struct {
	builder *strings.Builder
	err     error
}

func (w *parquetWriter) String(schema *pqschema.Schema) string {
	w.builder = &strings.Builder{}
	w.err = nil
	w.writeSchema(schema)
	if w.err != nil {
		return w.err.Error()
	}
	return w.builder.String()
}

func (w *parquetWriter) writeLine(str string, level int) {
	if w.err != nil {
		return
	}
	indent := strings.Repeat("  ", level)
	if _, err := w.builder.WriteString(indent + str + "\n"); err != nil {
		w.err = err
	}
}

func (w *parquetWriter) writeSchema(schema *pqschema.Schema) {
	w.writeLine("message {", 0)
	root := schema.Root()
	for i := 0; i < root.NumFields(); i += 1 {
		w.writeNode(root.Field(i), 1)
	}
	w.writeLine("}", 0)
}

func (w *parquetWriter) writeNode(node pqschema.Node, level int) {
	switch n := node.(type) {
	case *pqschema.GroupNode:
		w.writeGroupNode(n, level)
	case *pqschema.PrimitiveNode:
		w.writePrimitiveNode(n, level)
	default:
		w.writeLine(fmt.Sprintf("unknown node type: %v", node), level)
	}
}

func (w *parquetWriter) writeGroupNode(node *pqschema.GroupNode, level int) {
	repetition := node.RepetitionType().String()
	name := node.Name()
	annotation := LogicalOrConvertedAnnotation(node)

	w.writeLine(fmt.Sprintf("%s group %s%s {", repetition, name, annotation), level)
	for i := 0; i < node.NumFields(); i += 1 {
		w.writeNode(node.Field(i), level+1)
	}
	w.writeLine("}", level)
}

func (w *parquetWriter) writePrimitiveNode(node *pqschema.PrimitiveNode, level int) {
	repetition := node.RepetitionType().String()
	name := node.Name()
	nodeType := physicalTypeString(node.PhysicalType())
	annotation := LogicalOrConvertedAnnotation(node)

	w.writeLine(fmt.Sprintf("%s %s %s%s;", repetition, nodeType, name, annotation), level)
}

func LogicalOrConvertedAnnotation(node pqschema.Node) string {
	logicalType := node.LogicalType()
	convertedType := node.ConvertedType()

	switch t := logicalType.(type) {
	case *pqschema.IntLogicalType:
		return fmt.Sprintf(" (INT (%d, %t))", t.BitWidth(), t.IsSigned())
	case *pqschema.DecimalLogicalType:
		return fmt.Sprintf(" (DECIMAL (%d, %d))", t.Precision(), t.Scale())
	case *pqschema.TimestampLogicalType:
		var unit string
		switch t.TimeUnit() {
		case pqschema.TimeUnitMillis:
			unit = "MILLIS"
		case pqschema.TimeUnitMicros:
			unit = "MICROS"
		case pqschema.TimeUnitNanos:
			unit = "NANOS"
		default:
			unit = "UNKNOWN"
		}
		return fmt.Sprintf(" (TIMESTAMP (%s, %t))", unit, t.IsAdjustedToUTC())
	}

	var annotation string
	_, invalid := logicalType.(pqschema.UnknownLogicalType)
	_, none := logicalType.(pqschema.NoLogicalType)

	if logicalType != nil && !invalid && !none {
		annotation = fmt.Sprintf(" (%s)", strings.ToUpper(logicalType.String()))
	} else if convertedType != pqschema.ConvertedTypes.None {
		annotation = fmt.Sprintf(" (%s)", strings.ToUpper(convertedType.String()))
	}

	return annotation
}

var physicalTypeLookup = map[string]string{
	"byte_array": "binary",
}

func physicalTypeString(physical parquet.Type) string {
	nodeType := strings.ToLower(physical.String())
	if altType, ok := physicalTypeLookup[nodeType]; ok {
		return altType
	}
	if physical == parquet.Types.FixedLenByteArray {
		nodeType += fmt.Sprintf(" (%d)", physical.ByteSize())
	}
	return nodeType
}


type ConvertOptions struct {
	InputPrimaryColumn string
	Compression        string
	RowGroupLength     int
}

func getMetadata(fileReader *file.Reader, convertOptions *ConvertOptions) *Metadata {
	metadata, err := GetMetadata(fileReader.MetaData().KeyValueMetadata())
	if err != nil {
		primaryColumn := DefaultGeometryColumn
		if convertOptions.InputPrimaryColumn != "" {
			primaryColumn = convertOptions.InputPrimaryColumn
		}
		metadata = &Metadata{
			Version:       Version,
			PrimaryColumn: primaryColumn,
			Columns: map[string]*GeometryColumn{
				primaryColumn: getDefaultGeometryColumn(),
			},
		}
	}
	if convertOptions.InputPrimaryColumn != "" && metadata.PrimaryColumn != convertOptions.InputPrimaryColumn {
		metadata.PrimaryColumn = convertOptions.InputPrimaryColumn
	}
	return metadata
}

func FromParquet(input parquet.ReaderAtSeeker, output io.Writer, convertOptions *ConvertOptions) error {
	if convertOptions == nil {
		convertOptions = &ConvertOptions{}
	}

	var compression *compress.Compression
	if convertOptions.Compression != "" {
		c, err := GetCompression(convertOptions.Compression)
		if err != nil {
			return err
		}
		compression = &c
	}

	datasetInfo := NewDatasetStats(true)
	transformSchema := func(fileReader *file.Reader) (*schema.Schema, error) {
		inputSchema := fileReader.MetaData().Schema
		inputRoot := inputSchema.Root()
		metadata := getMetadata(fileReader, convertOptions)
		for geomColName := range metadata.Columns {
			if inputRoot.FieldIndexByName(geomColName) < 0 {
				message := fmt.Sprintf(
					"expected a geometry column named %q,"+
						" use the --input-primary-column to supply a different primary geometry",
					geomColName,
				)
				return nil, errors.New(message)
			}
		}
		for fieldNum := 0; fieldNum < inputRoot.NumFields(); fieldNum += 1 {
			field := inputRoot.Field(fieldNum)
			name := field.Name()
			if _, ok := metadata.Columns[name]; !ok {
				continue
			}
			if field.LogicalType() == ParquetStringType {
				datasetInfo.AddCollection(name)
			}
		}

		if datasetInfo.NumCollections() == 0 {
			return inputSchema, nil
		}

		numFields := inputRoot.NumFields()
		fields := make([]schema.Node, numFields)
		for fieldNum := 0; fieldNum < numFields; fieldNum += 1 {
			inputField := inputRoot.Field(fieldNum)
			if !datasetInfo.HasCollection(inputField.Name()) {
				fields[fieldNum] = inputField
				continue
			}
			outputField, err := schema.NewPrimitiveNode(inputField.Name(), inputField.RepetitionType(), parquet.Types.ByteArray, -1, -1)
			if err != nil {
				return nil, err
			}
			fields[fieldNum] = outputField
		}

		outputRoot, err := schema.NewGroupNode(inputRoot.Name(), inputRoot.RepetitionType(), fields, -1)
		if err != nil {
			return nil, err
		}
		return schema.NewSchema(outputRoot), nil
	}

	transformColumn := func(inputField *arrow.Field, outputField *arrow.Field, chunked *arrow.Chunked) (*arrow.Chunked, error) {
		if !datasetInfo.HasCollection(inputField.Name) {
			return chunked, nil
		}
		chunks := chunked.Chunks()
		transformed := make([]arrow.Array, len(chunks))
		builder := array.NewBinaryBuilder(memory.DefaultAllocator, arrow.BinaryTypes.Binary)
		defer builder.Release()

		collectionInfo := NewGeometryStats(false)
		for i, arr := range chunks {
			stringArray, ok := arr.(*array.String)
			if !ok {
				return nil, fmt.Errorf("expected a string array for %q, got %v", inputField.Name, arr)
			}
			for rowNum := 0; rowNum < stringArray.Len(); rowNum += 1 {
				if outputField.Nullable && stringArray.IsNull(rowNum) {
					builder.AppendNull()
					continue
				}
				str := stringArray.Value(rowNum)
				geometry, wktErr := wkt.Unmarshal(str)
				if wktErr != nil {
					return nil, wktErr
				}
				value, wkbErr := wkb.Marshal(geometry)
				if wkbErr != nil {
					return nil, wkbErr
				}
				collectionInfo.AddType(geometry.GeoJSONType())
				bounds := geometry.Bound()
				collectionInfo.AddBounds(&bounds)
				builder.Append(value)
			}
			transformed[i] = builder.NewArray()
		}
		datasetInfo.AddBounds(inputField.Name, collectionInfo.Bounds())
		datasetInfo.AddTypes(inputField.Name, collectionInfo.Types())
		chunked.Release()
		return arrow.NewChunked(builder.Type(), transformed), nil
	}

	beforeClose := func(fileReader *file.Reader, fileWriter *pqarrow.FileWriter) error {
		metadata := getMetadata(fileReader, convertOptions)
		for name, geometryCol := range metadata.Columns {
			if !datasetInfo.HasCollection(name) {
				continue
			}
			bounds := datasetInfo.Bounds(name)
			geometryCol.Bounds = []float64{
				bounds.Left(), bounds.Bottom(), bounds.Right(), bounds.Top(),
			}
			geometryCol.GeometryTypes = datasetInfo.Types(name)
		}
		encodedMetadata, jsonErr := json.Marshal(metadata)
		if jsonErr != nil {
			return fmt.Errorf("trouble encoding %q metadata: %w", MetadataKey, jsonErr)
		}
		if err := fileWriter.AppendKeyValueMetadata(MetadataKey, string(encodedMetadata)); err != nil {
			return fmt.Errorf("trouble appending %q metadata: %w", MetadataKey, err)
		}
		return nil
	}

	config := &TransformConfig{
		Reader:          input,
		Writer:          output,
		TransformSchema: transformSchema,
		TransformColumn: transformColumn,
		BeforeClose:     beforeClose,
		Compression:     compression,
		RowGroupLength:  convertOptions.RowGroupLength,
	}

	return TransformByColumn(config)
}
const (
	Version                     = "1.0.0"
	MetadataKey                 = "geo"
	EdgesPlanar                 = "planar"
	EdgesSpherical              = "spherical"
	OrientationCounterClockwise = "counterclockwise"
	DefaultGeometryColumn       = "geometry"
	DefaultGeometryEncoding     = EncodingWKB
)

var GeometryTypes = []string{
	"Point",
	"LineString",
	"Polygon",
	"MultiPoint",
	"MultiLineString",
	"MultiPolygon",
	"GeometryCollection",
	"Point Z",
	"LineString Z",
	"Polygon Z",
	"MultiPoint Z",
	"MultiLineString Z",
	"MultiPolygon Z",
	"GeometryCollection Z",
}

type Metadata struct {
	Version       string                     `json:"version"`
	PrimaryColumn string                     `json:"primary_column"`
	Columns       map[string]*GeometryColumn `json:"columns"`
}

func (m *Metadata) Clone() *Metadata {
	clone := &Metadata{}
	*clone = *m
	clone.Columns = make(map[string]*GeometryColumn, len(m.Columns))
	for i, v := range m.Columns {
		clone.Columns[i] = v.clone()
	}
	return clone
}

type ProjId struct {
	Authority string `json:"authority"`
	Code      any    `json:"code"`
}

type Proj struct {
	Name string  `json:"name"`
	Id   *ProjId `json:"id"`
}

func (p *Proj) String() string {
	id := ""
	if p.Id != nil {
		if code, ok := p.Id.Code.(string); ok {
			id = p.Id.Authority + ":" + code
		} else if code, ok := p.Id.Code.(float64); ok {
			id = fmt.Sprintf("%s:%g", p.Id.Authority, code)
		}
	}
	if p.Name != "" {
		return p.Name
	}
	if id == "" {
		return "Unknown"
	}
	return id
}

type GeometryColumn struct {
	Encoding      string    `json:"encoding"`
	GeometryType  any       `json:"geometry_type,omitempty"`
	GeometryTypes any       `json:"geometry_types"`
	CRS           *Proj     `json:"crs,omitempty"`
	Edges         string    `json:"edges,omitempty"`
	Orientation   string    `json:"orientation,omitempty"`
	Bounds        []float64 `json:"bbox,omitempty"`
	Epoch         float64   `json:"epoch,omitempty"`
}

func (g *GeometryColumn) clone() *GeometryColumn {
	clone := &GeometryColumn{}
	*clone = *g
	clone.Bounds = make([]float64, len(g.Bounds))
	copy(clone.Bounds, g.Bounds)
	return clone
}

func (col *GeometryColumn) GetGeometryTypes() []string {
	if multiType, ok := col.GeometryTypes.([]any); ok {
		types := make([]string, len(multiType))
		for i, value := range multiType {
			geometryType, ok := value.(string)
			if !ok {
				return nil
			}
			types[i] = geometryType
		}
		return types
	}

	if singleType, ok := col.GeometryType.(string); ok {
		return []string{singleType}
	}

	values, ok := col.GeometryType.([]any)
	if !ok {
		return nil
	}

	types := make([]string, len(values))
	for i, value := range values {
		geometryType, ok := value.(string)
		if !ok {
			return nil
		}
		types[i] = geometryType
	}

	return types
}

func getDefaultGeometryColumn() *GeometryColumn {
	return &GeometryColumn{
		Encoding:      DefaultGeometryEncoding,
		GeometryTypes: []string{},
	}
}

func DefaultMetadata() *Metadata {
	return &Metadata{
		Version:       Version,
		PrimaryColumn: DefaultGeometryColumn,
		Columns: map[string]*GeometryColumn{
			DefaultGeometryColumn: getDefaultGeometryColumn(),
		},
	}
}

var ErrNoMetadata = fmt.Errorf("missing %s metadata key", MetadataKey)
var ErrDuplicateMetadata = fmt.Errorf("found more than one %s metadata key", MetadataKey)

func GetMetadata(keyValueMetadata metadata.KeyValueMetadata) (*Metadata, error) {
	value, err := GetMetadataValue(keyValueMetadata)
	if err != nil {
		return nil, err
	}
	geoFileMetadata := &Metadata{}
	jsonErr := json.Unmarshal([]byte(value), geoFileMetadata)
	if jsonErr != nil {
		return nil, fmt.Errorf("unable to parse %s metadata: %w", MetadataKey, jsonErr)
	}
	return geoFileMetadata, nil
}

func GetMetadataValue(keyValueMetadata metadata.KeyValueMetadata) (string, error) {
	var value *string
	for _, kv := range keyValueMetadata {
		if kv.Key == MetadataKey {
			if value != nil {
				return "", ErrDuplicateMetadata
			}
			value = kv.Value
		}
	}
	if value == nil {
		return "", ErrNoMetadata
	}
	return *value, nil
}

type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

var (
	_ json.Marshaler = (*FeatureCollection)(nil)
)

func (c *FeatureCollection) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     "FeatureCollection",
		"features": c.Features,
	}
	return json.Marshal(m)
}

type Feature struct {
	Id         any            `json:"id,omitempty"`
	Type       string         `json:"type"`
	Geometry   orb.Geometry   `json:"geometry"`
	Properties map[string]any `json:"properties"`
}

var (
	_ json.Marshaler   = (*Feature)(nil)
	_ json.Unmarshaler = (*Feature)(nil)
)

func (f *Feature) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":       "Feature",
		"geometry":   orbjson.NewGeometry(f.Geometry),
		"properties": f.Properties,
	}
	if f.Id != nil {
		m["id"] = f.Id
	}
	return json.Marshal(m)
}

type jsonFeature struct {
	Id         any             `json:"id,omitempty"`
	Type       string          `json:"type"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties map[string]any  `json:"properties"`
}

var rawNull = json.RawMessage([]byte("null"))

func isRawNull(raw json.RawMessage) bool {
	if len(raw) != len(rawNull) {
		return false
	}
	for i, c := range raw {
		if c != rawNull[i] {
			return false
		}
	}
	return true
}

func (f *Feature) UnmarshalJSON(data []byte) error {
	jf := &jsonFeature{}
	if err := json.Unmarshal(data, jf); err != nil {
		return err
	}

	f.Type = jf.Type
	f.Id = jf.Id
	f.Properties = jf.Properties

	if isRawNull(jf.Geometry) {
		return nil
	}
	geometry := &orbjson.Geometry{}
	if err := json.Unmarshal(jf.Geometry, geometry); err != nil {
		return err
	}

	f.Geometry = geometry.Geometry()
	return nil
}

const (
	EncodingWKB = "WKB"
	EncodingWKT = "WKT"
)

func DecodeGeometry(value any, encoding string) (*orbjson.Geometry, error) {
	if value == nil {
		return nil, nil
	}
	if encoding == "" {
		if _, ok := value.([]byte); ok {
			encoding = EncodingWKB
		} else if _, ok := value.(string); ok {
			encoding = EncodingWKT
		}
	}
	if encoding == EncodingWKB {
		data, ok := value.([]byte)
		if !ok {
			return nil, fmt.Errorf("expected bytes for wkb geometry, got %T", value)
		}
		g, err := wkb.Unmarshal(data)
		if err != nil {
			return nil, err
		}
		return orbjson.NewGeometry(g), nil
	}
	if encoding == EncodingWKT {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for wkt geometry, got %T", value)
		}
		g, err := wkt.Unmarshal(str)
		if err != nil {
			return nil, err
		}
		return orbjson.NewGeometry(g), nil
	}
	return nil, fmt.Errorf("unsupported encoding: %s", encoding)
}

type GeometryStats struct {
	mutex *sync.RWMutex
	minX  float64
	maxX  float64
	minY  float64
	maxY  float64
	types map[string]bool
}

func NewGeometryStats(concurrent bool) *GeometryStats {
	var mutex *sync.RWMutex
	if concurrent {
		mutex = &sync.RWMutex{}
	}
	return &GeometryStats{
		mutex: mutex,
		types: map[string]bool{},
		minX:  math.MaxFloat64,
		maxX:  -math.MaxFloat64,
		minY:  math.MaxFloat64,
		maxY:  -math.MaxFloat64,
	}
}

func (i *GeometryStats) writeLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Lock()
}

func (i *GeometryStats) writeUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Unlock()
}

func (i *GeometryStats) readLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RLock()
}

func (i *GeometryStats) readUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RUnlock()
}

func (i *GeometryStats) AddBounds(bounds *orb.Bound) {
	i.writeLock()
	minPoint := bounds.Min
	minX := minPoint[0]
	minY := minPoint[1]
	maxPoint := bounds.Max
	maxX := maxPoint[0]
	maxY := maxPoint[1]
	i.minX = math.Min(i.minX, minX)
	i.maxX = math.Max(i.maxX, maxX)
	i.minY = math.Min(i.minY, minY)
	i.maxY = math.Max(i.maxY, maxY)
	i.writeUnlock()
}

func (i *GeometryStats) Bounds() *orb.Bound {
	i.readLock()
	bounds := &orb.Bound{
		Min: orb.Point{i.minX, i.minY},
		Max: orb.Point{i.maxX, i.maxY},
	}
	i.readUnlock()
	return bounds
}

func (i *GeometryStats) AddType(typ string) {
	i.writeLock()
	i.types[typ] = true
	i.writeUnlock()
}

func (i *GeometryStats) AddTypes(types []string) {
	i.writeLock()
	for _, typ := range types {
		i.types[typ] = true
	}
	i.writeUnlock()
}

func (i *GeometryStats) Types() []string {
	i.readLock()
	types := []string{}
	for typ, ok := range i.types {
		if ok {
			types = append(types, typ)
		}
	}
	i.readUnlock()
	return types
}

type DatasetStats struct {
	mutex       *sync.RWMutex
	collections map[string]*GeometryStats
}

func NewDatasetStats(concurrent bool) *DatasetStats {
	var mutex *sync.RWMutex
	if concurrent {
		mutex = &sync.RWMutex{}
	}
	return &DatasetStats{
		mutex:       mutex,
		collections: map[string]*GeometryStats{},
	}
}

func (i *DatasetStats) writeLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Lock()
}

func (i *DatasetStats) writeUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Unlock()
}

func (i *DatasetStats) readLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RLock()
}

func (i *DatasetStats) readUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RUnlock()
}

func (i *DatasetStats) NumCollections() int {
	i.readLock()
	num := len(i.collections)
	i.readUnlock()
	return num
}

func (i *DatasetStats) AddCollection(name string) {
	i.writeLock()
	i.collections[name] = NewGeometryStats(i.mutex != nil)
	i.writeUnlock()
}

func (i *DatasetStats) HasCollection(name string) bool {
	i.readLock()
	_, has := i.collections[name]
	i.readUnlock()
	return has
}

func (i *DatasetStats) AddBounds(name string, bounds *orb.Bound) {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	collection.AddBounds(bounds)
}

func (i *DatasetStats) Bounds(name string) *orb.Bound {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	return collection.Bounds()
}

func (i *DatasetStats) AddTypes(name string, types []string) {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	collection.AddTypes(types)
}

func (i *DatasetStats) Types(name string) []string {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	return collection.Types()
}
func GetCompression(codec string) (compress.Compression, error) {
	switch codec {
	case "uncompressed":
		return compress.Codecs.Uncompressed, nil
	case "snappy":
		return compress.Codecs.Snappy, nil
	case "gzip":
		return compress.Codecs.Gzip, nil
	case "brotli":
		return compress.Codecs.Brotli, nil
	case "zstd":
		return compress.Codecs.Zstd, nil
	case "lz4":
		return compress.Codecs.Lz4, nil
	default:
		return compress.Codecs.Uncompressed, fmt.Errorf("invalid compression codec %s", codec)
	}
}
type ColumnTransformer func(*arrow.Field, *arrow.Field, *arrow.Chunked) (*arrow.Chunked, error)

type SchemaTransformer func(*file.Reader) (*schema.Schema, error)

type TransformConfig struct {
	Reader          parquet.ReaderAtSeeker
	Writer          io.Writer
	Compression     *compress.Compression
	RowGroupLength  int
	TransformSchema SchemaTransformer
	TransformColumn ColumnTransformer
	BeforeClose     func(*file.Reader, *pqarrow.FileWriter) error
}

func getWriterProperties(config *TransformConfig, fileReader *file.Reader) (*parquet.WriterProperties, error) {
	var writerProperties []parquet.WriterProperty
	if config.Compression != nil {
		writerProperties = append(writerProperties, parquet.WithCompression(*config.Compression))
	} else {
		// retain existing column compression (from the first row group)
		if fileReader.NumRowGroups() > 0 {
			rowGroupMetadata := fileReader.RowGroup(0).MetaData()
			for colNum := 0; colNum < rowGroupMetadata.NumColumns(); colNum += 1 {
				colChunkMetadata, err := rowGroupMetadata.ColumnChunk(colNum)
				if err != nil {
					return nil, fmt.Errorf("failed to get column chunk metadata for column %d", colNum)
				}
				compression := colChunkMetadata.Compression()
				if compression != compress.Codecs.Uncompressed {
					colPath := colChunkMetadata.PathInSchema()
					writerProperties = append(writerProperties, parquet.WithCompressionPath(colPath, compression))
				}
			}
		}
	}

	if config.RowGroupLength > 0 {
		writerProperties = append(writerProperties, parquet.WithMaxRowGroupLength(int64(config.RowGroupLength)))
	}

	return parquet.NewWriterProperties(writerProperties...), nil
}

func TransformByColumn(config *TransformConfig) error {
	if config.Reader == nil {
		return errors.New("reader is required")
	}
	if config.Writer == nil {
		return errors.New("writer is required")
	}

	fileReader, fileReaderErr := file.NewParquetReader(config.Reader)
	if fileReaderErr != nil {
		return fileReaderErr
	}
	defer fileReader.Close()

	outputSchema := fileReader.MetaData().Schema
	if config.TransformSchema != nil {
		schema, err := config.TransformSchema(fileReader)
		if err != nil {
			return err
		}
		outputSchema = schema
	}

	arrowReadProperties := pqarrow.ArrowReadProperties{}

	arrowReader, arrowError := pqarrow.NewFileReader(fileReader, arrowReadProperties, memory.DefaultAllocator)
	if arrowError != nil {
		return arrowError
	}
	inputManifest := arrowReader.Manifest

	outputManifest, manifestErr := pqarrow.NewSchemaManifest(outputSchema, fileReader.MetaData().KeyValueMetadata(), &arrowReadProperties)
	if manifestErr != nil {
		return manifestErr
	}

	numFields := len(outputManifest.Fields)
	if numFields != len(inputManifest.Fields) {
		return fmt.Errorf("unexpected number of fields in the output schema, got %d, expected %d", numFields, len(inputManifest.Fields))
	}

	writerProperties, propErr := getWriterProperties(config, fileReader)
	if propErr != nil {
		return propErr
	}

	arrowSchema, arrowSchemaErr := pqarrow.FromParquet(outputSchema, &arrowReadProperties, fileReader.MetaData().KeyValueMetadata())
	if arrowSchemaErr != nil {
		return arrowSchemaErr
	}

	fileWriter, fileWriterErr := pqarrow.NewFileWriter(arrowSchema, config.Writer, writerProperties, pqarrow.DefaultWriterProps())
	if fileWriterErr != nil {
		return fileWriterErr
	}

	ctx := pqarrow.NewArrowWriteContext(context.Background(), nil)

	if config.RowGroupLength > 0 {
		columnReaders := make([]*pqarrow.ColumnReader, numFields)
		for fieldNum := 0; fieldNum < numFields; fieldNum += 1 {
			colReader, err := arrowReader.GetColumn(ctx, fieldNum)
			if err != nil {
				return err
			}
			columnReaders[fieldNum] = colReader
		}

		numRows := fileReader.NumRows()
		numRowsWritten := int64(0)
		for {
			fileWriter.NewRowGroup()
			numRowsInGroup := 0
			for fieldNum := 0; fieldNum < numFields; fieldNum += 1 {
				colReader := columnReaders[fieldNum]
				arr, readErr := colReader.NextBatch(int64(config.RowGroupLength))
				if readErr != nil {
					return readErr
				}
				if config.TransformColumn != nil {
					inputField := inputManifest.Fields[fieldNum].Field
					outputField := outputManifest.Fields[fieldNum].Field
					transformed, err := config.TransformColumn(inputField, outputField, arr)
					if err != nil {
						return err
					}
					if transformed.DataType() != outputField.Type {
						return fmt.Errorf("transform generated an unexpected type, got %s, expected %s", transformed.DataType().Name(), outputField.Type.Name())
					}
					arr = transformed
				}
				if numRowsInGroup == 0 {
					// TODO: propose fileWriter.RowGroupNumRows()
					numRowsInGroup = arr.Len()
				}
				if err := fileWriter.WriteColumnChunked(arr, 0, int64(arr.Len())); err != nil {
					return err
				}
			}
			numRowsWritten += int64(numRowsInGroup)
			if numRowsWritten >= numRows {
				break
			}
		}
	} else {
		numRowGroups := fileReader.NumRowGroups()
		for rowGroupIndex := 0; rowGroupIndex < numRowGroups; rowGroupIndex += 1 {
			rowGroupReader := arrowReader.RowGroup(rowGroupIndex)
			fileWriter.NewRowGroup()
			for fieldNum := 0; fieldNum < numFields; fieldNum += 1 {
				arr, readErr := rowGroupReader.Column(fieldNum).Read(ctx)
				if readErr != nil {
					return readErr
				}
				if config.TransformColumn != nil {
					inputField := inputManifest.Fields[fieldNum].Field
					outputField := outputManifest.Fields[fieldNum].Field
					transformed, err := config.TransformColumn(inputField, outputField, arr)
					if err != nil {
						return err
					}
					arr = transformed
				}
				if err := fileWriter.WriteColumnChunked(arr, 0, int64(arr.Len())); err != nil {
					return err
				}
			}
		}
	}

	if config.BeforeClose != nil {
		if err := config.BeforeClose(fileReader, fileWriter); err != nil {
			return err
		}
	}
	return fileWriter.Close()
}
type FeatureWriter struct {
	geoMetadata        *Metadata
	maxRowGroupLength  int64
	bufferedLength     int64
	fileWriter         *pqarrow.FileWriter
	recordBuilder      *array.RecordBuilder
	geometryTypeLookup map[string]map[string]bool
	boundsLookup       map[string]*orb.Bound
}
type WriterConfig struct {
	Writer             io.Writer
	Metadata           *Metadata
	ParquetWriterProps *parquet.WriterProperties
	ArrowWriterProps   *pqarrow.ArrowWriterProperties
	ArrowSchema        *arrow.Schema
}
func NewFeatureWriter(config *WriterConfig) (*FeatureWriter, error) {
	parquetProps := config.ParquetWriterProps
	if parquetProps == nil {
		parquetProps = parquet.NewWriterProperties()
	}

	arrowProps := config.ArrowWriterProps
	if arrowProps == nil {
		defaults := pqarrow.DefaultWriterProps()
		arrowProps = &defaults
	}

	geoMetadata := config.Metadata
	if geoMetadata == nil {
		geoMetadata = DefaultMetadata()
	}

	if config.ArrowSchema == nil {
		return nil, errors.New("schema is required")
	}

	if config.Writer == nil {
		return nil, errors.New("writer is required")
	}
	fileWriter, fileErr := pqarrow.NewFileWriter(config.ArrowSchema, config.Writer, parquetProps, *arrowProps)
	if fileErr != nil {
		return nil, fileErr
	}

	writer := &FeatureWriter{
		geoMetadata:        geoMetadata,
		fileWriter:         fileWriter,
		maxRowGroupLength:  parquetProps.MaxRowGroupLength(),
		bufferedLength:     0,
		recordBuilder:      array.NewRecordBuilder(parquetProps.Allocator(), config.ArrowSchema),
		geometryTypeLookup: map[string]map[string]bool{},
		boundsLookup:       map[string]*orb.Bound{},
	}

	return writer, nil
}
func (w *FeatureWriter) Write(feature *Feature) error {
	arrowSchema := w.recordBuilder.Schema()
numFields := arrowSchema.NumFields()
for i := 0; i < numFields; i++ {
	field := arrowSchema.Field(i)
	builder := w.recordBuilder.Field(i)
	if err := w.append(feature, field, builder); err != nil {
		return err
	}
}
w.bufferedLength += 1
if w.bufferedLength >= w.maxRowGroupLength {
	return w.writeBuffered()
}
return nil
}

func (w *FeatureWriter) writeBuffered() error {
record := w.recordBuilder.NewRecord()
defer record.Release()
if err := w.fileWriter.WriteBuffered(record); err != nil {
	return err
}
w.bufferedLength = 0
return nil
}

func (w *FeatureWriter) append(feature *Feature, field arrow.Field, builder array.Builder) error {
name := field.Name
if w.geoMetadata.Columns[name] != nil {
	return w.appendGeometry(feature, field, builder)
}

value, ok := feature.Properties[name]
if !ok || value == nil {
	if !field.Nullable {
		return fmt.Errorf("field %q is required, but the property is missing in the feature", name)
	}
	builder.AppendNull()
	return nil
}

return w.appendValue(name, value, builder)
}

func (w *FeatureWriter) appendValue(name string, value any, builder array.Builder) error {
switch b := builder.(type) {
case *array.BooleanBuilder:
	v, ok := value.(bool)
	if !ok {
		return fmt.Errorf("expected %q to be a boolean, got %v", name, value)
	}
	b.Append(v)
case *array.StringBuilder:
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected %q to be a string, got %v", name, value)
	}
	b.Append(v)
case *array.Float64Builder:
	v, ok := value.(float64)
	if !ok {
		return fmt.Errorf("expected %q to be a float64, got %v", name, value)
	}
	b.Append(v)
case *array.ListBuilder:
	b.Append(true)
	valueBuilder := b.ValueBuilder()
	switch vb := valueBuilder.(type) {
	case *array.BooleanBuilder:
		v, ok := toUniformSlice[bool](value)
		if !ok {
			return fmt.Errorf("expected %q to be []bool, got %v", name, value)
		}
		vb.AppendValues(v, nil)
	case *array.StringBuilder:
		v, ok := toUniformSlice[string](value)
		if !ok {
			return fmt.Errorf("expected %q to be []string, got %v", name, value)
		}
		vb.AppendValues(v, nil)
	case *array.Float64Builder:
		v, ok := toUniformSlice[float64](value)
		if !ok {
			return fmt.Errorf("expected %q to be []float64, got %v", name, value)
		}
		vb.AppendValues(v, nil)
	case *array.StructBuilder:
		v, ok := value.([]any)
		if !ok {
			return fmt.Errorf("expected %q to be []any, got %v", name, value)
		}
		for _, item := range v {
			if err := w.appendValue(name, item, vb); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported list element builder type %#v", vb)
	}
case *array.StructBuilder:
	v, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("expected %q to be map[string]any, got %v", name, value)
	}
	t, ok := b.Type().(*arrow.StructType)
	if !ok {
		return fmt.Errorf("expected builder for %q to have a struct type, got %v", name, b.Type())
	}
	b.Append(true)
	for i := 0; i < b.NumField(); i += 1 {
		field := t.Field(i)
		name := field.Name
		fieldValue, ok := v[name]
		fieldBuilder := b.FieldBuilder(i)
		if !ok || fieldValue == nil {
			if !field.Nullable {
				return fmt.Errorf("field %q is required, but the property is missing", name)
			}
			fieldBuilder.AppendNull()
			continue
		}
		if err := w.appendValue(name, fieldValue, fieldBuilder); err != nil {
			return err
		}
	}
default:
	return fmt.Errorf("unsupported builder type %#v", b)
}

return nil
}

func toUniformSlice[T any](value any) ([]T, bool) {
if values, ok := value.([]T); ok {
	return values, true
}
slice, ok := value.([]any)
if !ok {
	return nil, false
}
values := make([]T, len(slice))
for i, v := range slice {
	t, ok := v.(T)
	if !ok {
		return nil, false
	}
	values[i] = t
}
return values, true
}

func (w *FeatureWriter) appendGeometry(feature *Feature, field arrow.Field, builder array.Builder) error {
name := field.Name
geomColumn := w.geoMetadata.Columns[name]

binaryBuilder, ok := builder.(*array.BinaryBuilder)
if !ok {
	return fmt.Errorf("expected column %q to have a binary type, got %s", name, builder.Type().Name())
}
var geometry orb.Geometry
if name == w.geoMetadata.PrimaryColumn {
	geometry = feature.Geometry
} else {
	if value, ok := feature.Properties[name]; ok {
		g, ok := value.(orb.Geometry)
		if !ok {
			return fmt.Errorf("expected %q to be a geometry, got %v", name, value)
		}
		geometry = g
	}
}
if geometry == nil {
	if !field.Nullable {
		return fmt.Errorf("feature missing required %q geometry", name)
	}
	binaryBuilder.AppendNull()
	return nil
}

if w.geometryTypeLookup[name] == nil {
	w.geometryTypeLookup[name] = map[string]bool{}
}
w.geometryTypeLookup[name][geometry.GeoJSONType()] = true

bounds := geometry.Bound()
if w.boundsLookup[name] != nil {
	bounds = bounds.Union(*w.boundsLookup[name])
}
w.boundsLookup[name] = &bounds

switch geomColumn.Encoding {
case EncodingWKB:
	data, err := wkb.Marshal(geometry)
	if err != nil {
		return fmt.Errorf("failed to encode %q as WKB: %w", name, err)
	}
	binaryBuilder.Append(data)
	return nil
case EncodingWKT:
	binaryBuilder.Append(wkt.Marshal(geometry))
	return nil
default:
	return fmt.Errorf("unsupported geometry encoding: %s", geomColumn.Encoding)
}
}

func (w *FeatureWriter) Close() error {
defer w.recordBuilder.Release()
if w.bufferedLength > 0 {
	if err := w.writeBuffered(); err != nil {
		return err
	}
}

geoMetadata := w.geoMetadata.Clone()
for name, bounds := range w.boundsLookup {
	if bounds != nil {
		if geoMetadata.Columns[name] == nil {
			geoMetadata.Columns[name] = getDefaultGeometryColumn()
		}
		geoMetadata.Columns[name].Bounds = []float64{
			bounds.Left(), bounds.Bottom(), bounds.Right(), bounds.Top(),
		}
	}
}
for name, types := range w.geometryTypeLookup {
	geometryTypes := []string{}
	if len(types) > 0 {
		for geometryType := range types {
			geometryTypes = append(geometryTypes, geometryType)
		}
	}
	if geoMetadata.Columns[name] == nil {
		geoMetadata.Columns[name] = getDefaultGeometryColumn()
	}
	geoMetadata.Columns[name].GeometryTypes = geometryTypes
}

data, err := json.Marshal(geoMetadata)
if err != nil {
	return fmt.Errorf("failed to encode %s file metadata", MetadataKey)
}
if err := w.fileWriter.AppendKeyValueMetadata(MetadataKey, string(data)); err != nil {
	return fmt.Errorf("failed to append %s file metadata", MetadataKey)
}
return w.fileWriter.Close()
}
*/
