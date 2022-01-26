package osexec

import (
	"os/exec"
	"encoding/xml"
	"io"
	"bytes"
	"strings"
	"io/ioutil"
)

type Payload struct {
	Message string `xml:"message"`
}


func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}


func getXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := ioutil.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

