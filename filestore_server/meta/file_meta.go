package meta

//FileMeta文件元信息格式
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UpLoadAt string
}

var fileMetas = map[string]FileMeta

func init() {
	fileMetas=make(map[string]FileMeta)
}

// UpdateFileMeta:更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1]=fmeta
}

// GetFileMeta:通过sha1值获取文件信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}