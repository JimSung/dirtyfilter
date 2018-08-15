package main

import (
	"log"
	"sync"
	"testing"
	"time"

	"dirtyfilter/proto"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	address  = "localhost:50002"
	testText = []string{
		"我操你大爷，法轮大法好",
		"Fuck you，fuck you sisters!",
		"1100y.com",
		"快来1100y.com",
		"水乳交融",
		"abc",
		"毛泽东毛泽东",
		"操你妈",
		"Ecru Game Eagle",
	}
	conn *grpc.ClientConn
)

// 模拟client测试
func init() {
	// Set up a connection to the server.
	_conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}
	conn = _conn
	//	fmt.Println(replaceByte, replaceLenth)
}

func TestWordFilter(t *testing.T) {

	c := proto.NewWordFilterServiceClient(conn)

	// Contact the server and print out its response.
	for i := 0; i < len(testText); i++ {
		r, err := c.Filter(context.Background(), &proto.WordFilter_Text{Text: testText[i]})
		if err != nil {
			t.Fatalf("could not query: %v", err)
		}
		t.Logf("Filtered Text: %s", r.Text)
	}
}

func Test_WorldFiter(t *testing.T) {
	last := time.Now()
	c := proto.NewWordFilterServiceClient(conn)
	wait := sync.WaitGroup{}
	for i := 0; i < len(testText); i++ {
		wait.Add(1)
		go func() {
			r, err := c.Filter(context.Background(), &proto.WordFilter_Text{Text: testText[i]})
			if err != nil {
				t.Fatal(err)
			}
			_ = r
			wait.Done()
		}()
	}

	wait.Wait()
	t.Logf("cost:%v", time.Now().Sub(last))
}

func TestCheck(t *testing.T) {
	i := strings.Index("1100y.com", "1100y.com")
	fmt.Println(i)
}

func TestTrim(t *testing.T) {
	s := "     shhh aaa    "
	fmt.Println(strings.Trim(s, " "))
}

func TestLen(t *testing.T) {
	a := "fuck"
	aByte := []byte(a)
	fmt.Println(string(aByte[0:4]))
	fmt.Println(utf8.RuneCount(aByte[0:4]))
}

func TestType(t *testing.T) {
	a := []int{1, 4, 2, 3, 5}
	sort.Slice(a, func(i, j int) bool {
		x, y := a[i], a[j]
		if x > y {
			return true
		}
		return false
	})
	for _, b := range a {
		fmt.Println(b)
	}
}
