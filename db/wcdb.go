package db

import (
	"fmt"
	"strings"
)

type WCDB struct {
	enmicromsg  *EnMicroMsg
	wxfileindex *WxFileIndex
}

func InitWCDB(basePath string) *WCDB {
	wcdb := &WCDB{}
	wcdb.enmicromsg = OpenEnMicroMsg(basePath + "/EnMicroMsg_plain.db")
	wcdb.wxfileindex = OpenWxFileIndex(basePath + "/WxFileIndex_plain.db")
	return wcdb
}

func (wcdb WCDB) ChatList(pageIndex int, pageSize int, all bool) *ChatList {
	return wcdb.enmicromsg.ChatList(pageIndex, pageSize, all)
}

func (wcdb WCDB) ChatDetailList(talker string, pageIndex int, pageSize int) *ChatDetailList {
	result := wcdb.enmicromsg.ChatDetailList(talker, pageIndex, pageSize)
	detailList := make([]ChatDetailListRow, 0)
	for _, v := range result.Rows {
		detailList = append(detailList, wcdb.getMediaPath(v))
	}
	result.Rows = detailList
	return result
}

func (wcdb WCDB) GetUserInfo(username string) UserInfo {
	return wcdb.enmicromsg.GetUserInfo(username)
}

func (wcdb WCDB) GetMyInfo() UserInfo {
	return wcdb.enmicromsg.GetMyInfo()
}

func (wcdb WCDB) GetImgPath(msgId string) string {
	return wcdb.wxfileindex.GetImgPath(msgId)
}

func (wcdb WCDB) GetVideoPath(msgId string) string {
	return wcdb.wxfileindex.GetVideoPath(msgId)
}

func (wcdb WCDB) GetVoicePath(msgId string) string {
	return wcdb.wxfileindex.GetVoicePath(msgId)
}

// func (wcdb WCDB) GetFilePath(msgId string) string {
// 	return wcdb.wxfileindex.GetFilePath(msgId)
// }

func (wcdb WCDB) getMediaPath(chat ChatDetailListRow) ChatDetailListRow {
	switch chat.Type {
	case 3:
		// 图片
		chat.MediaPath = wcdb.enmicromsg.formatImagePath(chat.ImgPath)
		chat.MediaBCKPath = wcdb.enmicromsg.formatImageBCKPath(chat)
		chat.MediaSourcePath = wcdb.wxfileindex.GetImgPath(chat.MsgId)
		break
	case 34:
		// 语音
		chat.MediaPath = wcdb.enmicromsg.formatVoicePath(chat.ImgPath)
		break
	case 43:
		// 视频
		chat.MediaPath = wcdb.enmicromsg.formatVideoPath(chat.ImgPath)
		break
	case 1090519089:
		fileInfo := FileInfo{}
		filepath, fileSize := wcdb.wxfileindex.GetFilePath(chat.MsgId)
		fileInfo.FilePath = filepath
		fileInfo.FileSize = formatFileSize(fileSize)
		p := strings.Split(filepath, "/")
		if len(p) > 1 {
			fileName := p[len(p)-1]
			fileInfo.FileName = fileName
			fext := strings.Split(fileName, ".")
			if len(fext) > 1 {
				fileInfo.FileExt = fext[len(fext)-1]
			}
		}
		chat.FileInfo = fileInfo
	default:
		break
	}
	return chat
}

func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
