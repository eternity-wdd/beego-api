package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ApiLog struct {
	Id           int       `orm:"column(id);auto" description:"ID"`
	Name         string    `orm:"column(name);size(64)" description:"接口名称URL"`
	Platform     string    `orm:"column(platform);size(32)" description:"接口所属平台"`
	Controller   string    `orm:"column(controller);size(32)" description:"controller控制器"`
	Action       string    `orm:"column(action);size(32)" description:"action方法"`
	Method       string    `orm:"column(method);size(16)" description:"URL请求方法"`
	Params       string    `orm:"column(params);size(2000)" description:"请求参数json"`
	Header       string    `orm:"column(header);size(255)" description:"请求头"`
	Client       string    `orm:"column(client);size(32)" description:"客户端IP或域名"`
	Server       string    `orm:"column(server);size(64)" description:"服务端名称（端口不是80时，会添加端口号）"`
	ResponseCode string    `orm:"column(response_code);size(8)" description:"响应状态码"`
	ResponseMsg  string    `orm:"column(response_msg);size(255)" description:"响应消息"`
	ResponseData string    `orm:"column(response_data)" description:"响应数据"`
	CreateTime   time.Time `orm:"column(create_time);type(timestamp);auto_now" description:"请求时间（在NAVICAT中检索日志时使用）"`
	StartTime    int       `orm:"column(start_time);null" description:"接口开始时间"`
	EndTime      int       `orm:"column(end_time);null" description:"接口结束时间"`
	Life         int       `orm:"column(life)" description:"接口耗时/单位毫秒"`
}

func (t *ApiLog) TableName() string {
	return "bg_api_log"
}

func init() {
	orm.RegisterModel(new(ApiLog))
}

// AddApiLog insert a new ApiLog into database and returns
// last inserted Id on success.
func AddApiLog(m *ApiLog) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetApiLogById retrieves ApiLog by Id. Returns error if
// Id doesn't exist
func GetApiLogById(id int) (v *ApiLog, err error) {
	o := orm.NewOrm()
	v = &ApiLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApiLog retrieves all ApiLog matches certain condition. Returns empty list if
// no records exist
func GetAllApiLog(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ApiLog))
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

	var l []ApiLog
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

// UpdateApiLog updates ApiLog by Id and returns error if
// the record to be updated doesn't exist
func UpdateApiLogById(m *ApiLog) (err error) {
	o := orm.NewOrm()
	v := ApiLog{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApiLog deletes ApiLog by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApiLog(id int) (err error) {
	o := orm.NewOrm()
	v := ApiLog{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ApiLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
