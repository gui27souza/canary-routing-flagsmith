## 22/8

- Comecei o projeto entendendo um pouco melhor sobre o gin, um framework para APIs REST em Go
    - Segui o tutorial simplezinho do [web-service-gin](https://go.dev/doc/tutorial/web-service-gin)
- Iniciei a implementação real agora, comecei no simples ainda, nas rotas de `health check`
    - Os endpoints `/healthz` e `/readyz`
    - Gerenciador de estado com RWMutex, com flags booleanas sobre o estado da aplicação e Snapshot de estado, para prover ao endpoint de `/readyz`
    - Inicialização do cliente do `Flagsmith`
- Pensando em Dev-Ex, também criei um `Makefile` e um `lefthook.yml` que formata o código ao commitar

```
.
│   .gitignore
│   Dockerfile
│   go.mod
│   go.sum
│   lefthook.yml
│   Makefile
│   PROJECT.md
│   README.md
├───cmd
│   └───api
│           main.go
└───internal
    ├───handlers
    │       apphandler.go
    └───state
            state.go
```

## 23/8

- De início, tive uma conversa com o Gemini para entender os benefícios de configurações específicas na inicialização do sdk do `flagsmith`
    - A ideia é que ele baixe todas as flags localmente ao inicializar e tenha um melhor gerenciamento da leitura dessas flags, em memória
- Além disso, isolei de fato o inicializador do sdk por trás de uma struct/interface, em `internal/services/flags.go` , facilitando a manutenção e testabilidade
- Descobri 3 ferramentas para tornar a saída de testes mais agradável:
```
// gotestsum
go install gotest.tools/gotestsum@latest
gotestsum --format pkgname

// tparse
go install github.com/mfridman/tparse@latest
go test -json ./... | tparse -all

// coverage nativo
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 24/8

- Fiz uns refinamentos no DevEx, cons uns ajustezinhos no `Makefile`, e com checagem de env e alguns encadeamentos de comando
- Pude finalmente testar a aplicação! E deu tudo certo!!
    - Funcionou o start do servidor, hidradação das flags e os endpoints de saúde `/healthz` e `/readyz`!!

- Para o encorpar o objetivo do projeto, tomei a decisão de implementar um roteamento dinâmico de requisições.
- Basicamente, o microserviço será responsável por dizer se o cliente X deve ser direcionado a versão v1 ou v2 de dada aplicação:
    - Com base no seu ID
    - Levando em conta as flags de roteamento:
        - Se a flag geral de v2 esta ligada
        - Flag de configuração sobre o roteamento: percentual de clientes a serem direcionados à v2, países permitidos, etc

## 25/8

- Hoje comecei a implementar a lógica real de roteamento, com base nas flags e configs
    - Comecei com os kill switches e denied routing
    - Criei também métodos auxiliares, para o snapshot do state e gerar o response de denied
    - Fiz os casos de denied:
        - Feature v2 desabilitada
        - Erro ao buscar regras de canary routing
        - Erro ao parsear regras de canary routing
        - Feature desabilitada para o país do usuário
    - O próximo passo é implementar o roteamento real com o hash em cima do id do cliente

- Fiz uma normalização otimizada da verificação do país do usuário no roteamento
- Fiz a implementação do roteamento real usando hash!!
    - criei um utils/hash e implementei um normalizador em hash de 0 a 100
    - Com essa normalização de porcentagem, consegui fazer o roteamento de Canary percentual com base no bucket do hash do user id!!
- Implementei a busca de config json no flagsmith

- Com isso, melhorei a documentação do engine.go e flagsmith.go

- Criei também testes automatizados do hash
    - Com ele, aprendi a filosofia de que um teste também é importante para alertar mudanças em comportamentos estáticos/determinísticos!

## 26/8 a 28/8

- Aprendi a como criar mocks e como me organizar melhor, junto da proximidade de quem os usa:
  - Se o uso for apenas do mesmo pacote, o ideal é criar um mocks.go e manter o princípio de localidade
  - Agora, se for algo usado por outras partes do código, o ideal é padronizar um `testutil/` com geradores de mocks

- Ajustei também o Makefile para refletir melhor a cobertura de testes, ignorando os arquivosde teste, `cmd/` e os mocks

- Com esses novos conhecimentos, criei o RouteHandler test, restando apenas o Router test

- Essa é a visão até agr:

```
.
│   .gitignore
│   AGENTS.md
│   Dockerfile
│   go.mod
│   go.sum
│   lefthook.yml
│   Makefile
│   NOTES.md
│   PROJECT.md
│   README.md
├───cmd
│   └───api
│           main.go
└───internal
    ├───handlers
    │       apphandler.go
    │       apphandler_test.go
    │       route_handler.go
    │       route_handler_test.go
    ├───service
    │   ├───flags
    │   │       flags.go
    │   │       flagsmith.go
    │   │       flagsmith_test.go
    │   └───router
    │           engine.go
    │           engine_test.go
    ├───state
    │       state.go
    │       state_test.go
    ├───testutil
    │       flags.go
    │       router.go
    └───util
        └───hash
                hash.go
                hash_test.go
```

## 29/8

- Acho que finalmente entendi na prática o que é Dependency Inversion:
  - Para tornar os testes de router.Engine mais "testáveis", precisei remover a dependência direta que a estrutura tinha com o cálculo de bucket, que estava preso à função de Hash Normalizado e também ao time.Now() do response
  - A ideia foi passar tipos de função meio genéricos, especificando apenas a interface do que precisava delas
  - Com isso, consegui moldar essas dependências ao meu dispor, facilitando os testes, sem precisar ficar fazendo gambiarras:
    - Principalmente com as comparações do time.Now, complicava bastante no assert
    - Ou em ter q deixar o bucket que dado hash para dado input cairia hardcoded
  - Isso abriu caminho para o próximo passo dos testes: Table-Driven
- Também tornei o Makefile agnóstico ao SO, funcionando para windows ou linux/mac

## 31/8 a 1/9

- Entendi um pouco melhor sobre tratativa de erros, principalmente quando retornar ou nao um error em uma função
- Avancei bastante nos testes:
    - Aprendi Table-Driven testing como árvore de decisão testável
    - O limite dos testes, subir cobertura só para aumentar o número é pior -> torna implementações que deveriam ser simples completas e não dão real segurança apenas por ter uma cobertura ampla
- Entendi muito melhor sobre dependências e a importância das interfaces, como desacoplar
    - Quebrei as structs atreladas à engine em structs de domínio, dessa forma o route_handler não precisa conhecer a engine nem vice-versa!
