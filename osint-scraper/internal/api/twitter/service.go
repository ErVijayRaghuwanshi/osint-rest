package twitter

import (
    "net/http"
    "time"
)

type Service struct{}

func NewService() *Service {
    return &Service{}
}

func (s *Service) CheckWebsite() bool {
    client := http.Client{
        Timeout: 9 * time.Second,
    }

    resp, err := client.Get("https://www.twitter.com")
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode == 200
}
