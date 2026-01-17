package app

import (
    "fmt"
    vt "github.com/VirusTotal/vt-go"
)

func CheckImgByUrl(urlToCheck string) error {
    client := vt.NewClient(VirusTotalAPI)

    scanner := client.NewURLScanner()

    obj, err := scanner.Scan(urlToCheck)
    if err != nil {
        return fmt.Errorf("scan error: %w", err)
    }

    analysisID := obj.ID()
    fmt.Println("Analysis ID:", analysisID)

    urlReport := vt.URL("analyses/%s", analysisID)

    resp, err := client.Get(urlReport)
    if err != nil {
        return fmt.Errorf("get report error: %w", err)
    }

    fmt.Printf("Raw report JSON: %s\n", resp.Data)
    return nil
}
