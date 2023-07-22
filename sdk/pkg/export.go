package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/language"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"reflect"
	"strings"
	"sync/atomic"
)

const (
	DefaultTagName   = "json"
	DefaultDescTag   = "description"
	DefaultPageSize  = 5000
	DefaultSheetSize = 100000
	DefaultLimitSize = 100000
	DefaultSheetName = "Sheet"
)

type FieldFunc func(ctx context.Context) []string
type TranslateFunc func(ctx context.Context) map[string]string
type StatusFunc func(ctx context.Context) map[string]map[string]string
type HeaderFunc func(ctx context.Context, excel *excelize.File) ([]interface{}, error)
type TotalFunc func(ctx context.Context) (int, error)
type ListFunc func(ctx context.Context, page, pageSize int) ([]interface{}, error)
type RawFunc func(ctx context.Context) (any, error)
type FinishFunc func(ctx context.Context) error
type ErrorFunc func(ctx context.Context, err error) error
type LangFunc func(ctx context.Context, key string) string
type PictureFunc func(ctx context.Context, url string) string
type RawStruct struct {
	FieldList   []string
	HeaderCols  []interface{}
	StatusEnums map[string]map[string]string
	Obj         any
}

type GraphicOptions excelize.GraphicOptions

type ExportOptions struct {
	FileName       string
	PageSize       int
	SheetSize      int
	LimitSize      int
	SheetPrefix    string
	TagName        string
	DescTag        string
	PictureKeys    []string
	RespItemStruct any
	FieldFunc      FieldFunc
	TranslateFunc  TranslateFunc
	StatusFunc     StatusFunc
	HeaderFunc     HeaderFunc
	TotalFunc      TotalFunc
	ListFunc       ListFunc
	SummaryFunc    RawFunc
	ParamsFunc     RawFunc
	FinishFunc     FinishFunc
	ErrorFunc      ErrorFunc
	PictureFunc    PictureFunc
	i18n           *gi18n.Manager
	GraphicOptions *GraphicOptions
}

func getDefaultExportOptions() ExportOptions {
	return ExportOptions{
		FileName:    "export.xlsx",
		PageSize:    DefaultPageSize,
		SheetSize:   DefaultSheetSize,
		SheetPrefix: DefaultSheetName,
		LimitSize:   DefaultLimitSize,
	}
}
func WithExportOptionsFileName(fileName string) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.FileName = fileName
	}
}
func WithExportOptionsSheetPrefix(prefix string) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.SheetPrefix = prefix
	}
}
func WithExportOptionsDescTag(tag string) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.DescTag = tag
	}
}
func WithExportOptionsTagName(tag string) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.TagName = tag
	}
}
func WithExportOptionsPictureKeys(keys []string) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.PictureKeys = keys
	}
}
func WithExportOptionsRespItemStruct(obj any) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.RespItemStruct = obj
	}
}
func WithExportOptionsPageSize(pageSize int) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.PageSize = pageSize
	}
}
func WithExportOptionsSheetSize(sheetSize int) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.SheetSize = sheetSize
	}
}
func WithExportOptionsLimitSize(limit int) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.LimitSize = limit
	}
}
func WithExportOptionsFieldFunc(f FieldFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.FieldFunc = f
	}
}
func WithExportOptionsTranslateFunc(f TranslateFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.TranslateFunc = f
	}
}
func WithExportOptionsStatusFunc(f StatusFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.StatusFunc = f
	}
}
func WithExportOptionsHeaderFunc(f HeaderFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.HeaderFunc = f
	}
}
func WithExportOptionsTotalFunc(f TotalFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.TotalFunc = f
	}
}
func WithExportOptionsListFunc(f ListFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.ListFunc = f
	}
}
func WithExportOptionsSummaryFunc(f RawFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.SummaryFunc = f
	}
}
func WithExportOptionsParamsFunc(f RawFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.ParamsFunc = f
	}
}
func WithExportOptionsFinishFunc(f FinishFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.FinishFunc = f
	}
}
func WithExportOptionsErrorFunc(f ErrorFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.ErrorFunc = f
	}
}
func WithExportOptionsI18n(i18n *gi18n.Manager) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.i18n = i18n
	}
}
func WithExportOptionsPictureFunc(f PictureFunc) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.PictureFunc = f
	}
}
func WithExportOptionsGraphicOptions(g *GraphicOptions) func(*ExportOptions) {
	return func(options *ExportOptions) {
		options.GraphicOptions = g
	}
}

