# 🚀 Exercícios Go - Capítulos 1-3 (Contexto Real de APIs)

## 📋 Como usar esta lista:
- Faça os exercícios em ordem (do mais fácil ao mais complexo)
- Crie um arquivo `.go` para cada exercício
- Commit no GitHub após cada um
- Se travar mais de 30min, pesquisa e segue
- **Objetivo: praticar Go, não sofrer!**

---

## 🟢 NÍVEL 1: Fundamentos (Warming Up)

### Ex 1.1: Validador de Status HTTP
**Conceitos:** variáveis, condicionais, inteiros
```
Crie uma função que recebe um código HTTP (int) e retorna:
- "Success" para códigos 200-299
- "Redirect" para códigos 300-399
- "Client Error" para códigos 400-499
- "Server Error" para códigos 500-599
- "Invalid" para outros

Teste com: 200, 404, 500, 301, 999
```

### Ex 1.2: Calculadora de Rate Limit
**Conceitos:** variáveis, operações matemáticas, float
```
Você tem:
- Total de requests permitidas por hora: 1000
- Requests já feitas: variável

Calcule:
1. Quantas requests ainda pode fazer
2. Qual a porcentagem usada (float)
3. Se passou do limite, quantas requests a mais foram feitas

Teste com: 500, 1000, 1200 requests
```

### Ex 1.3: Formatador de Log Level
**Conceitos:** strings, condicionais, iota (constantes)
```
Crie constantes usando iota para níveis de log:
DEBUG = 0, INFO = 1, WARNING = 2, ERROR = 3, CRITICAL = 4

Função que recebe um level (int) e retorna a string correspondente.

Bônus: Crie outra função que verifica se deve logar 
(ex: só loga se level >= WARNING)
```

---

## 🟡 NÍVEL 2: Strings e Validações (Mais Realista)

### Ex 2.1: Validador de Email Simples
**Conceitos:** strings, condicionais, funções de string
```
Crie uma função que valida se um email é válido:
- Deve ter exatamente um @
- Deve ter pelo menos um . depois do @
- Não pode começar ou terminar com espaços
- Não pode estar vazio

Retorne: bool (válido ou não) + string (mensagem de erro se inválido)

Teste com:
"user@example.com" ✅
"invalid.email" ❌
" user@example.com" ❌
"user@domain" ❌
```

### Ex 2.2: Parser de Query String
**Conceitos:** strings, loops, string literals
```
Receba uma query string tipo: "name=John&age=30&city=NYC"

Extraia e imprima:
- Cada chave
- Cada valor

Exemplo de saída:
name: John
age: 30
city: NYC

Dica: Use strings.Split()
```

### Ex 2.3: Sanitizador de Input de API
**Conceitos:** strings, bytes, loops
```
Crie uma função que "limpa" um input de usuário:
1. Remove espaços do início e fim
2. Converte para lowercase
3. Remove caracteres especiais (!@#$%&*)
4. Limita a 50 caracteres

Retorne a string limpa.

Teste com: "  Hello World! @#$  ", "TESTE123!@#"
```

---

## 🟠 NÍVEL 3: Processamento de Dados (Como em APIs)

### Ex 3.1: Contador de Requests por Método HTTP
**Conceitos:** arrays/slices, loops, condicionais
```
Você tem um slice de métodos HTTP que chegaram:
methods := []string{"GET", "POST", "GET", "DELETE", "GET", "POST", "PUT"}

Conte quantos de cada método chegaram e imprima:
GET: 3
POST: 2
DELETE: 1
PUT: 1

Dica: Use um mapa (map) ou contadores separados
```

### Ex 3.2: Calculadora de Tempo de Resposta Médio
**Conceitos:** slices, loops, float, operações matemáticas
```
Você tem tempos de resposta de API em ms:
responseTimes := []float64{120.5, 340.2, 89.7, 450.1, 200.3}

Calcule e imprima:
1. Tempo médio
2. Menor tempo
3. Maior tempo
4. Quantas requests foram mais rápidas que 200ms
```

