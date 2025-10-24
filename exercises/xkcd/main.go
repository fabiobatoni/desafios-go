// main.go
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// Estrutura simplificada do JSON do xkcd (campos que vamos usar)
type XKCD struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
}

const (
	defaultCacheDirName = ".xkcd-cache"
	indexFilename       = "index.json"
)

// Tipo do índice: token -> []int (lista de números dos quadrinhos)
type Index map[string][]int

func main() {
	// subcomandos: index e search
	if len(os.Args) < 2 {
		usageAndExit()
	}

	cmd := os.Args[1]
	switch cmd {
	case "index":
		indexCmd(os.Args[2:])
	case "search":
		searchCmd(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n\n", cmd)
		usageAndExit()
	}
}

func usageAndExit() {
	fmt.Println(`Uso:
  xkcd index [--cache DIR] [--workers N] [--rebuild]
    Baixa (uma vez) todos os JSON do xkcd e cria/atualiza o índice invertido.

  xkcd search [--cache DIR] TERM [TERM ...]
    Busca TERM(s) no índice e exibe URL + transcrição dos quadrinhos que casam.

Exemplos:
  xkcd index --cache ~/.xkcd-cache
  xkcd search --cache ~/.xkcd-cache "quantum" "cat"
`)
	os.Exit(1)
}

func indexCmd(args []string) {
	fs := flag.NewFlagSet("index", flag.ExitOnError)
	cacheDir := fs.String("cache", filepath.Join(os.Getenv("HOME"), defaultCacheDirName), "diretório de cache")
	workers := fs.Int("workers", runtime.NumCPU(), "número de workers para download")
	rebuild := fs.Bool("rebuild", false, "forçar rebuild do índice (re-indexa arquivos em cache)")
	fs.Parse(args)

	if err := os.MkdirAll(*cacheDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "erro criando cache dir: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Obtendo número do quadrinho mais recente...")
	latest, err := fetchLatestNum()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro obtendo latest: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("último quadrinho: %d\n", latest)

	// baixar todos JSONs com cache (pula se já existir)
	fmt.Println("Baixando JSONs (se ainda não existirem)...")
	if err := downloadAll(latest, *cacheDir, *workers); err != nil {
		fmt.Fprintf(os.Stderr, "erro download: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Download concluído (cache).")

	// construir índice a partir dos JSONs no cache
	idxPath := filepath.Join(*cacheDir, indexFilename)
	if *rebuild {
		fmt.Println("Rebuild forçado do índice.")
		if err := os.Remove(idxPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "erro removendo índice antigo: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Construindo índice...")
	index, err := buildIndexFromCache(*cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro construindo índice: %v\n", err)
		os.Exit(1)
	}
	if err := saveIndex(idxPath, index); err != nil {
		fmt.Fprintf(os.Stderr, "erro salvando índice: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Índice salvo em %s (tokens: %d)\n", idxPath, len(index))
}

func searchCmd(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	cacheDir := fs.String("cache", filepath.Join(os.Getenv("HOME"), defaultCacheDirName), "diretório de cache")
	fs.Parse(args)

	terms := fs.Args()
	if len(terms) == 0 {
		fmt.Fprintln(os.Stderr, "forneça pelo menos um termo de busca")
		usageAndExit()
	}

	idxPath := filepath.Join(*cacheDir, indexFilename)
	index, err := loadIndex(idxPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro carregando índice (%s): %v\n", idxPath, err)
		fmt.Fprintf(os.Stderr, "rodar `xkcd index --cache %s` primeiro\n", *cacheDir)
		os.Exit(1)
	}

	// tokenizar termos de busca, buscar interseção
	tokens := tokenize(strings.Join(terms, " "))
	if len(tokens) == 0 {
		fmt.Fprintln(os.Stderr, "nenhum token válido nos termos")
		os.Exit(1)
	}

	// obter lista de quadrinhos que satisfazem todos tokens (AND)
	var resultIDs []int
	for i, tok := range tokens {
		ids := index[tok]
		if i == 0 {
			resultIDs = append([]int{}, ids...)
		} else {
			resultIDs = intersect(resultIDs, ids)
		}
		if len(resultIDs) == 0 {
			break
		}
	}
	if len(resultIDs) == 0 {
		fmt.Println("Nenhum resultado encontrado.")
		return
	}

	// ordenar e imprimir cada quadrinho com URL + transcrição
	sort.Ints(resultIDs)
	for _, id := range resultIDs {
		path := filepath.Join(*cacheDir, fmt.Sprintf("%d.json", id))
		c, err := loadComicFromFile(path)
		if err != nil {
			// se não tiver no cache, apenas pular
			fmt.Fprintf(os.Stderr, "erro lendo %s: %v\n", path, err)
			continue
		}
		printComicResult(c)
	}
}

// fetchLatestNum pega o número do quadrinho mais recente via https://xkcd.com/info.0.json
func fetchLatestNum() (int, error) {
	const url = "https://xkcd.com/info.0.json"
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}
	var c XKCD
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return 0, err
	}
	if c.Num == 0 {
		return 0, errors.New("num zero no latest")
	}
	return c.Num, nil
}

// downloadAll baixa todos os JSONs de 1..latest se ainda não existirem no cache
func downloadAll(latest int, cacheDir string, workers int) error {
	type job struct{ n int }
	jobs := make(chan job, workers*2)
	wg := sync.WaitGroup{}
	errCh := make(chan error, workers)

	client := &http.Client{Timeout: 20 * time.Second}

	// workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				path := filepath.Join(cacheDir, fmt.Sprintf("%d.json", j.n))
				if _, err := os.Stat(path); err == nil {
					// já existe
					continue
				}
				if err := downloadComic(client, j.n, path); err != nil {
					errCh <- fmt.Errorf("erro baixando %d: %w", j.n, err)
					return
				}
				// pequeno sleep para não sobrecarregar
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	// enfileira
	go func() {
		for n := 1; n <= latest; n++ {
			jobs <- job{n: n}
		}
		close(jobs)
	}()

	wg.Wait()
	close(errCh)
	// se houver erro, retornar primeiro
	for e := range errCh {
		return e
	}
	return nil
}

func downloadComic(client *http.Client, n int, path string) error {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", n)
	// Retry simples
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			fmt.Printf("quadrinho %d não existe (404) — ignorando\n", n)
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
			resp.Body.Close()
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			continue
		}
		// salvar em arquivo temporário e renomear (segurança)
		tmp := path + ".tmp"
		f, err := os.Create(tmp)
		if err != nil {
			resp.Body.Close()
			return err
		}
		_, err = io.Copy(f, resp.Body)
		resp.Body.Close()
		f.Close()
		if err != nil {
			os.Remove(tmp)
			lastErr = err
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			continue
		}
		if err := os.Rename(tmp, path); err != nil {
			return err
		}
		return nil
	}
	return lastErr
}

