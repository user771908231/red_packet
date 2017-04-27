package main

import (
	"log"
	"strings"
	"path/filepath"
	"os"
	"archive/zip"
	"io/ioutil"
	"casino_common/proto/ddproto"
	"casino_common/utils/redis"
	"casino_common/utils/redisUtils"
	"github.com/golang/protobuf/proto"
	"sort"
	"fmt"
	"crypto/md5"
	"io"
	"encoding/json"
	"strconv"
	"time"
	//"runtime/debug"
)

const (
)


var (
	gFileIdMap = make( map[string] string)
	gNewFiles = make( map[string] string)

	PROJ_ROOT_PATH = "/Users/kory/Documents/Dev/workspace/Git/casino/DDZ/"
	//PROJ_ROOT_PATH = "/Users/kory/Documents/Dev/cocos2d-x-3.12/casino/casino/DDZ/"
	BUILD_NATIVE_PATH = PROJ_ROOT_PATH+"build_native/"

	APP_ASSET_FILE=PROJ_ROOT_PATH+"assets/resources/HotUpdate/AssetsInfo.dat"
	IOS_BUILD_ASSET_FILE=BUILD_NATIVE_PATH+"jsb-default/res/raw-assets/resources/HotUpdate/AssetsInfo.dat"

	ROOT_PATH = BUILD_NATIVE_PATH + "jsb-default/"

	//OUTPUT_PATH = "/Users/kory/Documents/Dev/workspace/Git/GameUpdate/ios_lyqmj_3/" //iOS AppStore channelid=3
	OUTPUT_ROOT = "/Users/kory/Documents/Dev/workspace/Git/GameUpdate/"
	OUTPUT_PATH = OUTPUT_ROOT + "android_hunan_32/" //channelid=32

	AssetsVer = int32(1)
	FILEID_LIST_JSON = OUTPUT_PATH + "FileIdList_v1.json"
	FILEID_LIST_JSON_NEW = OUTPUT_PATH + "FileIdList_v2.json"


)


type FileIdInfo struct {
	FileId           int32
	FilePath         string
}

func isDir(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func getMd5(fileName string ) (result string, err error) {
	file, inerr := os.Open( fileName )
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)

		result = fmt.Sprintf("%x", md5h.Sum([]byte("")))
		//log.Printf("%s => result md5: %s\n", fileName, result) //md5
		return result, nil
	}

	return "", nil
}


//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}

		if( len(suffix) > 0 ) {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				files = append(files, filename)
			}
		} else {
			files = append(files, filename)
		}

		return nil
	})
	return files, err
}


//获取指定目录下的目录名(不递归子目录)
func GetDirs(dirPth string) (files []string, err error) {
	f, err := os.Open(dirPth)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func innerZipdir(path string, wholeDir string, w *zip.Writer) {

	//log.Printf("===== innerZipdir >> path:%v  wholeDir: %v\n", path, wholeDir)

	// 打开文件夹
	dir, err := os.Open(path)
	if err != nil {
		panic(nil)
	}
	defer dir.Close()

	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// 遍历文件列表
	for _, fi := range fis {
		// 递归
		if fi.IsDir() {
			//log.Printf("call innerZipdir >>> wholeDir:%v srcFile=%v nextWhole:%v\n", wholeDir, dir.Name() + "/" + fi.Name(), wholeDir+"/"+fi.Name() )
			innerZipdir(dir.Name() + "/" + fi.Name(), wholeDir+"/"+fi.Name(), w)
			continue
		}

		// 打印文件名称
		// log.Println(fi.Name())

		// 打开文件
		fr, err := os.Open(dir.Name() + "/" + fi.Name())
		if err != nil {
			panic(err)
		}
		defer fr.Close()
		fd,err := ioutil.ReadAll(fr)
		//log.Println(string(fd))
		filename := fi.Name()
		if wholeDir!="" {
			filename = wholeDir+"/"+fi.Name()
		}
		//log.Println("filename="+filename)

		f, err := w.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(fd))
		if err != nil {
			log.Fatal(err)
		}

	}
}

