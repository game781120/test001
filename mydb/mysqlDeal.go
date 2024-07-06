package mydb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"thundersoft.com/brain/DigitalVisitor/conf"
	"thundersoft.com/brain/DigitalVisitor/logs"
	"thundersoft.com/brain/DigitalVisitor/utils"
	"time"
)

type MySqlDB struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DbName    string `yaml:"dbName"`
	TableName string `yaml:"tableName"`
	FileIDS   string `yaml:"fileIDS"`
	Knowledge string `yaml:"knowledge"`
	DB        *sql.DB
}

var MyDB *MySqlDB

func Start() (bool, error) {
	MyDB = &MySqlDB{
		Username: conf.ConfigInfo.Mysql.Username,
		Password: conf.ConfigInfo.Mysql.Password,
		Host:     conf.ConfigInfo.Mysql.Host,
		Port:     conf.ConfigInfo.Mysql.Port,
		DbName:   conf.ConfigInfo.Mysql.DbName,
	}
	return Connect()
}

func Close() {
	MyDB.DB.Close()
}
func Connect() (bool, error) {
	var err error
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", MyDB.Username,
		MyDB.Password, MyDB.Host, MyDB.Port, MyDB.DbName)
	MyDB.DB, err = sql.Open("mysql", dbInfo)
	if err != nil {
		fmt.Printf("mysql 连接失败,%s err: %s", dbInfo, err.Error())
		return false, err
	}
	fmt.Printf("mysql 连接成功 %s", dbInfo)
	return true, nil
}

func Query(req *utils.QueryRecordReq) ([]utils.Appointment, error) {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return nil, err
		}
	}

	var sqlStr string
	sqlStr = "SELECT appointment_id,customer_mobile,customer_name," +
		"customer_company,customer_department,customer_position,visiting_datetime," +
		"visiting_address,visiting_description,visiting_totals," +
		"reception_mobile,reception_name,reception_company," +
		"reception_department,reception_position,sign_in_state,sign_in_datetime," +
		"appointment_mobile,appointment_name,appointment_company," +
		"appointment_department,appointment_position," +
		"update_datetime,create_datetime FROM t_appointments where is_del=0 "

	if len(req.KeyWords) > 0 {
		sqlStr += fmt.Sprintf(" and (customer_name like  '%%%s%%' or reception_name like '%%%s%%' or customer_mobile like '%%%s%%')",
			req.KeyWords, req.KeyWords, req.KeyWords)
	}

	if len(req.StartTime) > 0 {
		sqlStr += fmt.Sprintf(" and visiting_datetime >= '%s'", req.StartTime)
	}

	if len(req.EndTime) > 0 {
		sqlStr += fmt.Sprintf(" and visiting_datetime <= '%s'", req.EndTime)
	}

	if req.SignInState == 1 || req.SignInState == 2 {
		sqlStr += fmt.Sprintf(" and sign_in_state = %d", req.SignInState)
	}
	if len(req.ReceptionName) > 0 {
		sqlStr += fmt.Sprintf(" and reception_name = '%s'", req.ReceptionName)
	}

	if len(req.UserOpenId) > 0 {
		sqlStr += fmt.Sprintf(" and user_open_id = '%s'", req.UserOpenId)
	}

	fmt.Printf("sqlStr=%s\n", sqlStr)

	// 查询数据
	rows, err := MyDB.DB.Query(sqlStr)
	if err != nil {
		logs.Errorf(err.Error())
		return nil, err
	}
	defer rows.Close()
	appointmentList := make([]utils.Appointment, 0)
	for rows.Next() {
		var appointment utils.Appointment
		var signInDataTime, visitingDataTime, createDatetime any

		err := rows.Scan(&appointment.AppointmentId,
			&appointment.CustomerMobile, &appointment.CustomerName,
			&appointment.CustomerCompany, &appointment.CustomerDepartment, &appointment.CustomerPosition,
			&visitingDataTime,
			&appointment.VisitingAddress,
			&appointment.VisitingDescription,
			&appointment.VisitingTotals,
			&appointment.ReceptionMobile,
			&appointment.ReceptionName,
			&appointment.ReceptionCompany,
			&appointment.ReceptionDepartment,
			&appointment.ReceptionPosition,
			&appointment.SignInState,
			&signInDataTime,
			&appointment.AppointmentMobile, &appointment.AppointmentName, &appointment.AppointmentCompany,
			&appointment.AppointmentDepartment, &appointment.AppointmentPosition,
			&appointment.UpdateDatetime,
			&createDatetime)
		if err != nil {
			logs.Error(err)
			continue
		}
		if signInDataTime == nil {
			appointment.SignInDataTime = ""
		} else {
			appointment.SignInDataTime = signInDataTime.(time.Time).Format("2006-01-02 15:04:05")
		}

		if visitingDataTime == nil {
			appointment.VisitingDataTime = ""
		} else {
			appointment.VisitingDataTime = visitingDataTime.(time.Time).Format("2006-01-02 15:04:05")
		}

		if createDatetime == nil {
			appointment.CreateDatetime = ""
		} else {
			appointment.CreateDatetime = createDatetime.(time.Time).Format("2006-01-02 15:04:05")
		}

		appointmentList = append(appointmentList, appointment)

	}
	return appointmentList, nil

}