type Export struct {
	ctx          context.Context
	excel        *excelize.File
	streamWriter *excelize.StreamWriter
	options      *ExportOptions
	fieldList    []string
	pictureMap   map[string]string
	statusEnums  map[string]map[string]string
	translates   map[string]string
	headerCols   []interface{}
	offset       int
	total        int
	page         int
	pageTotal    int
	sheetTotal   int
	bodyStyleId  int
	count        int32
}

func NewExport(ctx context.Context, lang string, optionFuncs ...func(*ExportOptions)) *Export {
	defaultOptions := getDefaultExportOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	// 分表大小小于分页大小，则以分页大小为分表大小
	if options.SheetSize < options.PageSize {
		glog.Warning(ctx, fmt.Sprintf("export excel SheetSize(%d) < PageSize(%d), SheetSize will be set to PageSize", options.SheetSize, options.PageSize))
		options.SheetSize = options.PageSize
	}
	if options.TagName == "" {
		options.TagName = DefaultTagName
	}
	if options.DescTag == "" {
		options.DescTag = DefaultDescTag
	}
	// 语言
	var preferred []language.Tag
	if lang != "" {
		var err error
		preferred, _, err = language.ParseAcceptLanguage(lang)
		if err != nil {
			// err
			lang = "zh"
		}
	}
	if preferred == nil {
		lang = "zh"
	}
	matcher := language.NewMatcher([]language.Tag{
		language.English,           // 英语
		language.Chinese,           // 中文
		language.SimplifiedChinese, // 简体中文
		language.Malay,             // 马来西亚语 Malay
	})
	code, _, _ := matcher.Match(preferred...)
	base, _ := code.Base()
	lang = base.String()
	ctx = gi18n.WithLanguage(ctx, lang)

	return &Export{
		ctx:     ctx,
		options: options,
	}
}

func DefaultStatusFunc(ctx context.Context, l LangFunc, fieldList []string, translateMap map[string]string) (statusEnums map[string]map[string]string) {
	statusEnums = make(map[string]map[string]string)
	for _, key := range fieldList {
		value := translateMap[key]
		// 解析枚举参数
		if strings.Contains(value, ":") {
			list := strings.Split(value, ":")
			if len(list) == 2 {
				eList := strings.Split(list[1], ",")
				for _, s := range eList {
					vList := strings.Split(s, "=")
					if len(vList) == 2 {
						if statusEnums[key] == nil {
							statusEnums[key] = make(map[string]string)
						}
						statusEnums[key][vList[0]] = l(ctx, vList[1])
					}
				}
			}
		}
	}
	return statusEnums
}

func DefaultHeaderFunc(ctx context.Context, l LangFunc, excel *excelize.File, fieldList []string, translateMap map[string]string) (headerCols []interface{}, err error) {
	styleID, err := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Family: "宋体", Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		glog.Warning(ctx, "export excel header NewStyle error:", err.Error())
		return nil, err
	}
	headerCols = make([]interface{}, 0)
	for _, key := range fieldList {
		value := translateMap[key]
		if strings.Contains(value, ":") {
			list := strings.Split(value, ":")
			value = list[0]
		}
		col := excelize.Cell{Value: l(ctx, value), StyleID: styleID}
		headerCols = append(headerCols, col)
	}
	return headerCols, nil
}

func DefaultFieldFunc(obj any, tag string) []string {
	return GetTags(obj, tag)
}

func DefaultTranslateFunc(obj any, tag, desc string) map[string]string {
	return GetTagsMap(obj, tag, desc)
}