### Ex 3.3: Validador de Payload JSON (Simplificado)
**Conceitos:** strings, bool, condicionais, bytes
```
Receba uma string que deveria ser JSON válido.

Verifique se:
1. Começa com { e termina com }
2. Contém pelo menos um :
3. Contém pelo menos um ,

Retorne: bool (parece JSON válido ou não)

Teste com:
`{"name":"John","age":30}` ✅
`{name:John}` ❌
`{"name":"John"` ❌

Obs: Não precisa validar JSON de verdade, só verificação básica!
```

---

## 🔴 NÍVEL 4: Cenários Complexos (Desafios)

### Ex 4.1: Gerador de Token Simples
**Conceitos:** strings, loops, constantes, random (novo!)
```
Gere um token aleatório para autenticação:
- 16 caracteres
- Apenas letras (A-Z, a-z) e números (0-9)
- Retorne como string

Dica: Você vai precisar do pacote "math/rand" (pode pesquisar!)
```

### Ex 4.2: Rate Limiter Simples
**Conceitos:** loops, condicionais, inteiros, bool
```
Simule um rate limiter:
- Limite: 5 requests por "janela"
- Receba um slice de requests: []int{1,2,3,4,5,6,7,8}
  (cada número representa uma request)

Para cada request:
- Se ainda não passou o limite: aceita
- Se passou: rejeita

Imprima:
Request 1: Accepted
Request 2: Accepted
...
Request 6: Rejected (limit exceeded)
...
```

### Ex 4.3: Construtor de URL de API
**Conceitos:** strings, structs (se já viu!), condicionais
```
Crie uma função que monta URLs de API:

Entrada:
- Base: "https://api.example.com"
- Endpoint: "/users"
- Params: map ou variáveis (id=123, active=true)

Saída:
"https://api.example.com/users?id=123&active=true"

Teste com diferentes combinações!
```

---

## 🎯 PROJETO FINAL DE SEMANA (Opcional, mas incrível!)

### Mini API Logger
**Conceitos:** TUDO que você viu até agora

Crie um programa que simula logs de uma API:

1. Receba um slice de "requests" (pode ser string com método + endpoint)
2. Para cada request:
   - Gere um timestamp fake (pode ser só contador)
   - Valide o método HTTP (GET, POST, etc)
   - Gere um status code "fake" (200, 404, 500)
   - Calcule um tempo de resposta fake (random entre 50-500ms)
   - Formate e imprima um log bonito

Exemplo de saída:
```
[1] GET /users - 200 OK (120ms)
[2] POST /users - 201 Created (340ms)
[3] GET /invalid - 404 Not Found (89ms)
[4] DELETE /users/1 - 500 Error (450ms)

Summary:
Total requests: 4
Success: 2
Errors: 2
Avg response time: 249.75ms
```

---

## 📝 Dicas Importantes:

1. **Não precisa fazer tudo de uma vez!** Escolha 3-4 que te interessam mais
2. **Pesquise sem culpa** - "how to do X in Go" é seu amigo
3. **Commit cada exercício** - mesmo que simples
4. **Se travar, pula** - você pode voltar depois
5. **Foque em entender, não em perfeição**

---

## 🎓 O que você vai aprender com esses exercícios:

✅ Manipulação de strings (comum em APIs)  
✅ Validações de input (essencial em qualquer API)  
✅ Processamento de dados (slices, loops)  
✅ Lógica de negócio (rate limiting, logs)  
✅ Construção de utilidades reais  

**Esses exercícios são o tipo de coisa que você usa TODO DIA em APIs!**

---

## 📊 Meta pro Final de Semana:

- [ ] Completar pelo menos 5 exercícios (qualquer nível)
- [ ] Commits no GitHub
- [ ] Se divertir! (sério, isso importa)

**Qualquer dúvida, me chama! Boa diversão! 🚀**