func InsertData(appointment *utils.Appointment) string {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return ""
		}
	}
	u := uuid.New()
	appointment.AppointmentId = u.String()
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	appointment.UpdateDatetime = formattedTime
	appointment.CreateDatetime = formattedTime

	var signInDataTime any
	if appointment.IsSignIn {
		signInDataTime = formattedTime
		appointment.SignInState = 1
	} else {
		signInDataTime = nil
	}

	// 插入数据
	sqlStr := "INSERT INTO t_appointments(appointment_id,customer_mobile,customer_name," +
		"customer_company,customer_department,customer_position,visiting_datetime," +
		"visiting_address,visiting_description,visiting_totals," +
		"reception_mobile,reception_name,user_open_id,reception_company," +
		"reception_department,reception_position,sign_in_state,sign_in_datetime," +
		"appointment_mobile,appointment_name,appointment_company," +
		"appointment_department,appointment_position," +
		"update_datetime,create_datetime)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	insertStmt, err := MyDB.DB.Prepare(sqlStr)
	if err != nil {
		logs.Error(err)
	}
	defer insertStmt.Close()
	_, err = insertStmt.Exec(appointment.AppointmentId,
		appointment.CustomerMobile, appointment.CustomerName,
		appointment.CustomerCompany, appointment.CustomerDepartment, appointment.CustomerPosition,
		appointment.VisitingDataTime, appointment.VisitingAddress,
		appointment.VisitingDescription, appointment.VisitingTotals,
		appointment.ReceptionMobile, appointment.ReceptionName, appointment.UserOpenId,
		appointment.ReceptionCompany, appointment.ReceptionDepartment, appointment.ReceptionPosition,
		appointment.SignInState, signInDataTime,
		appointment.AppointmentMobile, appointment.AppointmentName,
		appointment.AppointmentCompany, appointment.AppointmentDepartment, appointment.AppointmentPosition,
		appointment.UpdateDatetime, appointment.CreateDatetime)
	if err == nil {
		fmt.Println("成功插入数据", appointment.AppointmentId)
	} else {
		fmt.Println("插入数据失败", err)
		return ""
	}
	return appointment.AppointmentId
}

func UpdateSingIn(appointmentId string) error {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return err
		}
	}
	// 更新数据
	updateStmt, err := MyDB.DB.Prepare("UPDATE t_appointments SET sign_in_state=1,sign_in_datetime=? WHERE appointment_id=?")
	if err != nil {
		logs.Error(err)
	}
	defer updateStmt.Close()
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	_, err = updateStmt.Exec(formattedTime, appointmentId)
	if err == nil {
		fmt.Println("成功更新数据")
	}
	return nil
}

func QueryAppointments(appointmentId string) (string, string, string, error) {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return "", "", "", err
		}
	}
	var sqlStr string
	sqlStr = fmt.Sprintf("SELECT user_open_id,customer_name,customer_company FROM t_appointments WHERE appointment_id= '%s'", appointmentId)
	fmt.Printf("sqlStr=%s\n", sqlStr)
	// 查询数据
	rows, err := MyDB.DB.Query(sqlStr)
	if err != nil {
		logs.Errorf(err.Error())
		return "", "", "", err
	}
	defer rows.Close()
	for rows.Next() {
		var userOpenid, name, company string
		err := rows.Scan(&userOpenid, &name, &company)
		if err != nil {
			logs.Error(err)
			return "", "", "", err
		} else {
			return userOpenid, name, company, nil
		}
	}
	return "", "", "", nil

}
func DelData(appointmentId string) error {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err0 := Connect(); err0 != nil {
			return err0
		}
	}
	// 更新数据
	updateStmt, err1 := MyDB.DB.Prepare("UPDATE t_appointments SET is_del=1,update_datetime=? WHERE appointment_id = ?")
	if err1 != nil {
		logs.Error(err1)
	}
	defer updateStmt.Close()
	formattedTime := time.Now().Format("2006-01-02 15:04:05")

	_, err = updateStmt.Exec(formattedTime, appointmentId)
	if err == nil {
		fmt.Println("成功更新数据")
	}
	return nil
}