func (e *Export) before(ctx context.Context) (err error) {
	e.excel = excelize.NewFile()
	if e.options.TotalFunc != nil {
		e.total, err = e.options.TotalFunc(ctx)
		if err != nil {
			return err
		}
	}
	// 0. 计算分页
	e.pageTotal = int(math.Ceil(float64(e.total) / float64(e.options.PageSize)))
	if e.pageTotal == 0 {
		e.pageTotal = 1
	}
	// 计算分表
	e.sheetTotal = int(math.Ceil(float64(e.total) / float64(e.options.SheetSize)))
	e.page = 0

	// 1. 字段列表
	if e.options.FieldFunc != nil {
		e.fieldList = e.options.FieldFunc(ctx)
	} else {
		if e.options.RespItemStruct == nil {
			return errors.New("export excel RespItemStruct is nil")
		}
		e.fieldList = DefaultFieldFunc(e.options.RespItemStruct, e.options.TagName)
	}
	// 2. 翻译列表
	if e.options.TranslateFunc != nil {
		e.translates = e.options.TranslateFunc(ctx)
	} else {
		if e.options.RespItemStruct == nil {
			return errors.New("export excel RespItemStruct is nil")
		}
		e.translates = DefaultTranslateFunc(e.options.RespItemStruct, e.options.TagName, e.options.DescTag)
	}
	// 3. 状态列表
	if e.options.StatusFunc != nil {
		e.statusEnums = e.options.StatusFunc(ctx)
	} else {
		e.statusEnums = DefaultStatusFunc(ctx, e.lang, e.fieldList, e.translates)
	}
	// 4. 表头构建
	if e.options.HeaderFunc != nil {
		e.headerCols, err = e.options.HeaderFunc(ctx, e.excel)
	} else {
		e.headerCols, err = DefaultHeaderFunc(ctx, e.lang, e.excel, e.fieldList, e.translates)
	}
	// 3. Set table body style
	e.bodyStyleId, err = e.excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "Calibri", Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		glog.Warning(ctx, "export excel body NewStyle error:", err.Error())
		return err
	}
	return err
}

func (e *Export) processorRaw(ctx context.Context, name string, r *RawStruct) (err error) {
	if r.FieldList == nil || len(r.FieldList) == 0 || r.HeaderCols == nil || len(r.HeaderCols) == 0 {
		return nil
	}

	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		glog.Warning(ctx, "export excel processorRaw CoordinatesToCellName error:", err.Error())
		return err
	}
	sheetName := name
	sheet, err := e.excel.NewSheet(sheetName)
	if err != nil {
		glog.Warning(ctx, "export excel processorRaw NewSheet error:", err.Error())
		return err
	}
	e.excel.SetActiveSheet(sheet)
	// 1. NewStreamWriter
	e.streamWriter, err = e.excel.NewStreamWriter(sheetName)
	if err != nil {
		glog.Warning(ctx, "export excel processorRaw NewStreamWriter error:", err.Error())
		return err
	}
	// 2. Write Header
	if err = e.streamWriter.SetRow(cell, r.HeaderCols, excelize.RowOpts{OutlineLevel: 1}); err != nil {
		glog.Warning(ctx, "export excel processorRaw SetRow error:", err.Error())
		return err
	}
	// 3. Write Body Data
	cell, err = excelize.CoordinatesToCellName(1, 2)
	if err != nil {
		glog.Warning(ctx, "export excel processorRaw CoordinatesToCellName error:", err.Error())
		return err
	}
	if r.Obj != nil {
		BodyRow := make([]interface{}, 0)
		mapObj := gconv.Map(r.Obj)
		for _, v := range r.FieldList {
			col := excelize.Cell{StyleID: e.bodyStyleId, Value: mapObj[v]}
			if r.StatusEnums[v] != nil {
				col.Value = r.StatusEnums[v][gconv.String(mapObj[v])]
			}
			BodyRow = append(BodyRow, col)
		}

		if err = e.streamWriter.SetRow(cell, BodyRow); err != nil {
			glog.Warning(ctx, "export excel processorRaw body SetRow error:", err.Error())
			return err
		}
	}

	// 4. Flush
	if err = e.streamWriter.Flush(); err != nil {
		glog.Warning(ctx, "export excel processorRaw Flush error:", err.Error())
		return err
	}
	return nil
}

func (e *Export) preProcessor(ctx context.Context, ptr any) (r *RawStruct, err error) {
	fields, err := GetNonEmptyFields(ptr, e.options.TagName)
	if err != nil {
		return nil, err
	}
	fieldList := make([]string, 0)
	for key := range fields {
		fieldList = append(fieldList, key)
	}
	translateMap := GetTagsMap(ptr, e.options.TagName, e.options.DescTag)
	statusEnums := DefaultStatusFunc(ctx, e.lang, fieldList, translateMap)
	headerCols, err := DefaultHeaderFunc(ctx, e.lang, e.excel, fieldList, translateMap)
	if err != nil {
		return nil, err
	}
	raw := RawStruct{
		FieldList:   fieldList,
		StatusEnums: statusEnums,
		HeaderCols:  headerCols,
		Obj:         ptr,
	}

	return &raw, err
}

