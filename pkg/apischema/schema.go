package schema

// QueryParam 通用查询参数
type QueryParam struct {
	User       string              `json:"user"`
	Fields     map[string]bool     `json:"-"`
	NoFields   map[string]bool     `json:"-"`
	SortAsc    []string            `json:"sort_asc"`
	SortDesc   []string            `json:"sort_desc"`
	Match      map[string]string   `json:"match"`
	In         map[string][]string `json:"in"`
	Exclude    map[string]string   `json:"exclude"`
	Contains   map[string]string   `json:"contains"`
	NotIn      map[string][]string `json:"not_in"`
	Limit      int                 `json:"limit"`
	Offset     int                 `json:"offset"`
	GetRelated bool                `json:"get_related"`
}

// UpdateParam ...
type UpdateParam struct {
	ID      uint32                 `json:"id"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// RespData ...
type RespData struct {
	RetCode int         `json:"retcode"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// StatusText ...
type StatusText string

func (t StatusText) String() string {
	return string(t)
}

const (
	// OKStatus  ...
	OKStatus StatusText = "OK"
	// ErrorStatus ...
	ErrorStatus StatusText = "ERROR"
	// FailStatus ...
	FailStatus StatusText = "FAIL"
)

// StatusResult ...
type StatusResult struct {
	Status StatusText `json:"status"`
}

// ErrorResult ...
type ErrorResult struct {
	Error ErrorItem `json:"error"`
}

// ErrorItem ...
type ErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListResult ...
type ListResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

// PaginationResult ...
type PaginationResult struct {
	Total  int64 `json:"total"`
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
}

// PaginationParam ...
type PaginationParam struct {
	Pagination bool `form:"-"`
	OnlyCount  bool `form:"-"`
	Offset     int  `form:"offset,default=1"`
	Limit      int  `form:"limit,default=10" binding:"max=100"`
}

// GetCurrent ...
func (a PaginationParam) GetCurrent() int {
	return a.Offset
}

// GetPageSize ...
func (a PaginationParam) GetPageSize() int {
	pageSize := a.Limit
	if a.Limit <= 0 {
		pageSize = 100
	}
	return pageSize
}

// OrderDirection ...
type OrderDirection int

const (
	// OrderByASC ...
	OrderByASC OrderDirection = iota + 1
	// OrderByDESC ...
	OrderByDESC
)

// NewOrderFieldWithKeys Create order fields key and define key index direction
func NewOrderFieldWithKeys(keys []string, directions ...map[int]OrderDirection) []*OrderField {
	m := make(map[int]OrderDirection)
	if len(directions) > 0 {
		m = directions[0]
	}

	fields := make([]*OrderField, len(keys))
	for i, key := range keys {
		d := OrderByASC
		if v, ok := m[i]; ok {
			d = v
		}

		fields[i] = NewOrderField(key, d)
	}

	return fields
}

// NewOrderFields  ...
func NewOrderFields(orderFields ...*OrderField) []*OrderField {
	return orderFields
}

// NewOrderField ...
func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

// OrderField ...
type OrderField struct {
	Key       string
	Direction OrderDirection
}

// NewIDResult ...
func NewIDResult(id uint32) *IDResult {
	return &IDResult{
		ID: id,
	}
}

// IDResult ...
type IDResult struct {
	ID uint32 `json:"id"`
}

// DeleteFlag mark sql row deleted
type DeleteFlag int8

// DeleteFlag value
const (
	// NotDeleted ...
	NotDeleted DeleteFlag = 1
	// Deleted ...
	Deleted DeleteFlag = 2
)

// SecretDataForClient the info is encrypted for client using, not for real user
// just keep in memory, no logging, do not save to file/db
type SecretDataForClient struct {
	Data              string            `json:"data"`
	SerializationType SerializationType `json:"serialization_type"`
	ObjectType        string            `json:"object_type"`
}

// SerializationType ...
type SerializationType string

const (
	// SerializationJSON ...
	SerializationJSON SerializationType = "json"
	// SerializationYAML ...
	SerializationYAML SerializationType = "yaml"
	// SerializationTOML ...
	SerializationTOML SerializationType = "toml"
)
