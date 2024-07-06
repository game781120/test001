package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// Config 定义了YAML文件中的配置结构
type Config struct {
	Server struct {
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Model struct {
		LlmHost      string `yaml:"llmHost"`
		LlmAuthToken string `yaml:"llmAuthToken"`
		IsLLm        bool   `yaml:"isLLm"`
		Name         string `yaml:"name"`
	} `yaml:"model"`
	Email struct {
		Imap            string `yaml:"imap"`
		ImapPort        string `yaml:"imapPort"`
		Smtp            string `yaml:"smtp"`
		SmtpPort        string `yaml:"smtpPort"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		IntervalTime    string `yaml:"intervalTime"`
		CheckTime       string `yaml:"checkTime"`
		ReAnalysis      bool   `yaml:"reAnalysis"`
		Uuids           string `yaml:"uuids"`
		IsConfirm       bool   `yaml:"isConfirm"`
		ReNotify        int    `yaml:"reNotify"`
		DeleteAbstract  string `yaml:"deleteAbstract"`
		ImportantEmails string `yaml:"importantEmails"`
		SpamEmails      string `yaml:"spamEmails"`
		DeleteCount     int    `yaml:"deleteCount"`
		FinancialFile   string `yaml:"financialFile"`
		FinancialFileEn string `yaml:"financialFileEn"`
	} `yaml:"email"`

	Elastic struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Index    string `yaml:"index"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"elastic"`

	Embedding struct {
		Address       string `yaml:"address"`
		ModelName     string `yaml:"model_name"`
		Authorization string `yaml:"authorization"`
	} `yaml:"embedding"`

	Milvus struct {
		Address        string `yaml:"address"`
		DbName         string `yaml:"dbName"`
		CollectionName string `yaml:"collectionName"`
		Dim            int    `yaml:"dim"`
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		DatasetIds     string `yaml:"datasetIds"`
	} `yaml:"milvus"`
	FlyBook struct {
		AppId           string `yaml:"appId"`
		AppSecret       string `yaml:"appSecret"`
		TokenUrl        string `yaml:"tokenUrl"`
		BatchUrl        string `yaml:"batchUrl"`
		UserUrl         string `yaml:"userUrl"`
		UserSearchUrl   string `yaml:"userSearchUrl"`
		DepartmentIdUrl string `yaml:"departmentIdUrl"`
		MessagesUrl     string `yaml:"messagesUrl"`
		UserTokenUrl    string `yaml:"userTokenUrl"`
		DepartmentUrl   string `yaml:"departmentUrl"`
	} `yaml:"flyBook"`

	Nacos struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"nacos"`
	Nlu struct {
		Host string `yaml:"host"`
	} `yaml:"nlu"`
	Mysql struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbName"`
	} `yaml:"mysql"`

	JuShi struct {
		Url           string `yaml:"url"`
		ModelName     string `yaml:"model_name"`
		Authorization string `yaml:"authorization"`
	} `yaml:"juShi"`
}

// UpdateData 定义了updateData.yaml文件中需要实时更新的数据
type UpdateData struct {
	Email struct {
		LastTimestamp string `yaml:"lastTimestamp"`
	} `yaml:"email"`
	Uuids struct {
		LastUuid string `yaml:"lastUuid"`
	} `yaml:"uuids"`
}

var ConfigInfo = Config{}
var UpdateDataInfo = UpdateData{}

func LoadYaml() bool {
	dataBytes, err := os.ReadFile("conf/conf.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return false
	}
	fmt.Println("yaml 文件的内容: \n", string(dataBytes))

	err = yaml.Unmarshal(dataBytes, &ConfigInfo)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return false
	}

	dataBytes, err = os.ReadFile("conf/updateData.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return false
	}
	fmt.Println("yaml 文件的内容: \n", string(dataBytes))

	err = yaml.Unmarshal(dataBytes, &UpdateDataInfo)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return false
	}
	fmt.Println("config内容为:", ConfigInfo)
	fmt.Println("UpdateDataInfo:", UpdateDataInfo)

	return true

}

func UpdateYaml(updateDataInfo *UpdateData) bool {
	// 将更新后的Config结构体转换为YAML格式
	updatedYamlData, err := yaml.Marshal(*updateDataInfo)
	if err != nil {
		fmt.Printf("无法转换为YAML格式：%s\n", err)
		return false
	}
	// 将更新后的YAML数据写入文件
	err = os.WriteFile("conf/updateData.yaml", updatedYamlData, 0644)
	if err != nil {
		fmt.Printf("无法写入文件：%s\n", err)
		return false
	}
	fmt.Println("YAML文件更新成功！")
	return true

}