func (e *Export) processor(ctx context.Context) (err error) {
	// excel 分表处理
	for currentSheet := 1; currentSheet <= e.sheetTotal; currentSheet++ {
		e.offset = 1
		cell, err := excelize.CoordinatesToCellName(1, e.offset)
		if err != nil {
			glog.Warning(ctx, "export excel CoordinatesToCellName error:", err.Error())
			return err
		}
		sheetName := fmt.Sprintf("%s%d", e.options.SheetPrefix, currentSheet)
		sheet, err := e.excel.NewSheet(sheetName)
		if err != nil {
			glog.Warning(ctx, "export excel NewSheet error:", err.Error())
			return err
		}
		e.excel.SetActiveSheet(sheet)
		// 1. NewStreamWriter
		e.streamWriter, err = e.excel.NewStreamWriter(sheetName)
		if err != nil {
			glog.Warning(ctx, "export excel NewStreamWriter error:", err.Error())
			return err
		}
		// 2. Write Header
		if err = e.streamWriter.SetRow(cell, e.headerCols, excelize.RowOpts{OutlineLevel: 1}); err != nil {
			glog.Warning(ctx, "export excel SetRow error:", err.Error())
			return err
		}
		// 3. Write Body Data
		maxPage := int(math.Ceil(float64(e.pageTotal) / float64(e.sheetTotal)))
		e.pictureMap = make(map[string]string)
		for s := 1; s <= maxPage && e.page < e.pageTotal; s++ {
			e.page++
			// 分页查询
			list, err := e.options.ListFunc(ctx, e.page, e.options.PageSize)
			if err != nil {
				return err
			}
			err = e.exportList(ctx, list)
			if err != nil {
				return err
			}
			if e.count >= int32(e.options.LimitSize) {
				glog.Warning(ctx, "export excel LimitSize:", e.options.LimitSize, "current count:", e.total)
				break
			}

			glog.Info(ctx, fmt.Sprintf("export excel filename: %s, total row: %d, sheetName: %s, current page: %d, Query list len: %d, current query: %d, current offset: %d success", e.options.FileName, e.total, sheetName, e.page, len(list), s, e.offset))
		}
		// 4. Set Picture
		for col, val := range e.pictureMap {
			graphicOptions := excelize.GraphicOptions{}
			if e.options.GraphicOptions != nil {
				graphicOptions = excelize.GraphicOptions(*e.options.GraphicOptions)
			}
			err = e.excel.AddPicture(sheetName, col, val, &graphicOptions)
			if err != nil {
				glog.Warning(ctx, "export excel AddPicture error:", err.Error())
				return err
			}
		}
		// 5. Flush
		if err = e.streamWriter.Flush(); err != nil {
			glog.Warning(ctx, "export excel Flush error:", err.Error())
			return err
		}
	}
	// 6. Summary
	if e.options.SummaryFunc != nil {
		summaryPtr, err := e.options.SummaryFunc(ctx)
		if err != nil {
			glog.Warning(ctx, "export excel SummaryFunc error:", err.Error())
			return err
		}
		raw, err := e.preProcessor(ctx, summaryPtr)
		if err != nil {
			glog.Warning(ctx, "export excel PreProcessor error:", err.Error())
			return err
		}
		err = e.processorRaw(ctx, "summary", raw)
		if err != nil {
			glog.Warning(ctx, "export excel processorRaw error:", err.Error())
			return err
		}
	}
	// 7. Params
	if e.options.ParamsFunc != nil {
		paramsPtr, err := e.options.ParamsFunc(ctx)
		if err != nil {
			glog.Warning(ctx, "export excel ParamsFunc error:", err.Error())
			return err
		}
		raw, err := e.preProcessor(ctx, paramsPtr)
		if err != nil {
			glog.Warning(ctx, "export excel PreProcessor error:", err.Error())
			return err
		}
		err = e.processorRaw(ctx, "params", raw)
		if err != nil {
			glog.Warning(ctx, "export excel processorRaw error:", err.Error())
			return err
		}
	}
	// Delete default Sheet1
	if e.options.SheetPrefix != "Sheet" {
		err = e.excel.DeleteSheet("Sheet1")
		if err != nil {
			glog.Warning(ctx, "export excel DeleteSheet error:", err.Error())
			return err
		}
	}
	// Set first sheet active
	e.excel.SetActiveSheet(0)
	// 8. Save
	if err = e.excel.SaveAs(e.options.FileName); err != nil {
		glog.Warning(ctx, "export excel SaveAs error:", err.Error())
		return err
	}
	// 9. Finish
	if e.options.FinishFunc != nil {
		err = e.options.FinishFunc(ctx)
		if err != nil {
			glog.Warning(ctx, "export excel FinishFunc error:", err.Error())
			return err
		}
	}
	return nil
}

