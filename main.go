package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/unicode"
)

var (
	bof      = flag.String("bof", "", "BOF Input File")
	outfile  = flag.String("out", "", "Output file to write to")
	ArgsFile = flag.String("args", "", "Args Input File")

	MagicHdr  = []byte{0x1c, 0x3f, 0xe6, 0x90}
	MagicaHdr = []byte{0x1c, 0x3f, 0xe6, 0x80}
)

func init() {
	flag.Parse()
}

func main() {
	if *bof == "" {
		flag.Usage()
		return
	}
	if *outfile == "" {
		*outfile = *bof + ".bin"
	}

	payload := BofLdr
	bof_payload, e := os.ReadFile(*bof)
	if e != nil {
		panic(e)
	}
	payload = append(payload, MagicHdr...)

	fmt.Println(len(bof_payload))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(len(bof_payload)))

	payload = append(payload, bs...)
	payload = append(payload, bof_payload...)

	payload = append(payload, MagicaHdr...)

	var ArgsBuff []byte

	if *ArgsFile == "" {
		binary.LittleEndian.PutUint32(bs, uint32(0))
	} else {
		ArgsBuff = bofConf(*ArgsFile)
		binary.LittleEndian.PutUint32(bs, uint32(len(ArgsBuff)))
	}

	payload = append(payload, bs...)

	payload = append(payload, ArgsBuff...)

	fmt.Printf("Writing %s\n", *outfile)
	os.WriteFile(*outfile, payload, 0644)

}

type BofArgs struct {
	ArgType string      `json:"type"`
	Value   interface{} `json:"value"`
}

type BOFArgsBuffer struct {
	Buffer *bytes.Buffer
}

func bofConf(conf string) []byte {

	//load conf
	config, _ := ioutil.ReadFile(conf)

	//config = []byte(CreatBofArgs("wstring", "C:\\"))
	var bargs []BofArgs
	argStr := strings.ReplaceAll(string(config), `\`, `\\`)
	fmt.Println(argStr)
	err := json.Unmarshal([]byte(argStr), &bargs)
	if err != nil {
		panic(err)
	}

	for i, a := range bargs {
		switch a.ArgType {
		case "binary":
			f := fmt.Sprintf("%v", a.Value)
			bargs[i].Value, err = ioutil.ReadFile(f)
			if err != nil {
				panic(err)
			}
		}
	}

	bofA := BOFArgsBuffer{
		Buffer: new(bytes.Buffer),
	}

	for _, a := range bargs {
		switch a.ArgType {
		case "integer":
			fallthrough
		case "int":
			if v, ok := a.Value.(float64); ok {
				err = bofA.AddInt(uint32(v))
			}
		case "string":
			if v, ok := a.Value.(string); ok {
				err = bofA.AddString(v)
			}
		case "wstring":
			if v, ok := a.Value.(string); ok {
				err = bofA.AddWString(v)
			}
		case "short":
			if v, ok := a.Value.(float64); ok {
				err = bofA.AddShort(uint16(v))
			}
		case "binary":
			if v, ok := a.Value.([]byte); ok {
				err = bofA.AddData([]byte(v))
			}
		}
		if err != nil {
			panic(err)
		}
	}

	parsedArgs, err := bofA.GetBuffer()
	if err != nil {
		panic(err)
	}
	return parsedArgs
}

func (b *BOFArgsBuffer) AddInt(d uint32) error {
	return binary.Write(b.Buffer, binary.LittleEndian, &d)
}
func (b *BOFArgsBuffer) AddData(d []byte) error {
	dataLen := uint32(len(d))
	err := binary.Write(b.Buffer, binary.LittleEndian, &dataLen)
	if err != nil {
		return err
	}
	return binary.Write(b.Buffer, binary.LittleEndian, &d)
}

func (b *BOFArgsBuffer) AddShort(d uint16) error {
	return binary.Write(b.Buffer, binary.LittleEndian, &d)
}

func (b *BOFArgsBuffer) AddString(d string) error {
	stringLen := uint32(len(d)) + 1
	err := binary.Write(b.Buffer, binary.LittleEndian, &stringLen)
	if err != nil {
		return err
	}
	dBytes := append([]byte(d), 0x00)
	return binary.Write(b.Buffer, binary.LittleEndian, dBytes)
}

func (b *BOFArgsBuffer) AddWString(d string) error {
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	strBytes := append([]byte(d), 0x00)
	utf16Data, err := encoder.Bytes(strBytes)
	if err != nil {
		return err
	}
	stringLen := uint32(len(utf16Data))
	err = binary.Write(b.Buffer, binary.LittleEndian, &stringLen)
	if err != nil {
		return err
	}
	return binary.Write(b.Buffer, binary.LittleEndian, utf16Data)
}

func (b *BOFArgsBuffer) GetBuffer() ([]byte, error) {
	outBuffer := new(bytes.Buffer)
	err := binary.Write(outBuffer, binary.LittleEndian, uint32(b.Buffer.Len()))
	if err != nil {
		return nil, err
	}
	err = binary.Write(outBuffer, binary.LittleEndian, b.Buffer.Bytes())
	if err != nil {
		return nil, err
	}
	return outBuffer.Bytes(), nil
}
