package repos

import (
	"reflect"

	"gorm.io/gorm"
)

// GenericRepo is a Generic repo for performing certain scopped yet generic operations
type GenericRepo struct {
	*gorm.DB
}

// NewGenericRepo returns a new GenericRepo using the given db
func NewGenericRepo(db *gorm.DB) *GenericRepo {
	return &GenericRepo{
		DB: db,
	}
}

// GetAllOfOptions contains all options for the GetAllOf function
type GetAllOfOptions struct {
	Where string
}

// GetAllOf fetches all of whatever type is passed in
func (m *GenericRepo) GetAllOf(typeRef interface{}, options *GetAllOfOptions) (interface{}, error) {
	out := reflect.New(reflect.SliceOf(reflect.TypeOf(typeRef))).Interface()

	db := m.DB

	if options != nil {
		if options.Where != "" {
			db = db.Where(options.Where)
		}
	}

	db = db.Find(out)

	err := db.Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetAllOfRaw fetches all of whatever type is passed in using the raw query
func (m *GenericRepo) GetAllOfRaw(typeRef interface{}, query string, params ...interface{}) (interface{}, error) {
	out := reflect.New(reflect.SliceOf(reflect.TypeOf(typeRef))).Interface()

	err := m.DB.Raw(query, params...).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetByUUID fetches any entity by UUID through the out variable
func (m *GenericRepo) GetByUUID(out interface{}, uuid string) error {
	return m.Where("uuid = ?", uuid).Find(out).Error
}
