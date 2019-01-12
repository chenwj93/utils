package utils

import (
	"regexp"
	"bytes"
	"strings"
	"github.com/satori/go.uuid"
	"strconv"
	"errors"
	"math"
)

type RandomCode struct {
	chars        string
	charsLen     int64
	length       int
	pos          int
	rearrangeFun []func(code *string)
}

func NewRandom(chars string, length int) (*RandomCode, error) {
	exp := "^[0-9a-zA-Z]+$"
	reg, _ := regexp.Compile(exp)
	if !reg.MatchString(chars) {
		return nil, errors.New("source code only support number and English letter")
	}
	return &RandomCode{chars: chars, charsLen: int64(len(chars)), length: length}, nil
}

// mount custom rearrange function
func (r *RandomCode) MountRearrangeFunc(fun ...func(code *string)) {
	r.rearrangeFun = append(r.rearrangeFun, fun...)
}

// 仅支持 数字 + 英文字母
func (r *RandomCode) GenerateCode() string {
	i := float64(0)
	code := &bytes.Buffer{}
	length := r.length
	for length != 0 {
		if length&1 != 0 {
			l := math.Pow(2, i)
			code.Write(r.generate(int(l)))
		}
		length = length >> 1
		i++
	}
	return r.rearrange(code)
}

func (r *RandomCode) generate(length int) []byte {
	u := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	var code bytes.Buffer
	off := 32 / length
	for i := 0; i < length; i++ {
		str := u[(i * off):(i*off + off)]
		i1, _ := strconv.ParseInt(str, 16, 64)
		code.WriteByte(r.chars[i1%r.charsLen])
	}
	return code.Bytes()
}

func (r *RandomCode) rearrange(buf *bytes.Buffer) string {
	code := r.disrupt(buf)
	for _, fun :=range r.rearrangeFun {
		fun(&code)
	}
	return code
}

func (r *RandomCode) disrupt(buf *bytes.Buffer) string {
	pos := r.pos
	r.pos = (r.pos + 1) % r.length
	code := buf.Bytes()
	next := (pos + r.length/2) % r.length
	temp := code[pos]
	code[pos] = code[next]
	code[next] = temp
	return string(code)
}
