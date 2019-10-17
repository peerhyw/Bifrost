package server

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/brokercap/Bifrost/plugin"
	"github.com/brokercap/Bifrost/config"
	"sync"
	"io"
)


var l sync.RWMutex

type recovery struct {
	Version string
	ToServer *json.RawMessage
	DbInfo *json.RawMessage
}

type recoveryDataSturct struct {
	Version string
	ToServer interface{}
	DbInfo interface{}
}

func DoRecoverySnapshotData(){

	var DataFile string = config.DataDir+"/db.Bifrost"
	//DataTmpFile = dataDir+"/db.Bifrost.tmp"

	fi, err := os.Open(DataFile)
	if err != nil {
		return
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)

	if err != nil {
		return
	}
	if string(fd) == ""{
		return
	}
	var data recovery
	errors := json.Unmarshal(fd,&data)
	if errors != nil{
		log.Printf("recovery error:%s, data:%s \r\n",errors,string(fd))
		return
	}
	if string(*data.ToServer) != "{}"{
		plugin.Recovery(data.ToServer)
	}
	if string(*data.DbInfo) != "{}"{
		Recovery(data.DbInfo,false)
	}
}

func GetSnapshotData() []byte{
	l.Lock()
	defer func(){
		l.Unlock()
		if err :=recover();err!=nil{
			log.Println(err)
		}
	}()
	data := recoveryDataSturct{
		Version:config.VERSION,
		ToServer:plugin.SaveToServerData(),
		DbInfo:SaveDBInfoToFileData(),
	}
	b,_:= json.Marshal(data)
	return b
}


func DoSaveSnapshotData(){
	var DataFile string = config.DataDir+"/db.Bifrost"
	var DataTmpFile string = config.DataDir+"/db.Bifrost.tmp"

	b := GetSnapshotData()

	f, err2 := os.OpenFile(DataTmpFile, os.O_CREATE|os.O_RDWR, 0700) //打开文件
	if err2 !=nil{
		log.Println("open file error:",err2)
		return
	}
	_, err1 := io.WriteString(f, string(b)) //写入文件(字符串)
	if err1 != nil {
		f.Close()
		log.Printf("save data to file error:%s, data:%s \r\n",err1,string(b))
		return
	}
	f.Close()
	err := os.Rename(DataTmpFile,DataFile)
	if err != nil{
		log.Println("doSaveDbInfo os.Rename err:",err)
	}
}


func DoRecoveryByBackupData(fileContent string){
	var data recovery
	errors := json.Unmarshal([]byte(fileContent),&data)
	if errors != nil{
		log.Printf("recovery error:%s, data:%s \r\n",errors,fileContent)
		return
	}
	if string(*data.ToServer) != "{}"{
		plugin.Recovery(data.ToServer)
	}
	if string(*data.DbInfo) != "{}"{
		Recovery(data.DbInfo,true)
	}
}