func IsManager(req *utils.ManagerData) (bool, error) {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return false, err
		}
	}
	var sqlStr string
	sqlStr = fmt.Sprintf("SELECT mobile_number,name FROM t_manager where 1=1")
	if len(req.Name) > 0 {
		sqlStr += " and name = " + "'" + req.Name + "'"
	}
	if len(req.MobileNumber) > 0 {
		sqlStr += " and mobile_number = " + "'" + req.MobileNumber + "'"
	}

	// 查询数据
	rows, err := MyDB.DB.Query(sqlStr)
	if err != nil {
		logs.Errorf(err.Error())
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var res utils.ManagerData
		err := rows.Scan(&res.MobileNumber, &res.Name)
		if err != nil {
			logs.Error(err)
			return false, err
		} else {
			return true, nil
		}

	}
	return false, nil

}

func InsertManager(managerData utils.ManagerData) (bool, error) {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return false, err
		}
	}
	b, err := IsManager(&managerData)
	if err == nil && b {
		return true, err
	}
	// 插入数据
	sqlStr := "INSERT INTO t_manager(mobile_number,name)VALUES(?,?)"
	insertStmt, err := MyDB.DB.Prepare(sqlStr)
	if err != nil {
		logs.Error(err)
		return false, err
	}
	defer insertStmt.Close()
	_, err = insertStmt.Exec(managerData.MobileNumber, managerData.Name)
	if err == nil {
		fmt.Println("成功插入数据", managerData)
	} else {
		fmt.Println("插入数据失败", err)
		return false, err
	}
	return true, err
}

func QueryManager(keyWords string, pCurrent, pSize int) ([]utils.ManagerData, error) {
	managerList := make([]utils.ManagerData, 0)
	resManagerList := make([]utils.ManagerData, 0)
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return managerList, err
		}
	}
	var sqlStr string
	sqlStr = fmt.Sprintf("SELECT mobile_number,name FROM t_manager ")
	if len(keyWords) > 0 {
		sqlStr += fmt.Sprintf(" where  name like  '%%%s%%' or mobile_number like '%%%s%%'", keyWords, keyWords)
	}
	fmt.Printf("sqlStr=%s\n", sqlStr)
	// 查询数据
	rows, err := MyDB.DB.Query(sqlStr)
	if err != nil {
		logs.Errorf(err.Error())
		return managerList, err
	}
	defer rows.Close()
	for rows.Next() {
		var res utils.ManagerData
		err := rows.Scan(&res.MobileNumber, &res.Name)
		if err != nil {
			logs.Error(err)
			return managerList, err
		} else {
			managerList = append(managerList, res)
		}
	}

	if pSize != 0 {
		if len(managerList) < pSize {
			pSize = len(managerList)
		}
		if pCurrent == 0 {
			resManagerList = managerList[(pCurrent)*pSize : (pCurrent+1)*pSize]
		} else {
			resManagerList = managerList[(pCurrent-1)*pSize : (pCurrent)*pSize]
		}

	} else {
		resManagerList = managerList
	}
	return resManagerList, nil

}

func DelManager(req utils.ManagerData) (bool, error) {
	// 测试连接
	err := MyDB.DB.Ping()
	if err != nil {
		if _, err := Connect(); err != nil {
			return false, err
		}
	}

	// 插入数据
	sqlStr := "DELETE FROM t_manager WHERE 1=1 "
	if len(req.Name) > 0 {
		sqlStr += " and name = " + "'" + req.Name + "'"
	}
	if len(req.MobileNumber) > 0 {
		sqlStr += " and mobile_number = " + "'" + req.MobileNumber + "'"
	}

	delStmt, err := MyDB.DB.Prepare(sqlStr)
	if err != nil {
		logs.Error(err)
		return false, err
	}
	defer delStmt.Close()
	_, err = delStmt.Exec()
	if err == nil {
		fmt.Println("成功删除数据")
	} else {
		fmt.Println("删除数据失败", err)
		return false, err
	}
	return true, err
}