func (e *Export) exportList(ctx context.Context, list []any) error {
	// 4. For loop Set table body
	for _, item := range list {
		e.offset++
		cell, err := excelize.CoordinatesToCellName(1, e.offset)
		if err != nil {
			glog.Warning(ctx, "export excel CoordinatesToCellName error:", err.Error())
			return err
		}
		BodyRow := make([]interface{}, 0)
		mapObj := gconv.Map(item)

		for k, v := range e.fieldList {
			col := excelize.Cell{StyleID: e.bodyStyleId, Value: mapObj[v]}
			if e.statusEnums[v] != nil {
				// 状态翻译
				col.Value = e.statusEnums[v][gconv.String(mapObj[v])]
			}
			// 图片处理
			if e.options.PictureFunc != nil && InSlice(v, e.options.PictureKeys) {
				imgPath := e.options.PictureFunc(ctx, gconv.String(mapObj[v]))
				colName, err := excelize.CoordinatesToCellName(k+1, e.offset)
				if err != nil {
					glog.Warning(ctx, "export excel CoordinatesToCellName error:", err.Error())
					return err
				}
				e.pictureMap[colName] = imgPath
				col.Value = ""
			}
			BodyRow = append(BodyRow, col)
		}

		if err = e.streamWriter.SetRow(cell, BodyRow); err != nil {
			glog.Warning(ctx, "export excel body SetRow error:", err.Error())
			return err
		}
		atomic.AddInt32(&e.count, 1)
	}
	return nil
}

func (e *Export) after(ctx context.Context) (err error) {
	// 5. Close
	if err = e.excel.Close(); err != nil {
		glog.Warning(ctx, "export excel Close error:", err.Error())
		return err
	}
	return nil
}
func (e *Export) lang(ctx context.Context, key string) string {
	if e.options.i18n != nil {
		return e.options.i18n.Translate(ctx, key)
	}
	return key
}

func (e *Export) Run() (err error) {
	err = e.before(e.ctx)
	if err != nil {
		if e.options.ErrorFunc != nil {
			return e.options.ErrorFunc(e.ctx, err)
		}
		return err
	}
	err = e.processor(e.ctx)
	if err != nil {
		if e.options.ErrorFunc != nil {
			return e.options.ErrorFunc(e.ctx, err)
		}
		return err
	}
	defer func() {
		err = e.after(e.ctx)
		if err != nil {
			if e.options.ErrorFunc != nil {
				_ = e.options.ErrorFunc(e.ctx, err)
			} else {
				glog.Warning(e.ctx, "export excel after error:", err.Error())
			}
			return
		}
	}()

	return nil
}

// GetNonEmptyFields returns a map of non-empty fields of a struct
func GetNonEmptyFields(obj interface{}, tagName string) (map[string]any, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj must be a pointer to a struct")
	}

	nonEmptyFields := make(map[string]any)
	v := reflect.ValueOf(obj).Elem()
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" && tag != "-" {
			value := v.Field(i)
			if value.IsValid() && value.CanInterface() && !isEmpty(value) {
				nonEmptyFields[tag] = value.Interface()
			}
		}
	}

	return nonEmptyFields, nil
}

// IsEmpty returns true if the given value is empty
func isEmpty(value reflect.Value) bool {
	zero := reflect.Zero(value.Type())
	return reflect.DeepEqual(value.Interface(), zero.Interface())
}

// GetTags returns a slice of any tags of a struct
func GetTags(obj any, tagName string) []string {
	t := reflect.TypeOf(obj)
	if t == nil || (t.Kind() == reflect.Ptr && t.Elem().Kind() != reflect.Struct) && t.Kind() != reflect.Struct {
		return nil
	}

	tags := make([]string, 0)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" && tag != "-" {
			tags = append(tags, tag)
		}
	}

	return tags
}

// GetTagsMap returns a map of any tags of a struct
func GetTagsMap(obj any, tagName, vTagName string) map[string]string {
	t := reflect.TypeOf(obj)
	if t == nil || (t.Kind() == reflect.Ptr && t.Elem().Kind() != reflect.Struct) && t.Kind() != reflect.Struct {
		return nil
	}

	tagsValue := make(map[string]string, 0)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" && tag != "-" {
			value := field.Tag.Get(vTagName)
			tagsValue[tag] = value
		}
	}

	return tagsValue
}

// InSlice returns true if the needle string is found in the haystack slice
func InSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
