package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

// Issue representa um issue do GitHub
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
	Milestone *Milestone
}

// User representa um usuário do GitHub
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// Milestone representa um milestone do GitHub
type Milestone struct {
	Title string
	URL   string `json:"html_url"`
}

// IssuesSearchResult representa o resultado da busca
type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

var (
	cachedIssues *IssuesSearchResult
	issueList    = template.Must(template.New("issuelist").Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>GitHub Issues</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; }
        h1 { color: #333; }
        .nav { margin: 20px 0; padding: 10px; background: #e9ecef; border-radius: 4px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #007bff; font-weight: bold; }
        .nav a:hover { text-decoration: underline; }
        .issue { border: 1px solid #ddd; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .issue h3 { margin-top: 0; }
        .issue-meta { color: #666; font-size: 0.9em; }
        .state-open { color: #28a745; font-weight: bold; }
        .state-closed { color: #dc3545; font-weight: bold; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background: #f8f9fa; font-weight: bold; }
        tr:hover { background: #f8f9fa; }
        .count { color: #666; font-size: 0.9em; margin: 10px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>GitHub Issues - {{.Repo}}</h1>
        <div class="nav">
            <a href="/">Todas as Issues</a>
            <a href="/milestones">Milestones</a>
            <a href="/users">Usuários</a>
        </div>
        <p class="count">Total de issues: {{.TotalCount}}</p>
        {{if eq .View "issues"}}
            {{range .Items}}
            <div class="issue">
                <h3><a href="{{.HTMLURL}}" target="_blank">#{{.Number}} {{.Title}}</a></h3>
                <div class="issue-meta">
                    <span class="state-{{.State}}">{{.State}}</span> | 
                    Criado por <a href="{{.User.HTMLURL}}" target="_blank">{{.User.Login}}</a> | 
                    {{.CreatedAt.Format "02/01/2006"}}
                    {{if .Milestone}} | Milestone: {{.Milestone.Title}}{{end}}
                </div>
                {{if .Body}}
                <p>{{printf "%.200s" .Body}}{{if gt (len .Body) 200}}...{{end}}</p>
                {{end}}
            </div>
            {{end}}
        {{else if eq .View "milestones"}}
            <table>
                <tr>
                    <th>Milestone</th>
                    <th>Issues</th>
                </tr>
                {{range $milestone, $issues := .Milestones}}
                <tr>
                    <td><strong>{{if $milestone}}{{$milestone}}{{else}}Sem Milestone{{end}}</strong></td>
                    <td>{{len $issues}}</td>
                </tr>
                {{end}}
            </table>
        {{else if eq .View "users"}}
            <table>
                <tr>
                    <th>Usuário</th>
                    <th>Issues Criadas</th>
                    <th>Perfil</th>
                </tr>
                {{range $user, $count := .Users}}
                <tr>
                    <td><strong>{{$user.Login}}</strong></td>
                    <td>{{$count}}</td>
                    <td><a href="{{$user.HTMLURL}}" target="_blank">Ver Perfil</a></td>
                </tr>
                {{end}}
            </table>
        {{end}}
    </div>
</body>
</html>
`))
)

// SearchIssues consulta a API do GitHub
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha na busca: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Handler para a página principal (todas as issues)
func handleIssues(w http.ResponseWriter, r *http.Request) {
	if cachedIssues == nil {
		http.Error(w, "Nenhuma consulta realizada ainda", http.StatusInternalServerError)
		return
	}

	data := struct {
		Repo       string
		TotalCount int
		Items      []*Issue
		View       string
	}{
		Repo:       "golang/go",
		TotalCount: cachedIssues.TotalCount,
		Items:      cachedIssues.Items,
		View:       "issues",
	}

	if err := issueList.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

// Handler para milestones
func handleMilestones(w http.ResponseWriter, r *http.Request) {
	if cachedIssues == nil {
		http.Error(w, "Nenhuma consulta realizada ainda", http.StatusInternalServerError)
		return
	}

	milestones := make(map[string][]*Issue)
	for _, issue := range cachedIssues.Items {
		key := "Sem Milestone"
		if issue.Milestone != nil {
			key = issue.Milestone.Title
		}
		milestones[key] = append(milestones[key], issue)
	}

	data := struct {
		Repo       string
		TotalCount int
		Milestones map[string][]*Issue
		View       string
	}{
		Repo:       "golang/go",
		TotalCount: cachedIssues.TotalCount,
		Milestones: milestones,
		View:       "milestones",
	}

	if err := issueList.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

// Handler para usuários
func handleUsers(w http.ResponseWriter, r *http.Request) {
	if cachedIssues == nil {
		http.Error(w, "Nenhuma consulta realizada ainda", http.StatusInternalServerError)
		return
	}

	users := make(map[*User]int)
	for _, issue := range cachedIssues.Items {
		if issue.User != nil {
			users[issue.User]++
		}
	}

	data := struct {
		Repo       string
		TotalCount int
		Users      map[*User]int
		View       string
	}{
		Repo:       "golang/go",
		TotalCount: cachedIssues.TotalCount,
		Users:      users,
		View:       "users",
	}

	if err := issueList.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Faz a consulta inicial ao GitHub (uma única vez)
	fmt.Println("Consultando issues do repositório golang/go...")
	result, err := SearchIssues([]string{"repo:golang/go", "is:open"})
	if err != nil {
		log.Fatal(err)
	}
	cachedIssues = result
	fmt.Printf("Encontradas %d issues\n", result.TotalCount)

	// Configura os handlers
	http.HandleFunc("/", handleIssues)
	http.HandleFunc("/milestones", handleMilestones)
	http.HandleFunc("/users", handleUsers)

	// Inicia o servidor
	fmt.Println("Servidor rodando em http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}