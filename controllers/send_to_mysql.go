package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/dao/mysql"
	"go_web/models"
	"io"
	"net/http"
	"os"
	"strings"
)

// WordBookToMySQL 导入单词书
func WordBookToMySQL(c *gin.Context) {
	// 打开json文件
	jsonFile, err := os.Open("pkg/word_books/KaoYan_2.json")

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	info := string(byteValue)
	result := strings.Split(info, "\n")
	for _, v := range result {
		var jsWordInfo models.WordInfo
		dict := JsonToMap(v)
		if dict == nil {
			continue
		}
		// 获取单词序号
		jsWordInfo.No = int(dict["wordRank"].(float64))
		zap.L().Info("序号")

		// 获取单词
		jsWordInfo.Word = dict["headWord"].(string)
		zap.L().Info("单词")

		content := dict["content"]
		word, _ := content.(map[string]interface{})["word"]

		contentNext, _ := word.(map[string]interface{})["content"]

		//获取音标
		phonesInfo, _ := contentNext.(map[string]interface{})
		_, ok := phonesInfo["phone"]
		if ok {
			phones, _ := contentNext.(map[string]interface{})["phone"]
			jsWordInfo.Phone = strings.Split(phones.(string), ",")[0]
			zap.L().Info("音标")
		}

		// 获取单词翻译
		synoInfo, _ := contentNext.(map[string]interface{})
		_, ok = synoInfo["syno"]
		if ok {
			syno, _ := contentNext.(map[string]interface{})["syno"]
			synos, _ := syno.(map[string]interface{})["synos"]
			pos := synos.([]interface{})
			for _, singlePos := range pos {
				wordPos := singlePos.(map[string]interface{})["pos"]
				tran := singlePos.(map[string]interface{})["tran"]
				jsWordInfo.Translation = jsWordInfo.Translation + wordPos.(string) + "." + tran.(string) + "\t"
			}
			jsWordInfo.Translation = fmt.Sprintf("%s\n", jsWordInfo.Translation)
			zap.L().Info("翻译")
		}

		// 获取示例句子
		sentenceInfo, _ := contentNext.(map[string]interface{})
		_, ok = sentenceInfo["sentence"]
		if ok {
			sentences, _ := contentNext.(map[string]interface{})["sentence"]
			sentence, _ := sentences.(map[string]interface{})["sentences"]
			sentenceSplit := sentence.([]interface{})
			for _, sen := range sentenceSplit {
				singleSentence := sen.(map[string]interface{})
				sContent := singleSentence["sContent"]
				sCn := singleSentence["sCn"]
				jsWordInfo.Sentence = jsWordInfo.Sentence + sContent.(string) + "\n" + sCn.(string) + "\n"
			}
			zap.L().Info("句子")
		}

		// 获取短语
		phraseInfo, _ := contentNext.(map[string]interface{})
		_, ok = phraseInfo["phrase"]
		if ok {
			phrase, _ := contentNext.(map[string]interface{})["phrase"]
			phrases, _ := phrase.(map[string]interface{})["phrases"]
			phraseSplit := phrases.([]interface{})
			for _, singelPhrase := range phraseSplit {
				pContent := singelPhrase.(map[string]interface{})["pContent"]
				pCn := singelPhrase.(map[string]interface{})["pCn"]
				jsWordInfo.Phrase = jsWordInfo.Phrase + pContent.(string) + "  " + pCn.(string) + "\t"
			}
			jsWordInfo.Phrase = jsWordInfo.Phrase + "\n"
			zap.L().Info("短语")
		}

		// 获取变形
		relWords, _ := contentNext.(map[string]interface{})
		_, ok = relWords["relWord"]
		if ok {
			relWord := relWords["relWord"]
			rels, _ := relWord.(map[string]interface{})["rels"]
			relSplit := rels.([]interface{})
			for _, rel := range relSplit {
				relPos := rel.(map[string]interface{})["pos"]
				words := rel.(map[string]interface{})["words"]
				wordSplit := words.([]interface{})
				for _, relPara := range wordSplit {
					hwd := relPara.(map[string]interface{})["hwd"]
					tran := relPara.(map[string]interface{})["tran"]
					if hwd == nil || relPos == nil || tran == nil {
						jsWordInfo.WordMorph = ""
					} else {
						jsWordInfo.WordMorph = jsWordInfo.WordMorph + hwd.(string) + "\t" + relPos.(string) + "." + tran.(string) + "\t"
					}
				}
				jsWordInfo.WordMorph = fmt.Sprintf("%s\n", jsWordInfo.WordMorph)
			}
			zap.L().Info("变形")
		}

		// 导入单词书到数据库
		err = mysql.WordBookToMySQL(jsWordInfo)
		if err != nil {
			zap.L().Error("导入单词书失败")
			return
		}
	}

	// 5.登录成功返回
	zap.L().Info("导入单词书成功")
	c.JSON(http.StatusOK, gin.H{
		"status": -1,
		"msg":    "导入单词书成功",
	})
}

func JsonToMap(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		return nil
	}
	return tempMap
}
