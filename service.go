package main

import (
	"bufio"
	pb "dirtyfilter/proto"
	"os"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"unicode/utf8"
)

type server struct {
	replaceWord string
	dirty_words map[string]bool
	//segmenter   sego.Segmenter
}

func (s *server) init() {
	s.replaceWord = "*"
	s.dirty_words = make(map[string]bool)

	//dictionary := "./dictionary.txt"
	dirty := "./dirty.txt"

	//// 载入字典
	//log.Info("Loading Dictionary...")
	//s.segmenter.LoadDictionary(dictionary)
	//log.Info("Dictionary Loaded")

	// 读取脏词库
	log.Info("Loading Dirty Words...")
	f, err := os.Open(dirty)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// 逐行扫描
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		words := strings.ToUpper(strings.TrimSpace(scanner.Text())) // 均处理为大写
		if words != "" {
			s.dirty_words[words] = true
		}
	}
	log.Info("Dirty Words Loaded")
}

func (s *server) Filter(ctx context.Context, in *pb.WordFilter_Text) (*pb.WordFilter_Text, error) {
	//bin := []byte(in.Text)
	//segments := s.segmenter.Segment(bin)
	//clean_text := make([]byte, 0, len(bin))
	//for _, seg := range segments {
	//	word := bin[seg.Start():seg.End()]
	//	if s.dirty_words[strings.ToUpper(string(word))] {
	//		clean_text = append(clean_text, []byte(strings.Repeat(s.replaceWord, utf8.RuneCount(word)))...)
	//	} else {
	//		clean_text = append(clean_text, word...)
	//	}
	//}
	//
	//
	//// 分词后没找到脏词，整句话再找一次
	//if string(clean_text) == in.Text {
	//	log.Info(in.Text)
	//	if s.dirty_words[strings.ToUpper(in.Text)] {
	//		log.Info(in.Text)
	//		clean_text = []byte(strings.Repeat(s.replaceWord, utf8.RuneCount([]byte(in.Text))))
	//	}
	//}
	//return &pb.WordFilter_Text{string(clean_text)}, nil

	words := in.Text
	indexMap := make(map[int]struct{})
	for key := range s.dirty_words {
		for _, index := range checkWord(strings.ToUpper(words), key) {
			indexMap[index] = struct{}{}
		}
	}
	wordsRune := []rune(words)
	for i := range indexMap {
		wordsRune[i] = '*'
		//b = append(b, wordsBytes[i])
		//if utf8.Valid(b) {
		//	clear += string(wordsBytes[i:]) + "*"
		//	b = []byte{}
		//}
	}
	return &pb.WordFilter_Text{Text: string(wordsRune)}, nil
}

func checkWord(words, dirtyKey string) (indexList []int) {
	n := 0
	for i := 0; i+len(dirtyKey) <= len(words); i++ {
		if words[i] == dirtyKey[0] && (len(dirtyKey) == 1 || words[i:i+len(dirtyKey)] == dirtyKey) {
			wordsByte := []byte(words)
			if i > 0 {
				n = utf8.RuneCount(wordsByte[0:i])
			}
			for j := 0; j < utf8.RuneCount(wordsByte[i:i+len(dirtyKey)]); j++ {
				indexList = append(indexList, n+j)
			}
			i += len(dirtyKey) - 1
		}
	}
	return
}