// buildIndexFromCache varre os arquivos *.json no cache e constroi o índice invertido
func buildIndexFromCache(cacheDir string) (Index, error) {
	idx := make(Index)
	err := filepath.WalkDir(cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}
		// pular o próprio arquivo index.json
		if d.Name() == indexFilename {
			return nil
		}
		c, err := loadComicFromFile(path)
		if err != nil {
			// ignore arquivos inválidos
			fmt.Fprintf(os.Stderr, "warn: não foi possível ler %s: %v\n", path, err)
			return nil
		}
		// obter tokens do conjunto de campos que queremos pesquisar
		text := strings.Join([]string{c.Title, c.SafeTitle, c.Alt, c.Transcript}, " ")
		toks := tokenize(text)
		unique := map[string]struct{}{}
		for _, t := range toks {
			unique[t] = struct{}{}
		}
		for t := range unique {
			idx[t] = append(idx[t], c.Num)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// ordenar listas e remover duplicatas (segurança)
	for k := range idx {
		arr := uniqueInts(idx[k])
		sort.Ints(arr)
		idx[k] = arr
	}
	return idx, nil
}

func loadComicFromFile(path string) (*XKCD, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c XKCD
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func saveIndex(path string, idx Index) error {
	b, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func loadIndex(path string) (Index, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var idx Index
	if err := json.Unmarshal(b, &idx); err != nil {
		return nil, err
	}
	return idx, nil
}

// tokenize: separa palavras por qualquer rune que não seja letra ou dígito e normaliza pra minúsculas
func tokenize(s string) []string {
	f := func(r rune) bool {
		// considera letras e dígitos como parte do token
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return false
		}
		// tratar acentos: usar unicode.IsLetter seria mais completo, mas evita dependências aqui
		return true
	}
	raw := strings.FieldsFunc(strings.ToLower(s), f)
	// filtrar tokens vazios e curtos
	out := make([]string, 0, len(raw))
	for _, t := range raw {
		t = strings.TrimSpace(t)
		if t == "" || len(t) < 2 {
			continue
		}
		out = append(out, t)
	}
	return out
}

func uniqueInts(a []int) []int {
	m := map[int]struct{}{}
	for _, v := range a {
		m[v] = struct{}{}
	}
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func intersect(a, b []int) []int {
	set := map[int]struct{}{}
	for _, v := range a {
		set[v] = struct{}{}
	}
	var out []int
	for _, v := range b {
		if _, ok := set[v]; ok {
			out = append(out, v)
		}
	}
	return uniqueInts(out)
}

func printComicResult(c *XKCD) {
	url := fmt.Sprintf("https://xkcd.com/%d/", c.Num)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("Num: %d\nTitle: %s\nURL: %s\nImage: %s\n", c.Num, c.Title, url, c.Img)
	if strings.TrimSpace(c.Transcript) != "" {
		fmt.Println("\n--- Transcript ---")
		fmt.Println(strings.TrimSpace(c.Transcript))
	} else if strings.TrimSpace(c.Alt) != "" {
		fmt.Println("\n--- Alt / Hover text ---")
		fmt.Println(strings.TrimSpace(c.Alt))
	} else {
		fmt.Println("\n(sem transcript nem alt)")
	}
	fmt.Println()
}
