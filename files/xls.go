package files

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"thundersoft.com/brain/DigitalVisitor/utils"

	//"thundersoft.com/brain/lzu-agent/config"
	//"thundersoft.com/brain/lzu-agent/db"
	//"thundersoft.com/brain/lzu-agent/utils"
	"time"
)

func WriteFile(appointment []utils.Appointment) string {
	// 创建一个新的XLSX文件
	f := excelize.NewFile()
	// 创建一个工作表
	index, _ := f.NewSheet("Sheet1")

	startColASCII := 'A'
	endColASCII := 'W'
	// 列宽
	colWidth := 20
	// 使用 for 循环设置列宽
	for colChar := startColASCII; colChar <= endColASCII; colChar++ {
		// 将 ASCII 字符转换为字符串
		colStr := string(colChar)
		// 调用 SetColWidth 方法设置列宽
		f.SetColWidth("Sheet1", colStr, colStr, float64(colWidth))
	}

	f.SetCellValue("Sheet1", "A1", "预约ID")
	f.SetCellValue("Sheet1", "B1", "来访人手机号")
	f.SetCellValue("Sheet1", "C1", "来访人姓名")
	f.SetCellValue("Sheet1", "D1", "来访人公司")
	f.SetCellValue("Sheet1", "E1", "来访人部门")
	f.SetCellValue("Sheet1", "F1", "来访人职务")
	f.SetCellValue("Sheet1", "G1", "来访时间")
	f.SetCellValue("Sheet1", "H1", "来访地点")
	f.SetCellValue("Sheet1", "I1", "来访事由")
	f.SetCellValue("Sheet1", "J1", "来访人数")

	f.SetCellValue("Sheet1", "K1", "接待人手机号")
	f.SetCellValue("Sheet1", "L1", "接待人姓名")
	f.SetCellValue("Sheet1", "M1", "接待人公司")
	f.SetCellValue("Sheet1", "N1", "接待人部门")
	f.SetCellValue("Sheet1", "O1", "接待人职务")

	f.SetCellValue("Sheet1", "P1", "状态")
	f.SetCellValue("Sheet1", "Q1", "签到时间")
	f.SetCellValue("Sheet1", "R1", "预约时间")
	f.SetCellValue("Sheet1", "S1", "预约人手机号")
	f.SetCellValue("Sheet1", "T1", "预约人姓名")
	f.SetCellValue("Sheet1", "U1", "预约人公司")
	f.SetCellValue("Sheet1", "V1", "预约人部门")
	f.SetCellValue("Sheet1", "W1", "预约人职务")

	i := 1
	for _, record := range appointment {
		i++
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), record.AppointmentId)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), record.CustomerMobile)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), record.CustomerName)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), record.CustomerCompany)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i), record.CustomerDepartment)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", i), record.CustomerPosition)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", i), record.VisitingDataTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", i), record.VisitingAddress)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", i), record.VisitingDescription)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", i), record.VisitingTotals)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", i), record.ReceptionMobile)
		f.SetCellValue("Sheet1", fmt.Sprintf("L%d", i), record.ReceptionName)
		f.SetCellValue("Sheet1", fmt.Sprintf("M%d", i), record.ReceptionCompany)
		f.SetCellValue("Sheet1", fmt.Sprintf("N%d", i), record.ReceptionDepartment)
		f.SetCellValue("Sheet1", fmt.Sprintf("O%d", i), record.ReceptionPosition)
		f.SetCellValue("Sheet1", fmt.Sprintf("P%d", i), record.SignInState)
		f.SetCellValue("Sheet1", fmt.Sprintf("Q%d", i), record.SignInDataTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("R%d", i), record.CreateDatetime)
		f.SetCellValue("Sheet1", fmt.Sprintf("S%d", i), record.AppointmentMobile)
		f.SetCellValue("Sheet1", fmt.Sprintf("T%d", i), record.AppointmentName)
		f.SetCellValue("Sheet1", fmt.Sprintf("U%d", i), record.AppointmentCompany)
		f.SetCellValue("Sheet1", fmt.Sprintf("V%d", i), record.AppointmentDepartment)
		f.SetCellValue("Sheet1", fmt.Sprintf("W%d", i), record.AppointmentPosition)
	}

	// 设置默认工作表
	f.SetActiveSheet(index)

	now := time.Now()
	// 格式化日期和时间
	date := now.Format("20060102")
	hour := now.Format("15")
	minute := now.Format("04")
	// 构建文件名
	filename := fmt.Sprintf("预约来访记录_%s-%s-%s.xlsx", date, hour, minute)
	path := "./temp/"
	// 创建路径
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println("无法创建路径:", err)
		return ""
	}
	// 保存文件
	err = f.SaveAs(path + filename)
	if err != nil {
		fmt.Println("无法保存文件:", err)
		return ""
	}

	fmt.Println("文件保存成功")

	return path + filename
}

