package main

import (
	"bufio"
	"os"
	"strings"
	pb "dirtyfilter/proto"

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
	dirty := "./dirty2017.txt"

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
		words := strings.Split(strings.ToUpper(strings.TrimSpace(scanner.Text())), " %") // 均处理为大写
		if words[0] != "" {
			s.dirty_words[words[0]] = true
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

	word := in.Text
	for key, _ := range s.dirty_words {
		i := strings.Index(strings.ToUpper(word), key)
		if i > -1 {
			word = checkWord(word, key, i)
		}
	}
	return &pb.WordFilter_Text{word}, nil
}

func checkWord(word, dirtyKey string, i int) (text string){
	bc := []byte(word)
	bd := []byte(dirtyKey)
	n := utf8.RuneCount(bd)
	rd := strings.Repeat("*",n)
	if strings.ToUpper(word) == dirtyKey {
		text = rd
	} else {
		text = string(bc[:i]) + rd + string(bc[i+len(dirtyKey):])
	}
	return
}