// main.go
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Movie struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	Poster string `json:"Poster"`
	Error  string `json:"Error"`
}

func main() {
	apiKey := flag.String("apikey", os.Getenv("OMDB_API_KEY"), "chave da API do OMDb (ou definir OMDB_API_KEY no ambiente)")
	outDir := flag.String("out", ".", "diretório onde salvar o pôster")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Uso: poster [--apikey CHAVE] [--out DIR] \"nome do filme\"")
		os.Exit(1)
	}
	movieName := strings.Join(flag.Args(), " ")

	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "Erro: é necessário fornecer uma API key (--apikey ou variável OMDB_API_KEY)")
		os.Exit(1)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Erro criando diretório de saída: %v\n", err)
		os.Exit(1)
	}

	movie, err := fetchMovie(movieName, *apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao buscar filme: %v\n", err)
		os.Exit(1)
	}

	if movie.Poster == "" || strings.ToLower(movie.Poster) == "n/a" {
		fmt.Printf("Filme encontrado: %s (%s)\nMas não há pôster disponível.\n", movie.Title, movie.Year)
		return
	}

	outPath, err := downloadPoster(movie.Poster, movie.Title, *outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao baixar pôster: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Filme: %s (%s)\nPôster salvo em: %s\n", movie.Title, movie.Year, outPath)
}

func fetchMovie(title, apiKey string) (*Movie, error) {
	url := fmt.Sprintf("https://www.omdbapi.com/?t=%s&apikey=%s", urlEncode(title), apiKey)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}

	if movie.Error != "" {
		return nil, errors.New(movie.Error)
	}

	return &movie, nil
}

func downloadPoster(url, title, outDir string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("erro HTTP %d ao baixar pôster", resp.StatusCode)
	}

	ext := filepath.Ext(url)
	if ext == "" {
		ext = ".jpg"
	}
	safeTitle := sanitizeFilename(title)
	outPath := filepath.Join(outDir, safeTitle+ext)

	outFile, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", err
	}

	return outPath, nil
}

func sanitizeFilename(name string) string {
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, ch := range invalid {
		name = strings.ReplaceAll(name, ch, "_")
	}
	return strings.TrimSpace(name)
}

func urlEncode(s string) string {
	s = strings.ReplaceAll(s, " ", "+")
	return s
}