// readXlsxData 读取 XLSX 文件数据并返回 XlsxRow 切片
func ReadXlsxData(filePath string) ([]utils.AppointmentTemp, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var data []utils.AppointmentTemp
	for _, row := range rows {

		data = append(data, utils.AppointmentTemp{
			AppointmentId:         row[0],
			CustomerMobile:        row[1],
			CustomerName:          row[2],
			CustomerCompany:       row[3],
			CustomerDepartment:    row[4],
			CustomerPosition:      row[5],
			VisitingDataTime:      row[6],
			VisitingAddress:       row[7],
			VisitingDescription:   row[8],
			VisitingTotals:        row[9],
			ReceptionMobile:       row[10],
			ReceptionName:         row[11],
			ReceptionCompany:      row[12],
			ReceptionDepartment:   row[13],
			ReceptionPosition:     row[14],
			SignInState:           row[15],
			SignInDataTime:        row[16],
			CreateDatetime:        row[17],
			AppointmentMobile:     row[18],
			AppointmentName:       row[19],
			AppointmentCompany:    row[20],
			AppointmentDepartment: row[21],
			AppointmentPosition:   row[22],
		})
	}

	return data, nil
}

//func CreateFile() (string, string, *excelize.File, int, int, error) {
//	// 创建一个新的XLSX文件
//	f := excelize.NewFile()
//	// 创建一个工作表
//	index, _ := f.NewSheet("Sheet1")
//	f.SetColWidth("Sheet1", "A", "A", 50)
//	f.SetColWidth("Sheet1", "B", "B", 100)
//	f.SetColWidth("Sheet1", "C", "C", 20)
//	f.SetColWidth("Sheet1", "D", "D", 100)
//	f.SetColWidth("Sheet1", "D", "E", 100)
//	f.SetCellValue("Sheet1", "A1", "文件名")
//	f.SetCellValue("Sheet1", "B1", "总结")
//	f.SetCellValue("Sheet1", "C1", "状态")
//	f.SetCellValue("Sheet1", "D1", "原因")
//	f.SetCellValue("Sheet1", "E1", "发给模型前的数据")
//
//	now := time.Now()
//	// 格式化日期和时间
//	date := now.Format("20060102")
//	hour := now.Format("15")
//	minute := now.Format("04")
//	// 构建文件名
//	filename := fmt.Sprintf("复选结果_%s_%s_%s.xlsx", date, hour, minute)
//	path := "./temp/"
//	// 创建路径
//	err := os.MkdirAll(path, os.ModePerm)
//	if err != nil {
//		fmt.Println("无法创建路径:", err)
//		return "", "", nil, -1, 0, err
//	}
//	return path, filename, f, index, 1, nil
//}
//
//func WriteFileV2(f *excelize.File, index, i int, record, eliminate *utils.FileData) int {
//	i++
//	if record != nil {
//		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), record.FileName)
//		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), record.Summary)
//		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), "符合")
//		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), record.Eliminate)
//		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i), record.DataModel)
//
//	}
//	if eliminate != nil {
//		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), eliminate.FileName)
//		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), eliminate.Summary)
//		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), "不符合")
//		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), eliminate.Eliminate)
//		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i), eliminate.DataModel)
//	}
//	// 设置默认工作表
//	f.SetActiveSheet(index)
//	return i
//
//}
//
//func SaveFile(path, filename string, f *excelize.File) (string, string) {
//	// 保存文件
//	err := f.SaveAs(path + filename)
//	if err != nil {
//		fmt.Println("无法保存文件:", err)
//		return "", ""
//	}
//
//	fmt.Println("文件保存成功")
//	outputFilePath := path + filename
//	cfg := config.CustomCfg.Minio
//	info, err := db.MinioClient.FPutObject(context.Background(), cfg.BucketName, filename,
//		outputFilePath, minio.PutObjectOptions{ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"})
//	if err != nil {
//		return "", ""
//	}
//	//_ = os.Remove(outputFilePath)
//	slog.Debug("Successfully uploaded to minio", "fileName", filename, "endpoint", cfg.Endpoint, "Size", info.Size)
//
//	fileUrl := fmt.Sprintf("http://%s/%s/%s", cfg.Endpoint, cfg.BucketName, filename)
//	fmt.Println("文件地址:", fileUrl)
//	return fileUrl, filename
//}
//
//func WriteFileTest(testDatas []utils.TestData) {
//	// 创建一个新的XLSX文件
//	f := excelize.NewFile()
//	// 创建一个工作表
//	index, _ := f.NewSheet("Sheet1")
//	f.SetColWidth("Sheet1", "A", "A", 30)
//	f.SetColWidth("Sheet1", "B", "B", 50)
//	f.SetColWidth("Sheet1", "C", "C", 50)
//	f.SetColWidth("Sheet1", "D", "D", 200)
//
//	f.SetCellValue("Sheet1", "A1", "文件ID")
//	f.SetCellValue("Sheet1", "B1", "文件名称")
//	f.SetCellValue("Sheet1", "C1", "模型返回")
//	f.SetCellValue("Sheet1", "D1", "发送模型前的数据")
//
//	i := 1
//	for _, test := range testDatas {
//		i++
//
//		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), strconv.FormatInt(test.FileId, 10))
//		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), test.FileName)
//		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), test.ModelRes)
//		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), test.DataModel)
//	}
//
//	// 设置默认工作表
//	f.SetActiveSheet(index)
//
//	now := time.Now()
//	// 格式化日期和时间
//	date := now.Format("20060102")
//	hour := now.Format("15")
//	minute := now.Format("04")
//	// 构建文件名
//	filename := fmt.Sprintf("struct_%s-%s-%s.xlsx", date, hour, minute)
//	path := "./temp/"
//	// 创建路径
//	err := os.MkdirAll(path, os.ModePerm)
//	if err != nil {
//		fmt.Println("无法创建路径:", err)
//		return
//	}
//	// 保存文件
//	err = f.SaveAs(path + filename)
//	if err != nil {
//		fmt.Println("无法保存文件:", err)
//		return
//	}
//
//	fmt.Println("文件保存成功")
//}
//
//func WriteSummFile(summaryMap map[string][]utils.FileData) {
//	// 创建一个新的XLSX文件
//	f := excelize.NewFile()
//	// 创建一个工作表
//	index, _ := f.NewSheet("Sheet1")
//
//	f.SetColWidth("Sheet1", "A", "A", 50)
//	f.SetColWidth("Sheet1", "B", "B", 100)
//	f.SetColWidth("Sheet1", "C", "C", 100)
//	f.SetCellValue("Sheet1", "A1", "文件名")
//	f.SetCellValue("Sheet1", "B1", "总结")
//	f.SetCellValue("Sheet1", "C1", "发送模型前的数据")
//
//	i := 1
//	for fileName, v := range summaryMap {
//		for _, summary := range v {
//			i++
//			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), fileName)
//			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), summary.Summary)
//			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), summary.DataModel)
//		}
//	}
//
//	// 设置默认工作表
//	f.SetActiveSheet(index)
//
//	now := time.Now()
//	// 格式化日期和时间
//	date := now.Format("20060102")
//	hour := now.Format("15")
//	minute := now.Format("04")
//	// 构建文件名
//	filename := fmt.Sprintf("文章总结_%s-%s-%s.xlsx", date, hour, minute)
//	path := "./temp/"
//	// 创建路径
//	err := os.MkdirAll(path, os.ModePerm)
//	if err != nil {
//		fmt.Println("无法创建路径:", err)
//		return
//	}
//	// 保存文件
//	err = f.SaveAs(path + filename)
//	if err != nil {
//		fmt.Println("无法保存文件:", err)
//		return
//	}
//
//	fmt.Println("文件保存成功")
//}
//
//func ReadFile() map[string]map[string]string {
//	fileMap := make(map[string]map[string]string)
//
//	// 打开XLSX文件
//	file, err := xlsx.OpenFile("./files/template.xlsx")
//	if err != nil {
//		fmt.Println("打开文件失败:", err)
//		return nil
//	}
//
//	// 获取Sheet1
//	sheet, ok := file.Sheet["Sheet1"]
//	if !ok {
//		fmt.Println("找不到Sheet1")
//		return nil
//	}
//
//	// 遍历Sheet1中的所有行
//	for _, row := range sheet.Rows {
//		reasonMap := make(map[string]string)
//		// 遍历行中的所有单元格
//		fileName := row.Cells[0].String()
//		reason := row.Cells[2].String()
//		reasonMap[reason] = row.Cells[3].String()
//		fileMap[fileName] = reasonMap
//
//	}
//	return fileMap
//}
