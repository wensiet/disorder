package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Discord struct {
	timeout time.Time
	channel string
	token   string
}

func New(token, channel string) *Discord {
	return &Discord{
		timeout: time.Now(),
		channel: channel,
		token:   token,
	}
}

func (d *Discord) SendChunk(filename string, chunk []byte) (map[string]interface{}, error) {
	requestLink := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", d.channel)

	body := &bytes.Buffer{}
	reader := bytes.NewReader(chunk)
	writer := multipart.NewWriter(body)

	// Add the file to the form-data
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return nil, err
	}

	nameField, err := writer.CreateFormField("name")
	if err != nil {
		return nil, err
	}
	_, err = nameField.Write([]byte(filename))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", requestLink, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bot "+d.token)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
