package main

import (
	"bufio"
	pb "dirtyfilter/proto"
	"os"
	"strings"

	"golang.org/x/net/context"

	"unicode/utf8"
	"log"
)

type Config struct {
	Host        string `envconfig:"default=loaclhost"`
	Port        int    `envconfig:"default=50002"`
	Path        string `envconfig:"default=./dirty.txt"`
	ReplaceWord string `envconfig:"default=*"`
}

type server struct {
	replaceWord string
	dirtyWords  map[string]bool
}

func (s *server) init() {
	s.replaceWord = "*"
	s.dirtyWords = make(map[string]bool)

	dirty := "./dirty.txt"

	// 读取脏词库
	log.Printf("Loading Dirty Words...")
	f, err := os.Open(dirty)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 逐行扫描
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		words := strings.ToUpper(strings.TrimSpace(scanner.Text())) // 均处理为大写
		if words != "" {
			s.dirtyWords[words] = true
		}
	}
	log.Printf("Dirty Words Loaded")
}

func (s *server) Filter(ctx context.Context, in *pb.WordFilter_Text) (*pb.WordFilter_Text, error) {
	words := in.Text
	indexMap := make(map[int]struct{})
	for key := range s.dirtyWords {
		for _, index := range checkWord(strings.ToUpper(words), key) {
			indexMap[index] = struct{}{}
		}
	}
	wordsRune := []rune(words)
	for i := range indexMap {
		wordsRune[i] = '*'
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
