package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AreaCode struct {
	Id           int       `orm:"column(id);auto" description:"主键"`
	AREACODE     string    `orm:"column(AREA_CODE);size(20);null" description:"地区编码"`
	AREAENG      string    `orm:"column(AREA_ENG);size(20);null" description:"地区英文"`
	AREANAME     string    `orm:"column(AREA_NAME);size(30);null" description:"地区名称"`
	COUNTRYCODE  string    `orm:"column(COUNTRY_CODE);size(20);null" description:"国家编码"`
	COUNTRYENG   string    `orm:"column(COUNTRY_ENG);size(20);null" description:"国家英文"`
	COUNTRYNAME  string    `orm:"column(COUNTRY_NAME);size(20);null" description:"国家名称"`
	PROVINCEENG  string    `orm:"column(PROVINCE_ENG);size(20);null" description:"省级英文"`
	PROVINCENAME string    `orm:"column(PROVINCE_NAME);size(20);null" description:"省级名称"`
	PARENTENG    string    `orm:"column(PARENT_ENG);size(20);null" description:"父级英文"`
	PARENTNAME   string    `orm:"column(PARENT_NAME);size(20);null" description:"父级名称"`
	LON          float64   `orm:"column(LON);null" description:"经度"`
	LAT          float64   `orm:"column(LAT);null" description:"纬度"`
	Makedate     time.Time `orm:"column(makedate);type(datetime);null;auto_now_add" description:"创建日期"`
	Modifydate   time.Time `orm:"column(modifydate);type(datetime);null;auto_now" description:"更新日期"`
}

func (t *AreaCode) TableName() string {
	return "bg_area_code"
}

func init() {
	orm.RegisterModel(new(AreaCode))
}

// QueryByArea select AreaCode by the area name
func QueryByArea(area string) map[string]interface{} {
	result := make(map[string]interface{})
	// var info []*models.AreaCode{}
	maps := AreaCode{}
	// var maps []orm.ParamsList
	o := orm.NewOrm()
	qs := o.QueryTable(new(AreaCode))
	err := qs.Filter("AREA_NAME", area).OrderBy("id").One(&maps, "Id", "AREACODE", "LON", "LAT")
	if err != nil {
		// 进行一次容错处理， 如果区不存在，则去查市。
		err = qs.Filter("PARENT_NAME", area).OrderBy("id").One(&maps, "Id", "AREACODE", "LON", "LAT")
		if err != nil {
			result["code"] = -100
			result["msg"] = "未知地区"
			return result
		}
	}

	result["code"] = 200
	result["msg"] = "查询成功"
	data := make(map[string]interface{})
	data["areaCode"] = maps.AREACODE
	data["lon"] = maps.LON
	data["lat"] = maps.LAT
	result["data"] = data
	return result
}

// AddAreaCode insert a new AreaCode into database and returns
// last inserted Id on success.
func AddAreaCode(m *AreaCode) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAreaCodeById retrieves AreaCode by Id. Returns error if
// Id doesn't exist
func GetAreaCodeById(id int) (v *AreaCode, err error) {
	o := orm.NewOrm()
	v = &AreaCode{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAreaCode retrieves all AreaCode matches certain condition. Returns empty list if
// no records exist
func GetAllAreaCode(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AreaCode))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []AreaCode
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateAreaCode updates AreaCode by Id and returns error if
// the record to be updated doesn't exist
func UpdateAreaCodeById(m *AreaCode) (err error) {
	o := orm.NewOrm()
	v := AreaCode{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAreaCode deletes AreaCode by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAreaCode(id int) (err error) {
	o := orm.NewOrm()
	v := AreaCode{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AreaCode{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