//压缩多个文件/目录
func zipFiles(srcFiles []string, basePath string, dstZip string) {
	fw, err := os.Create( dstZip )
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	w:=zip.NewWriter(fw)

	for _, srcFile := range srcFiles {
		//log.Printf("[%d] zipFiles >>> srcFile: %v base(srcFile):%v\n", i, srcFile, filepath.Base(srcFile))

		if isDir(srcFile) { // 目录
			realBasePath := basePath
			if basePath == "import" {
				realBasePath = filepath.Base(srcFile)
			}

			innerZipdir(srcFile, realBasePath, w)

		} else { // 文件
			fr, err := os.Open( srcFile )
			if err != nil {
				panic(err)
			}
			defer fr.Close()

			fd,err := ioutil.ReadAll(fr)
			filename := basePath + filepath.Base( srcFile )

			f, err := w.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.Write([]byte(fd))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	defer w.Close()
}

func zipFile(srcFile string, baseName, dstZip string) {
	log.Printf("---zipFile  >>> srcFile:%v dstZip:%v baseName:%v\n", srcFile, dstZip, baseName)

	os.MkdirAll(filepath.Dir(dstZip), 0777)

	fw, err := os.Create( dstZip )
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	w:=zip.NewWriter(fw)


	// 打开文件
	fr, err := os.Open( srcFile )
	if err != nil {
		panic(err)
	}
	defer fr.Close()
	fd,err := ioutil.ReadAll(fr)
	//log.Println(string(fd))
	filename := baseName + filepath.Base( srcFile )

	//log.Println("filename="+filename)

	f, err := w.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write([]byte(fd))
	if err != nil {
		log.Fatal(err)
	}

	defer w.Close()
}

func zipDir(srcPath, basePath, dstZip string) {
	fw, err := os.Create( dstZip )
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	w:=zip.NewWriter(fw)
	innerZipdir(srcPath, basePath, w)

	defer w.Close()
}

func makeInfoFile(isUpdateAppAsset bool, assets []*ddproto.AssetInfo, assetHost,
saveFile string, assetVer int32, redisHost string) (pkgData *ddproto.HotupdateAckAssetsInfo) {
	pkgData = new(ddproto.HotupdateAckAssetsInfo) //这里是全量的资源，不需要返回全量的，只需要返回变化的就行了...
	pkgData.AssetHost = proto.String(assetHost)

	pkgData.Assets = assets
	//copy(pkgData.Assets, assets)
	pkgData.LastestAssetsVersion = &assetVer

	for _, a := range pkgData.Assets {
		log.Println("-----a:%v", a)
	}

	//序列化
	filedata, err := proto.Marshal( pkgData )
	if err != nil {
		log.Println("error:", err)
		return nil
	}


	//if isFileExist(saveFile) {
	//	//备份文件
	//	os.Rename(saveFile, saveFile+".bak-"+time.Now().Format("2006-01-02 15:04:05"))
	//}

	//写入文件
	err = ioutil.WriteFile(saveFile, filedata, 0666)
	if err != nil {
		log.Printf("write文件失败:%v savefile:%V", err, saveFile)
		panic(err)
	}

	log.Println("=====AssetsInfo成功保存至文件:%v", saveFile)

	if ( isUpdateAppAsset ) {
		//同时保存一份至App资源目录
		err = ioutil.WriteFile(APP_ASSET_FILE, filedata, 0666)
		if err != nil {
			log.Printf("write文件失败:%v savefile:%V", err, APP_ASSET_FILE)
			panic(err)
		}

		//再保存一份至ios build 目录
		err = ioutil.WriteFile(IOS_BUILD_ASSET_FILE, filedata, 0666)
		if err != nil {
			log.Printf("write文件失败:%v savefile:%V", err, APP_ASSET_FILE)
			panic(err)
		}

		log.Printf("=== 已更新 App/Resource/AssetsInfo.dat ====\n")
	} else {
		log.Printf("===isUpdateAppAsset=false >> 无需更新 App/Resource/HotUpdate/AssetsInfo.dat\n")
	}


	return pkgData
}

func loadAssetInfoFromFile( saveFile string ) (assetInfo *ddproto.HotupdateAckAssetsInfo) {
	assetInfo = new(ddproto.HotupdateAckAssetsInfo)

	fr, err := os.Open(saveFile)
	if err != nil {
		log.Printf("Cannot open '%v', err:%v", saveFile, err)
		return nil
	}
	defer fr.Close()
	fileData, err := ioutil.ReadAll(fr)
	if err!= nil {
		log.Printf("read file err:%v",err)
		return nil
	}

	err = proto.Unmarshal( fileData, assetInfo )
	if err != nil {
		log.Printf("unmarshal err:%V",err)
		return nil
	}

	return assetInfo
}

func printAssetInfoFile(saveFile, logSaveFileId string ) string {
	assetInfo := loadAssetInfoFromFile( saveFile )

	text := "" + logSaveFileId
	text += fmt.Sprintf("\n=========================\nassetInfo.dat读取文件后打印: [%s] \n\t [ assetInfo: %v ]\n \t[ 资源文件数:%v]\t [上一版本:%v] [当前版本:%v]\n",
		time.Now().Format("2006-01-02 15:04:05"), *assetInfo.AssetHost, len(assetInfo.Assets), ConfData.LastVersion, *assetInfo.LastestAssetsVersion)
	for i, asset := range assetInfo.Assets {
		_, isNewFile := gNewFiles[*asset.FilePath]
		newFileTag := "      "
		if isNewFile {
			newFileTag = "[NEW] " //本版本跟上一版本比更新的文件
		}

		text += fmt.Sprintf("\t--%v[%d] asset >>> fid:%v fPath:%v fver:%v size:%v md5:%v gameId:%v isCode:%v\n",
			newFileTag, i,  *asset.FileId, *asset.FilePath,
			*asset.FileVer, *asset.FileSize, *asset.Md5, *asset.GameId, *asset.IsCode)
	}

	log.Printf(text)

	//写入文件
	outputName := "print_" + strings.Replace(filepath.Base(saveFile), ".dat", ".txt", -1)
	printFile := filepath.Dir(saveFile) + "/" + outputName
	err := ioutil.WriteFile(printFile, []byte(text), 0666)
	if err != nil {
		log.Printf("write文件失败:%v savefile:%v", err, printFile)
		panic(err)
	}

	log.Printf("============ %v 保存完毕=============\n", printFile )

	return text
}

func getFileVer(newAsset *ddproto.AssetInfo, oldAssetInfo *ddproto.HotupdateAckAssetsInfo) (fileVer *int32) {
	fileVer = new(int32)
	*fileVer = 1

	if oldAssetInfo == nil {
		log.Fatalf("未知错误：oldAssetInfo == nil")
		return fileVer
	}

	//TODO: 到assetInfo.dat中找到asset.FileId对应的oldAsset.FileVer
	isExistedFile := false
	for _, asset := range oldAssetInfo.GetAssets() {
		if ( *asset.FileId == *newAsset.FileId ) {
			if ( *asset.FilePath == *newAsset.FilePath ) {
				if ( *asset.Md5 == *newAsset.Md5 ) {
					//md5一致, 直接返回旧文件的FileVer
					*fileVer = *asset.FileVer
					//if *asset.FileId >= 82 {
					//	*fileVer = *asset.FileVer + 4
					//	log.Printf("  >>>>> fid:%v[%v] md5虽然未变, 但设置fileVer:%v", *asset.FileId, *asset.FilePath, *fileVer)
					//}

					//log.Printf("fid:%v[%v] md5未变,直接返回fileVer:%v", *asset.FileId, *asset.FilePath, *fileVer)
				} else {
					//新文件md5变了, 版本号+1
					*fileVer = *asset.FileVer + 1
					log.Printf("fid:%v[%v] > new fid:%v[%v]新文件md5变了,版本号+1=%v (size:%v -> %v) md5:%v -> %v\n",
						*asset.FileId, *asset.FilePath, *newAsset.FileId, *newAsset.FilePath, *fileVer,
						*asset.FileSize, *newAsset.FileSize, *asset.Md5, *newAsset.Md5)

					//记录已升级的文件列表
					_, ok := gNewFiles[*newAsset.FilePath]
					if !ok {
						gNewFiles[*newAsset.FilePath] = fmt.Sprintf("%d", *newAsset.FileId)
					}
				}

				//if *asset.FileId == 1 {
				//	*fileVer = 14
				//}
			} else {
				*fileVer = *asset.FileVer
				log.Fatalf("非法数据: fileId相同但filePath不同: fid:%v old:%v new:%v\n", *asset.FileId, *asset.FilePath, *newAsset.FilePath )
				//panic(nil)
			}

			isExistedFile = true //已存在fid的旧文件
			break
		}
	}

	if !isExistedFile {
		log.Printf("[新增文件]: fid:%v  path:%v", *newAsset.FileId,*newAsset.FilePath )
		_, ok := gNewFiles[*newAsset.FilePath]
		if !ok {
			gNewFiles[*newAsset.FilePath] = fmt.Sprintf("%d", *newAsset.FileId)
		}
	}

	return fileVer
}

func getGameId( module string ) (gameId ddproto.CommonEnumGame, isCode bool) {
	gameId = ddproto.CommonEnumGame_GID_HALL
	isCode = false
	//if( module == "Lobby" ) {
	//	gameId = 6 //朋友桌
	//} else if( module == "GameCommon" ) {
	//	gameId = 100 //游戏桌内公用资源
	//} else if( module == "Common" || module == "Hall" ) {
	//	gameId = ddproto.CommonEnumGame_GID_HALL
	//} else if( module == "Login" ) {
	//	gameId = ddproto.CommonEnumGame_GID_HALL
	//} else if( module == "DDZ" ) {
	//	gameId = ddproto.CommonEnumGame_GID_DDZ
	//} else if( module == "Mahjong" ) {
	//	gameId = ddproto.CommonEnumGame_GID_MAHJONG
	//} else if( module == "ZJH" ) {
	//	gameId = ddproto.CommonEnumGame_GID_ZJH
	//} else if( module == "PDK" ) {
	//	gameId = ddproto.CommonEnumGame_GID_PDK
	//}else if( strings.Contains(module, "src") ) { //源码
	//	//gameId = ddproto.CommonEnumGame_GID_SRC
	//	gameId = ddproto.CommonEnumGame_GID_HALL
	//	isCode = true
	//} else if( strings.Contains(module, "import") ) { //
	//	gameId = ddproto.CommonEnumGame_GID_HALL
	//	isCode = true
	//} else {
	//	log.Fatal("无法识别未知模块:" + module +". 终止进程!" )
	//}


	if( module == "Common" || module == "Hall" || module == "GameCommon" ) {
		gameId = ddproto.CommonEnumGame_GID_HALL
	} else if( module == "Lobby" ) {
		gameId = ddproto.CommonEnumGame_GID_HALL
	} else if( module == "Login" ) {
		gameId = ddproto.CommonEnumGame_GID_HALL
	} else if( module == "DDZ" ) {
		gameId = ddproto.CommonEnumGame_GID_DDZ
	} else if( module == "Mahjong" ) {
		gameId = ddproto.CommonEnumGame_GID_MAHJONG
	} else if( module == "ZJH" ) {
		gameId = ddproto.CommonEnumGame_GID_ZJH
	} else if( module == "PDK" ) {
		gameId = ddproto.CommonEnumGame_GID_PDK
	}else if( strings.Contains(module, "src") ) { //源码
		//gameId = ddproto.CommonEnumGame_GID_SRC
		gameId = ddproto.CommonEnumGame_GID_HALL
		isCode = true
	} else if( strings.Contains(module, "import") ) { //
		gameId = ddproto.CommonEnumGame_GID_HALL
		isCode = true
	}else {
		log.Fatal("无法识别未知模块:" + module +". 终止进程!" )
	}

	return gameId, isCode
}

//将多个散列文件打成1个zip包
func packSomeFiles(origAssetInfo *ddproto.HotupdateAckAssetsInfo, files[]string, module, outputPath,filePath,basePath string,  fid *int32 ) (assets *ddproto.AssetInfo, err error) {
	destFile := ""
	if filePath == "" {
		filePath = module + "/others.zip"
	}

	destFile = outputPath + filePath

	log.Printf(" ==== packSomeFiles >>>> files count:%v module:%v fid:%v filePath:%v, destFile:%v", len(files), module, *fid, filePath, destFile )

	//创建所需目录
	err = os.MkdirAll(filepath.Dir(destFile), 0777)


	//压缩生成zip包
	zipFiles(files, basePath, destFile)

	//获取文件大小
	fileSize := int32(0)
	fileInfo, err := os.Stat(destFile)
	if err != nil {
		panic(err)
	}
	fileSize = int32(fileInfo.Size())


	//纠正历史问题
	if filepath.Base(destFile) == "gameOverWindow.zip" {
		log.Printf("重命名: gameOverWindow.zip -> GameOverWindow.zip\n")
		os.Rename(destFile,  filepath.Dir(destFile) + "GameOverWindow.zip" )
	}

	//计算文件md5
	md5str, _ := getMd5( destFile )


	gameId, isCode := getGameId( module )


	//*fid = *fid + 1
	fileId := int32( getFileId(fid, &filePath) )

	asset := new( ddproto.AssetInfo )
	asset.FilePath = &filePath
	asset.FileId =  &fileId //文件id
	asset.FileSize = &fileSize
	asset.IsCode = &isCode
	asset.Md5 = &md5str
	asset.GameId = &gameId
	//asset.IsCompress = new(bool)
	//*asset.IsCompress = true
	asset.FileVer = getFileVer(  asset, origAssetInfo )  //文件版本

	log.Printf("===生成one asset >>> fid[%d]: %v\n", *asset.FileId, asset)

	return asset, nil
}

func packOneAsset(origAssetInfo *ddproto.HotupdateAckAssetsInfo,  resPath, module, outputPath, filePath string, fid *int32 ) (assets *ddproto.AssetInfo, err error) {
	srcPath := ROOT_PATH

	destFile := outputPath + filePath
	basePath := filepath.Base(resPath)

	if( strings.Contains(module, "src") ) { //"src/project.jsc"
		filePath = "src.zip"
		destFile = outputPath + "src.zip"
		basePath = "src/"
	}

	log.Printf(" ==== packOneAsset >>>> resPath:%v\n outputPath:%v filePath:%v\n module:%v basePath:%s fid:%v",
		resPath, outputPath, filePath, module, basePath, fid )

	//创建所需目录
	err = os.MkdirAll(filepath.Dir(destFile), 0777)

	//压缩生成zip包
	if( strings.Contains(module, "src") ) {
		zipFile(srcPath+module, basePath, destFile)
	} else {
		zipDir(resPath, basePath, destFile)
	}

	//获取文件大小
	fileSize := int32(0)
	fileInfo, err := os.Stat(destFile)
	if err != nil {
		panic(err)
	}
	fileSize = int32(fileInfo.Size())

	//计算文件md5
	md5str, _ := getMd5( destFile )

	gameId, isCode := getGameId( module )

	fileId := getFileId(fid, &filePath)

	asset := new( ddproto.AssetInfo )
	asset.FilePath = &filePath
	asset.FileId = &fileId //文件id
	asset.FileSize = &fileSize
	asset.IsCode = &isCode
	asset.Md5 = &md5str
	asset.GameId = &gameId
	//asset.IsCompress = new(bool)
	//*asset.IsCompress = true
	asset.FileVer = getFileVer( asset, origAssetInfo )  //文件版本

	log.Printf("===生成one asset >>> fid[%d]: %v\n", *asset.FileId, asset)

	return asset, nil
}

func packResources(importpath string, outputPath string, oldAssetFile string,  isOnlySource bool) (assets []*ddproto.AssetInfo, err error ) {
	resPath := ROOT_PATH + "/res/raw-assets/resources/"

	//读取上一次生成的资源信息
	origAssetInfo := loadAssetInfoFromFile( oldAssetFile )

	var srcFiles []string

	isRelease := isFileExist( ROOT_PATH+"src/project.jsc" )

	if isRelease {
		srcFiles = append(srcFiles, ROOT_PATH+"src/project.jsc")
		srcFiles = append(srcFiles, ROOT_PATH+"src/settings.jsc")
	} else {
		srcFiles = append(srcFiles, ROOT_PATH+"src/project.dev.js")
		srcFiles = append(srcFiles, ROOT_PATH+"src/settings.js")
	}

	fileId := int32(0)

	//先打包src源码
	basePath := "src/"
	filePath := "src.zip"
	asset, err := packSomeFiles(origAssetInfo, srcFiles, "src", outputPath, filePath, basePath, &fileId )
	if err == nil {
		assets = append( assets, asset)
	}

	if isOnlySource {
		return assets, nil
	}

	//资源目录
	var dirs []string

	if isOnlySource == false {
		resdirs, _ := GetDirs( resPath )
		dirs = append(dirs, resdirs...)
	}

	idx := 0
	for _, module := range dirs {
		if( strings.Contains(module, ".") && !strings.Contains(module, "src") ) {
			log.Println("===skip module: ", module)
			continue
		}

		idx ++
		log.Printf("---idx[%v]make module[一级目录]: %v \n", idx, module)

		if isDir( resPath + module ) {
			realOutputPath := outputPath + "res/raw-assets/resources/"

			subdirs, _ := GetDirs( resPath + module )
			singleFiles := make([]string, 0)
			for subidx, subdir := range subdirs {
				if( subdir == ".DS_Store" || subdir == ".git" || subdir=="." ||subdir=="..") {
					log.Printf("====skip subdir:'%v' ===", subdir)
					continue
				}

				//二级目录
				currFile := resPath + module + "/" + subdir
				if isDir( currFile ) {
					//fileId := int32(idx + subidx)
					//fileId ++
					log.Printf("\t---sub[%v] make sub module: fid:%v >>> %v/%v \n", subidx, fileId, module, subdir)

					filePath := module + "/" + subdir + ".zip"
					asset, err := packOneAsset(origAssetInfo, resPath + module + "/"+subdir+"/", module, realOutputPath, filePath, &fileId)
					if err == nil {
						assets = append(assets, asset)
					}
				} else { //是文件(散列的文件)
					if subdir=="AssetsInfo.dat" {
						log.Printf("===skip module: %v/%v, 跳过打包AssetsInfo.dat\n", module,subdir)
						continue
					}
					singleFiles = append(singleFiles, currFile)
				}
			}

			if( len(singleFiles) > 0) {
				//打包散列的文件
				filePath := ""
				basePath := ""
				asset, err := packSomeFiles(origAssetInfo, singleFiles, module, realOutputPath, filePath, basePath, &fileId )
				if err == nil {
					assets = append( assets, asset)
				}
			}

		} else { // 一级目录的单个文件( 是src )
			if strings.Contains(module, "src/project") { //module=="src/project.jsc"
				//fileId ++
				realOutputPath := outputPath
				log.Printf("单个文件打包: module=%v  fileId:%\n",  module, fileId)
				asset, err := packOneAsset(origAssetInfo, ROOT_PATH, module, realOutputPath, "", &fileId)
				if err == nil {
					assets = append( assets, asset)
				}
			} else {
				log.Printf("严重错误: 一级目录存在非法单独文件: %v", module)
				panic(new(error))
			}
		}
	}

	if isOnlySource {
		return assets, nil
	}

	///////////////////////////////////
	// 处理import资源

	importPath := ROOT_PATH + "res/import/"
	importDirs, _ := GetDirs( importPath )

	pathsMap := make(map[string] []string)
	for i, dir := range importDirs {
		if dir == "." || dir == ".." || dir==".DS_Store" {
			log.Printf("====skip dir:'%v' ===", dir)
			continue
		}

		if !isDir(importPath+dir) {
			log.Fatal("import目录下出现非法文件:%v", dir)
			panic(new(error))
		}

		key := string(dir[0])

		path := importPath + dir
		paths, ok := pathsMap[key]
		if !ok {
			paths := make([]string, 1)
			paths[0] = path
			pathsMap[key] = paths
		} else {
			pathsMap[key] = append(paths, path)
		}
		log.Printf("[%d] import path:%v", i, path)
	}


	module := "import"
	realOutputPath := outputPath

	log.Printf("=== import pathsMap : === \n")
	for key, files := range pathsMap {
		log.Printf("import prefix: %v >> ", key)
		//for i, file:= range files {
		//	log.Printf("\t [%v]: file:%v", i, file)
		//}

		filePath := "res/import/" + key + ".zip"
		basePath := "import"
		asset, err := packSomeFiles(origAssetInfo, files, module, realOutputPath, filePath, basePath, &fileId )
		if err == nil {
			assets = append( assets, asset)
		}
	}

	return assets, nil
}


//读取assetInfo文件,并写入redis
func setAssetsFileToRedis(assetFile, clientAppId, redisHost string, isGrey bool ) {
	//写入redis
	data.InitRedis(redisHost, "test")

	log.Printf("redisHost:%v assetFile:%v clientAppId:%v\n", redisHost, assetFile, clientAppId )

	fr, err := os.Open( assetFile )
	if err != nil {
		log.Printf("===== 读取文件失败: %v, err:%v\n", assetFile, err)
		panic(err)
	}
	defer fr.Close()
	fileData, err := ioutil.ReadAll(fr)
	if err!= nil {
		log.Printf("read file err:%v",err)
		return
	}

	assetInfo := new(ddproto.HotupdateAckAssetsInfo)
	err = proto.Unmarshal( fileData, assetInfo )
	if err != nil {
		log.Printf("unmarshal err:%V",err)
		return
	}

	key := "AssetsInfo"+clientAppId
	if isGrey { //灰度测试key
		key = "AssetsInfoWhite"+clientAppId
	}
	err = redisUtils.SetObj(key , assetInfo)
	if err!= nil {
		log.Printf("===== 写入redis失败: %v\n", err)
	} else {
		log.Printf("===== 写入redis成功! (key:%v) =====\n", key)
	}

	/////////////////
	//打印数据
	text := ""
	text += fmt.Sprintf("\n=========================\nassetInfo.dat读取文件后打印: [%s] \n\t [ assetInfo: %v ]\n \t[ 资源文件数:%v]\t[版本号:%v]\n",
		time.Now().Format("2006-01-02 15:04:05"), *assetInfo.AssetHost, len(assetInfo.Assets), *assetInfo.LastestAssetsVersion)

	for i, asset := range assetInfo.Assets {
		text += fmt.Sprintf("\t--[%d] asset >>> fid:%v fPath:%v fver:%v size:%v md5:%v gameId:%v isCode:%v compress:%v\n",
			i,  *asset.FileId, *asset.FilePath,
			*asset.FileVer, *asset.FileSize, *asset.Md5, *asset.GameId, *asset.IsCode, asset.IsCompress)
	}

	log.Printf(text)
}

func compareAssetInfo(fileOld, fileNew string) {
	log.Printf("比较两个文件：\n\t old:%v \n\t new:%v\n", fileOld, fileNew)
	infoOld := loadAssetInfoFromFile( fileOld )
	infoNew := loadAssetInfoFromFile( fileNew )

	md5matchCnt := 0
	badCnt := 0
	updatedCnt := 0
	newAddCnt := 0
	totalSize := int32(0)
	updateSize := int32(0)
	for i, assetNew := range infoNew.Assets {
		log.Printf("[%d] ==== loop assetNew >>> fid:%v ver:%v path:%v md5: %v gameId:%v\n",
			i, *assetNew.FileId, *assetNew.FileVer, *assetNew.FilePath, *assetNew.Md5, *assetNew.GameId)

		totalSize += *assetNew.FileSize
		bFound := false
		for _, assetOld := range infoOld.Assets {
			if *assetNew.FileId == *assetOld.FileId {
				if *assetNew.FilePath != *assetOld.FilePath {
					log.Printf("[%d]严重错误： fid:%v 对应的path不同 >> oldPath:%v newPath:%v GameId:%v\n",
						i, *assetNew.FileId , *assetOld.FilePath, *assetNew.FilePath, *assetNew.GameId)
					badCnt ++
				} else if *assetNew.FileVer == *assetOld.FileVer {
					if *assetNew.Md5 == *assetOld.Md5 {
						log.Printf("\t[%d] ==== fver 和 md5 均一致,合法记录 >>> fid:%v ver:%v fpath:%v md5: %v gameId:%v\n",
							i, *assetNew.FileId,*assetNew.FileVer, *assetNew.FilePath, *assetNew.Md5, *assetNew.GameId)
						md5matchCnt++
					} else {
						badCnt ++
						log.Printf("[%d]严重错误： fver相同 但 md5不同 >>> fid:%v fpath:%v oldMd5: %v newMd5:%v\n",
							i, *assetNew.FileId, *assetNew.FilePath, *assetOld.Md5, *assetNew.Md5)
					}
				} else {
					updatedCnt++
					updateSize += *assetNew.FileSize
					if *assetNew.FileVer < *assetOld.FileVer {
						log.Printf("\t[%d] [警告]fver有更新, 但newVer < oldVer. >>> fid:%v fpath:%v oldVer:%v newVer:%v, oldMd5: %v newMd5:%v\n",
							i, *assetNew.FileId, *assetNew.FilePath, *assetOld.FileVer, *assetNew.FileVer, *assetOld.Md5, *assetNew.Md5)
					} else {
						log.Printf("\t[%d] fver有更新 >>> fid:%v fpath:%v oldVer:%v newVer:%v, oldMd5: %v newMd5:%v\n",
							i, *assetNew.FileId, *assetNew.FilePath, *assetOld.FileVer, *assetNew.FileVer, *assetOld.Md5, *assetNew.Md5)
					}
				}

				bFound = true
				break
			}
		}

		if !bFound {
			newAddCnt++
			updateSize += *assetNew.FileSize
			log.Printf("\t[%d] ==== 新增文件 >>> fid:%v ver:%v size:%v path: %v md5: %v\n",
				i, *assetNew.FileId, *assetNew.FileVer, *assetNew.FileSize, *assetNew.FilePath, *assetNew.Md5)
		}
	}


	log.Printf("===== 总数: %v 未变动:%v 更新数:%v 新增数:%v 错误数:%v totalSize:%v updateSize:%v =====\n",
		len(infoNew.Assets), md5matchCnt, updatedCnt, newAddCnt, badCnt,  totalSize, updateSize)

}

func saveFileIdList( newAssetFile string ) (result bool, logstr string) {
	//if isFileExist( FILEID_LIST_JSON ) {
	//	log.Printf("%v fileId文件已存在. ", FILEID_LIST_JSON)
	//	return false
	//}
	log.Printf("saveFileIdList>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")
	logstr = ""
	logstr += fmt.Sprintf("==========检查是否有新增fileId=========\n")
	//debug.PrintStack()

	//outputFileidList := OUTPUT_PATH + "FileIdList" + "_v" + fmt.Sprintf("%d", assetsVer) // + ".new"
	var gFileIdMap map[string] string

	//读取新的assets文件
	newAssetInfo := loadAssetInfoFromFile( newAssetFile )

	//读取旧的FileidList文件
	fr, err := os.Open( FILEID_LIST_JSON )
	if err != nil {
		logstr += fmt.Sprintf("Cannot open '%v', err:%v", FILEID_LIST_JSON, err)
		panic(err)
	}
	defer fr.Close()
	fileData, err := ioutil.ReadAll(fr)
	if err!= nil {
		logstr += fmt.Sprintf("read file err:%v",err)
		panic(err)
	}
	err = json.Unmarshal( fileData, &gFileIdMap )
	if err != nil {
		logstr += fmt.Sprintf("json.Unmarshal err:%v\n",  err)
		panic(err)
	}

	//比较得出差异的fileId,添加进去
	for _, newasset := range newAssetInfo.Assets {
		_, ok := gFileIdMap[*newasset.FilePath]
		if !ok {
			logstr += fmt.Sprintf("====<<差异[新增]fileId:%v path:%v\n", *newasset.FileId, *newasset.FilePath)
			gFileIdMap[*newasset.FilePath] = fmt.Sprintf("%d", *newasset.FileId)
		}
	}
	//打印已删除的fileId
	for filePath, fileId := range gFileIdMap {
		isExists := false
		for _, newasset := range newAssetInfo.Assets {
			if filePath == *newasset.FilePath {
				isExists = true
				break
			}
		}
		if !isExists {
			logstr += fmt.Sprintf("====>>差异[已删除]fileId:%v path:%v\n", fileId, filePath)
		}
	}

	//序列化为json保存
	data, err := json.Marshal( gFileIdMap )
	if err != nil {
		log.Printf("json.Marshal err:%v\n",  err)
		panic(err)
	}

	err = ioutil.WriteFile(FILEID_LIST_JSON_NEW, []byte( data ), 0666)
	if err != nil {
		log.Printf("write文件失败:%v savefile:%v", err, FILEID_LIST_JSON_NEW)
		panic(err)
	}

	logstr += fmt.Sprintf("=====保存文件Id信息完成( %s )=====\n", FILEID_LIST_JSON_NEW)

	//log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")
	//log.Printf("logstr: %s\n", logstr)
	//log.Printf("===========>>>>>>>>>>>>>>>>==========\n")

	return true, logstr
}

func getFileId(globalFid *int32, filePath *string) (fileId int32) {

	if len(gFileIdMap)==0 && isFileExist( FILEID_LIST_JSON ) {
		fr, err := os.Open( FILEID_LIST_JSON )
		if err != nil {
			log.Printf("Cannot open '%v', err:%v", FILEID_LIST_JSON, err)
			*globalFid = *globalFid + 1
			return *globalFid
		}
		defer fr.Close()
		fileData, err := ioutil.ReadAll(fr)
		if err!= nil {
			log.Printf("read file err:%v",err)
			*globalFid = *globalFid + 1
			return *globalFid
		}


		log.Printf("读取FileId文件成功:%v!", FILEID_LIST_JSON)

		//var fileIdMap = make( map[string] int32)
		err = json.Unmarshal( fileData, &gFileIdMap )
		if err != nil {
			log.Printf("===getFileId >> unmarshal err:%v",err)
			*globalFid = *globalFid + 1
			return *globalFid
		}

		maxFid := int32(0)
		for _, fid := range gFileIdMap {
			fileId, _ := strconv.ParseInt(fid, 10, 64)
			if int32(fileId) > maxFid {
				maxFid = int32(fileId)
			}
		}
		*globalFid = maxFid
		log.Printf("将globalFid设为maxFid: %v \n", maxFid)
	}

	fid, exists := gFileIdMap[ *filePath ]
	if exists {
		log.Printf("===getFileId >> 已找到文件[%v]匹配的fileId: %v\n", *filePath, fid)
		fileId, _ := strconv.ParseInt(fid, 10, 64)
		//fileId = fid
		return int32(fileId)
	} else {
		log.Printf("===【警告】:getFileId >> 未找到文件[%v]匹配的fileId, 取用globalId+1: %v\n", *filePath, *globalFid+1)
		*globalFid = *globalFid + 1
		return *globalFid
	}

	*globalFid = *globalFid + 1
	return *globalFid
}

var ConfData struct {
	ChannelId	   int32
	CurrVersion	   int32
	LastVersion        int32
	RedisAddr          string //redis配置
	ProjectPath	   string //
	OutputPath         string
	AssetHost	   string

}

//载入配置
func LoadJsonConfig(cid string) {
	filename := "conf.json"
	if cid != "" {
		filename = "./conf."+cid+".json"
	}
	data, err := ioutil.ReadFile( filename )
	if err != nil {
		log.Fatal("%v", err)
	}
	//log.Printf("配置文件信息:%v\n", string(data))

	err = json.Unmarshal(data, &ConfData)
	if err != nil {
		log.Fatal("%v", err)
	}

	log.Printf("========================\n")
	log.Printf("配置信息: \nChannelId: %v\nCurrVersion: %v\nLastVersion:%v\nRedisAddr: %v\nProjectPath:%v\nOUTPUT_PATH:%v\nAssetHost: %v\n",
		ConfData.ChannelId, ConfData.CurrVersion, ConfData.LastVersion, ConfData.RedisAddr,
		ConfData.ProjectPath, ConfData.OutputPath,  ConfData.AssetHost)
	log.Printf("========================\n")
}



func main() {
	if len(os.Args) < 2 {
		fmt.Printf("用法1： 输入 channelId 打包资源，生成 AssetsInfo. \n")
		fmt.Printf("用法2： 输入 -print {AssetsInfoFile} 来打印 \n")
		fmt.Printf("用法3： 输入 -compare {OldAssetsFile} {NewAssetsFile} 比较差异 \n")
		return
	}

	cid := os.Args[1]    //channelId
	assetHost := ""

	//=======================================
	//读取AssetsInfo.dat => 更新redis数据
	if len(os.Args) < 2 {
		fmt.Printf("用法： 输入AssetsInfo{cid}_{ver}.dat {gray/release} {redishost}\n")
		return
	}
	AssetFile:= os.Args[1]
	redisHost := "127.0.0.1:6379"
	mode := "正式发布"
	isGrey := false
	if len(os.Args) > 2 {
		isGrey = (os.Args[2]!="release") //(os.Args[2]=="gray" || os.Args[2]=="grey")
		if isGrey {
			mode = "灰度发布"
		}
	}
	if len(os.Args) > 3 {
		redisHost = os.Args[3]
	}
	log.Printf("【"+mode+"】 cid="+ cid + " AssetFile=" + AssetFile + " redisHost="+redisHost+"\n")

	//文件名格式如：AssetsInfo3_v23.dat (其中的23为cid）
	ss1 := strings.Split(AssetFile, "_")
	ss2 := strings.Split(ss1[1], ".")
	cid = ss2[0][1:]

	_, error := strconv.Atoi( cid )
	if error != nil{
		log.Fatalf("[%v] 非法文件名格式，无法识别cid：%v\n", AssetFile, cid)
	}

	//log.Printf("cid=%v ss1:%v ss2:%v\n", cid, ss1,ss2)

	setAssetsFileToRedis( AssetFile, cid, redisHost, isGrey )

	log.Printf("【"+mode+"】已导入 cid="+ cid + " AssetFile=" + AssetFile + " redisHost="+redisHost+"\n")

	return
	//=======================================

	if os.Args[1] == "-print" {
		toPrintAssetFile :=  os.Args[2]

		log.Printf(">>>>> 即将读取打印：%v\n", toPrintAssetFile)

		printAssetInfoFile( toPrintAssetFile, "" )
		return
	} else if os.Args[1] == "-printsaveid" {
		toPrintAssetFile :=  os.Args[2]

		FILEID_LIST_JSON = filepath.Dir(toPrintAssetFile) + "/FileIdList_Old.json"
		FILEID_LIST_JSON_NEW = filepath.Dir(toPrintAssetFile)  + "/FileIdList_Curr.json"
		if len(os.Args) > 3 {
			FILEID_LIST_JSON = filepath.Dir(toPrintAssetFile) + "/" + os.Args[3]
		}

		log.Printf(">>>>> 即将读取打印：%v\n基于FILEID_LIST_JSON: %v\n", toPrintAssetFile, FILEID_LIST_JSON)

		_, logSaveFileId := saveFileIdList( toPrintAssetFile )

		printAssetInfoFile( toPrintAssetFile, logSaveFileId )
		return
	} else if os.Args[1] == "-compare" { //比较两个版本的AssetsInfo.dat
		oldAssetFile := os.Args[2]
		newAssetFile := os.Args[3]
		compareAssetInfo( oldAssetFile, newAssetFile )

		return
	}

	//读取配置文件
	LoadJsonConfig( cid )

	PROJ_ROOT_PATH = ConfData.ProjectPath

	BUILD_NATIVE_PATH = PROJ_ROOT_PATH + "build_native/"
	APP_ASSET_FILE=PROJ_ROOT_PATH + "assets/resources/HotUpdate/AssetsInfo.dat"
	IOS_BUILD_ASSET_FILE=BUILD_NATIVE_PATH + "jsb-default/res/raw-assets/resources/HotUpdate/AssetsInfo.dat"
	ROOT_PATH = BUILD_NATIVE_PATH + "jsb-default/"


	//全局资源版本号
	//AssetsVer = int32( 18 ) //2017.04.11
	//AssetsVer = int32( 20 ) //2017.04.15
	//AssetsVer = int32( 21 ) //2017.04.16
	//AssetsVer = int32( 22 ) //2017.04.17
	AssetsVer = ConfData.CurrVersion //int32( 23 ) //2017.04.18

	if ConfData.LastVersion <= 0 {
		ConfData.LastVersion = AssetsVer - 1
	}

	sAssetsVer := fmt.Sprintf("%d", AssetsVer)
	sAssetsVerOld := fmt.Sprintf("%d", ConfData.LastVersion)
	//sAssetsVerOld = "19"

	// 从配置文件赋值
	OUTPUT_PATH = ConfData.OutputPath
	assetHost = ConfData.AssetHost

	////channelid=32
	//if cid == "32" {
	//	OUTPUT_PATH = OUTPUT_ROOT + "android_hunan_32/"
	//	assetHost = "http://test2.tondeen.com/update/android_hunan_32/" //android湖南: cid=32
	//} else if cid == "3" {
	//	OUTPUT_PATH = OUTPUT_ROOT +"ios_lyqmj_3/"
	//	assetHost = "http://test2.tondeen.com/update/lyqmj3/"  //来一圈麻将
	//} else if cid == "33" {
	//	OUTPUT_PATH = OUTPUT_ROOT +"android_gz_33/"
	//	assetHost = "http://test2.tondeen.com/update/android_gz_33/"  //广州
	//}


	oldAssetFile := OUTPUT_PATH + "AssetsInfo"+ cid +"_v"+ sAssetsVerOld + ".dat"
	newAssetFile := OUTPUT_PATH + "AssetsInfo"+ cid +"_v"+ sAssetsVer + ".dat"

	FILEID_LIST_JSON = OUTPUT_PATH + "FileIdList"+ cid  + "_v" + sAssetsVerOld + ".json"
	FILEID_LIST_JSON_NEW = OUTPUT_PATH + "FileIdList"+ cid + "_v" + sAssetsVer + ".json"

	if true {
		isOnlySource := false      //是否只生产源码
		redisHost := "127.0.0.1:6379"

		isUpdateAppAsset := true  //是否更新App内置信息:/Resource/HotUpdate/AssetsInfo.dat

		assets, err := packResources("", OUTPUT_PATH, oldAssetFile, isOnlySource)
		if err != nil {
			return
		}

		//var assets []*ddproto.AssetInfo
		//assetHost := ASSET_HOST
		//assetHost := "http://test2.tondeen.com/sjddz_hotupdate/"  //神经斗地主
		//cid = "31"
		//AssetsVer = int32( 11 ) //2017.03.26
		//assetHost := "http://test2.tondeen.com/corp_hotupdate/"  //企业号: cid=31

		pkgData := makeInfoFile(isUpdateAppAsset, assets, assetHost, newAssetFile, AssetsVer, redisHost)

		if( pkgData != nil ) {
			////写入redis
			//log.Printf("将打印saveFile:%s\n")
			//data.InitRedis(redisHost, "test")
			//redisUtils.SetObj("AssetsInfo"+CLIENT_APPID, pkgData)
		}

	}

	//从文件AssetsInfo.dat.new生成FileId
	_, logSaveFileId := saveFileIdList( newAssetFile )
	logSaveFileId += "========= 生成时间: " + time.Now().Format("2006-01-02 15:04:05") + " ==========\n" +logSaveFileId

	//log.Printf("将打印saveFile:%s\n", newAssetFile)
	//log.Printf( logSaveFileId )

	printAssetInfoFile( newAssetFile, logSaveFileId )

	//写入redis
	//setAssetsFileToRedis( newAssetFile, cid, "127.0.0.1:6379")

}

