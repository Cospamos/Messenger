package app

import (
  "fmt"

  "github.com/skip2/go-qrcode"
)

type Tunnel struct {
  PublicURL string `json:"public_url"`
}
var NgrokTunnel Tunnel

type TunnelsResponse struct {
  Tunnels []Tunnel `json:"tunnels"`
}

var SaveDir = "www/asset/img/qrcode.png"

func GenQR() error {
  url := "https://zonia-interparliament-nonnormally.ngrok-free.dev"
  if url == "" {
    return fmt.Errorf("no public URL available. Is ngrok running?")
  }

  err := qrcode.WriteFile(url, qrcode.Medium, 256, SaveDir)
  if err != nil {
    return fmt.Errorf("failed to generate QR code: %v", err)
  }
  
  fmt.Println("QR code saved to", SaveDir)
  return nil
}