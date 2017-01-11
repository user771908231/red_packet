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
)

const (
	APP_ASSET_FILE="/Users/kory/Documents/Dev/cocos2d-x-3.12/casino/DDZ/assets/resources/HotUpdate/AssetsInfo.dat"
	BUILD_NATIVE_PATH = "/Users/kory/Documents/Dev/cocos2d-x-3.12/casino/DDZ/build_native/"
	ROOT_PATH = BUILD_NATIVE_PATH + "jsb-default/"
	OUTPUT_PATH = BUILD_NATIVE_PATH + "hotupdate/"

	FILEID_LIST_JSON = OUTPUT_PATH+"/FileIdList.json"

	//ASSET_HOST = "http://d.tondeen.com/hotupdate/"
	ASSET_HOST = "http://d.tondeen.com/testhot/"
	CLIENT_APPID = "1"
)


var (
	gFileIdMap = make( map[string] string)

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

func makeInfoFile(assets []*ddproto.AssetInfo, saveFile string, assetVer int32, redisHost string) (result bool) {
	pkgData := new(ddproto.HotupdateAckAssetsInfo) //这里是全量的资源，不需要返回全量的，只需要返回变化的就行了...
	pkgData.AssetHost = proto.String(ASSET_HOST)

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
		return false
	}


	//写入文件
	err = ioutil.WriteFile(saveFile, filedata, 0666)
	if err != nil {
		log.Printf("write文件失败:%v savefile:%V", err, saveFile)
		panic(err)
	}

	log.Println("=====AssetsInfo成功保存至文件:%v pkgData:%v", saveFile, pkgData)

	//同时保存一份至App资源目录
	err = ioutil.WriteFile(APP_ASSET_FILE, filedata, 0666)
	if err != nil {
		log.Printf("write文件失败:%v savefile:%V", err, APP_ASSET_FILE)
		panic(err)
	}

	//写入redis
	data.InitRedis(redisHost, "test")
	redisUtils.SetObj("AssetsInfo"+CLIENT_APPID, pkgData)

	//data.InitRedis("192.168.199.120:6379", "test")
	//redisUtils.SetObj("AssetsInfo"+CLIENT_APPID, pkgData)

	return true
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

func printAssetInfoFile(saveFile string ) {
	assetInfo := loadAssetInfoFromFile( saveFile )

	text := ""
	text += fmt.Sprintf("\n=========================\nassetInfo.dat读取文件后打印: \n\t [ assetInfo: %v ]\n \t[ 资源文件数:%v]\t[版本号:%v]\n",
	*assetInfo.AssetHost, len(assetInfo.Assets), *assetInfo.LastestAssetsVersion)
	for i, asset := range assetInfo.Assets {
		text += fmt.Sprintf("\t--[%d] asset >>> fid:%v fPath:%v fver:%v size:%v md5:%v gameId:%v isCode:%v compress:%v\n",
			i,  *asset.FileId, *asset.FilePath,
			*asset.FileVer, *asset.FileSize, *asset.Md5, *asset.GameId, *asset.IsCode, asset.IsCompress)
	}

	log.Printf(text)

	//写入文件
	printFile := filepath.Dir(saveFile)+"/print.txt"
	err := ioutil.WriteFile(printFile, []byte(text), 0666)
	if err != nil {
		log.Printf("write文件失败:%v savefile:%v", err, printFile)
		panic(err)
	}

	log.Printf("=========================\n")
}

func getFileVer(newAsset *ddproto.AssetInfo, oldAssetInfo *ddproto.HotupdateAckAssetsInfo) (fileVer *int32) {
	fileVer = new(int32)
	*fileVer = 1

	if oldAssetInfo == nil {
		return fileVer
	}

	//TODO: 到assetInfo.dat中找到asset.FileId对应的oldAsset.FileVer
	for _, asset := range oldAssetInfo.GetAssets() {
		if ( *asset.FileId == *newAsset.FileId ) {
			if ( *asset.FilePath == *newAsset.FilePath ) {
				if ( *asset.Md5 == *newAsset.Md5 ) {
					//md5一致, 直接返回旧文件的FileVer
					*fileVer = *asset.FileVer
					log.Printf("fid:%v[%v] md5未变,直接返回fileVer:%v", *asset.FileId, *asset.FilePath, *fileVer)
				} else {
					//新文件md5变了, 版本号+1
					*fileVer = *asset.FileVer + 1
					log.Printf("fid:%v[%v] > new fid:%v[%v]新文件md5变了,版本号+1=%v  md5:%v -> %v\n",
						*asset.FileId, *asset.FilePath, *newAsset.FileId, *newAsset.FilePath, *fileVer, *asset.Md5, *newAsset.Md5)
				}
			} else {
				*fileVer = *asset.FileVer
				log.Printf("非法数据: fileId相同但filePath不同: fid:%v old:%v new:%v\n", *asset.FileId, *asset.FilePath, *newAsset.FilePath )
				//panic(nil)
			}
			break
		}
	}

	return fileVer
}

//将多个散列文件打成1个zip包
func packSomeFiles(origAssetInfo *ddproto.HotupdateAckAssetsInfo, files[]string, module, outputPath,filePath,basePath string,  fid *int32 ) (assets *ddproto.AssetInfo, err error) {
	destFile := ""
	if filePath == "" {
		filePath = module + "/others.zip"
	}

	destFile = outputPath + filePath

	log.Printf(" ==== packSomeFiles >>>> files count:%v module:%v fid:%v filePath:%v, destFile:%v", len(files), module, fid, filePath, destFile )

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

	//计算文件md5
	md5str, _ := getMd5( destFile )

	isCode := false

	gameId := ddproto.CommonEnumGame_GID_HALL
	if( module == "Common" || module == "Hall" ) {
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
	} else if( strings.Contains(module, "src") ) { //源码
		//gameId = ddproto.CommonEnumGame_GID_SRC
		gameId = ddproto.CommonEnumGame_GID_HALL
		isCode = true
	} else if( strings.Contains(module, "import") ) { //
		gameId = ddproto.CommonEnumGame_GID_HALL
		isCode = true
	}

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

	isCode := false

	gameId := ddproto.CommonEnumGame_GID_HALL
	if( module == "Common" || module == "Hall" ) {
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
	} else if( strings.Contains(module, "src") ) { //源码
		//gameId = ddproto.CommonEnumGame_GID_SRC
		gameId = ddproto.CommonEnumGame_GID_HALL
		isCode = true
	} else if( strings.Contains(module, "import") ) { //
		gameId = ddproto.CommonEnumGame_GID_HALL
		isCode = true
	}

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

func packResources(importpath string, outputPath string,  isRelease, isOnlySource bool) (assets []*ddproto.AssetInfo, err error ) {
	resPath := ROOT_PATH + "/res/raw-assets/resources/"

	//读取上一次生成的资源信息
	origAssetInfo := loadAssetInfoFromFile( OUTPUT_PATH + "/AssetsInfo.dat" )


	//资源目录
	var dirs []string

	if isRelease {
		dirs = append(dirs, "src/project.jsc")
	} else {
		dirs = append(dirs, "src/project.dev.js")
	}

	if isOnlySource == false {
		resdirs, _ := GetDirs( resPath )
		dirs = append(dirs, resdirs...)
	}

	fileId := int32(0)
	idx := 0
	for _, module := range dirs {
		if( strings.Contains(module, ".") && !strings.Contains(module, "src") ) {
			log.Println("===skip module: ", module)
			continue
		}

		idx ++
		log.Printf("---idx[%v]make module: %v \n", idx, module)

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
func setAssetsFileToRedis(assetFile, clientAppId, redisHost string ) {
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

	err = redisUtils.SetObj("AssetsInfo"+clientAppId, assetInfo)
	if err!= nil {
		log.Printf("===== 写入redis失败: %v\n", err)
	} else {
		log.Printf("===== 写入redis成功! (key:%v) =====\n", "AssetsInfo"+clientAppId)
	}
}

func compareAssetInfo() {
	fileOld:="/Users/kory/Documents/Dev/cocos2d-x-3.12/casino/DDZ/build_native/hotupdate103/assetsInfo.dat"
	fileNew:="/Users/kory/Documents/Dev/cocos2d-x-3.12/casino/DDZ/build_native/hotupdate106/assetsInfo.dat"

	infoOld := loadAssetInfoFromFile( fileOld )
	infoNew := loadAssetInfoFromFile( fileNew )

	md5matchCnt := 0
	matchSize := int32(0)
	totalSize := int32(0)
	for i, assetNew := range infoNew.Assets {
		log.Printf("[%d] ==== loop >>> fid:%v fpath:%v md5: %v\n", i, *assetNew.FileId, *assetNew.FilePath, *assetNew.Md5)

		totalSize += *assetNew.FileSize
		bFound := false
		for _, assetOld := range infoOld.Assets {
			if *assetNew.FilePath == *assetOld.FilePath {
				if *assetNew.Md5 == *assetOld.Md5 {
					log.Printf("\t[%d] md5 is MATCH >>> fid:%v fpath:%v md5: %v\n", i, *assetNew.FileId, *assetNew.FilePath, *assetNew.Md5)
					md5matchCnt++
					matchSize += *assetNew.FileSize
				}
				bFound = true
				break
			}
		}

		if !bFound {
			log.Printf("\t[%d] ==== new File >>> fid:%v fpath:%v md5: %v\n", i, *assetNew.FileId, *assetNew.FilePath, *assetNew.Md5)
		}
	}


	log.Printf("md5match:%v totalCnt: %v  size:%v / %v (%v%%)\n", md5matchCnt, len(infoNew.Assets), matchSize, totalSize, (100*matchSize/totalSize) )

}

func saveFileIdList(assets []*ddproto.AssetInfo) bool {
	if isFileExist( FILEID_LIST_JSON ) {
		log.Printf("%v fileId文件已存在. ", FILEID_LIST_JSON)
		return false
	}

	gFileIdMap = make( map[string] string)
	for _, asset := range assets {
		//fileIdMap[*asset.FilePath] = *asset.FileId
		gFileIdMap[*asset.FilePath] = fmt.Sprintf("%d", *asset.FileId)
	}

	data, err := json.Marshal( gFileIdMap )
	if err != nil {
		log.Printf("json.Marshal err:%v\n",  err)
		return false
	}

	err = ioutil.WriteFile(FILEID_LIST_JSON, []byte( data ), 0666)
	if err != nil {
		log.Printf("write文件失败:%v savefile:%v", err, FILEID_LIST_JSON)
		panic(err)
	}

	log.Printf("=====保存文件Id信息完成( FileIdList.json )=====\n")

	return true
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

func main() {
	//compareAssetInfo()
	//return

	saveFile := OUTPUT_PATH + "/AssetsInfo.dat"
	//redisHost := "192.168.199.120:6379"
	redisHost := "127.0.0.1:6379"

//更新redis数据
	//redisHost = "127.0.0.1:6379"
	//saveFile = "./AssetsInfo.dat"
	//setAssetsFileToRedis( saveFile, redisHost, CLIENT_APPID)
	//return

	isRelease := true
	isOnlySource := false

	assets, err := packResources("", OUTPUT_PATH, isRelease, isOnlySource)
	if err != nil {
		return
	}

	//var assets []*ddproto.AssetInfo
	//全局资源版本号
	lastestAssetsVer := int32(1)
	makeInfoFile(assets, saveFile, lastestAssetsVer, redisHost)


	//保存文件Id
	saveFileIdList(assets)


	printAssetInfoFile( saveFile )
}

