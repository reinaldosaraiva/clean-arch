# Gitflow + Conventional Commits + Semantic Versioning

## Branches permanentes

| Branch | Proposito |
|--------|-----------|
| `main` | Producao estavel — apenas merges via PR de `develop` ou `hotfix/*` |
| `develop` | Integracao continua — merges de `feature/*` e `fix/*` |

## Nomenclatura de branches por fase

| Fase | Branch | Tipo de commit | Bump |
|------|--------|---------------|------|
| Fase 1 — Scaffolding + dominio | `feature/scaffolding-domain` | `feat(domain)` | minor |
| Fase 2 — UseCase ListOrders | `feature/usecase-list-orders` | `feat(usecase)` | minor |
| Fase 3A — REST GET /order | `feature/rest-list-orders` | `feat(rest)` | minor |
| Fase 3B — gRPC ListOrders | `feature/grpc-list-orders` | `feat(grpc)` | minor |
| Fase 3C — GraphQL listOrders | `feature/graphql-list-orders` | `feat(graphql)` | minor |
| Fase 4 — Wire DI update | `feature/wire-di-update` | `feat(di)` | minor |
| Fase 5 — Docker automation | `feature/docker-automation` | `feat(infra)` | minor |
| Fase 6 — api.http + README | `docs/api-readme` | `docs` | patch |
| Fase 7 — Revisao e qualidade | `fix/<issue>` | `fix` / `refactor` | patch |

## Fluxo por fase

```
develop
  └── feature/scaffolding-domain
        └── commits pequenos e atomicos
        └── PR -> develop (merge --no-ff)
              └── tag v0.1.0

develop
  └── feature/usecase-list-orders
        └── PR -> develop
              └── tag v0.2.0

... (idem para 3A, 3B, 3C em paralelo em branches separadas)

develop
  └── feature/docker-automation
        └── PR -> develop
              └── tag v0.7.0

develop -> main (release)
  tag v1.0.0
```

## Convencao de commits (Conventional Commits)

```
<tipo>(<escopo>): <descricao curta em minusculas>

[corpo opcional — bullet points com detalhes]

[rodape — Co-Authored-By, BREAKING CHANGE, etc]
```

### Tipos validos

| Tipo | Quando usar |
|------|-------------|
| `feat` | Nova funcionalidade |
| `fix` | Correcao de bug |
| `docs` | Apenas documentacao |
| `refactor` | Refatoracao sem mudanca funcional |
| `test` | Adicionar/corrigir testes |
| `chore` | Manutencao, deps, configs, CI |
| `perf` | Melhoria de performance |
| `ci` | Pipeline, Docker, automacao |

### Escopos por camada

| Escopo | Camada |
|--------|--------|
| `domain` | `internal/entity/` |
| `usecase` | `internal/usecase/` |
| `repository` | `internal/infra/database/` |
| `rest` | `internal/infra/web/` |
| `grpc` | `internal/infra/grpc/` |
| `graphql` | `internal/infra/graph/` |
| `di` | `cmd/ordersystem/wire*.go` |
| `infra` | `Dockerfile`, `docker-compose`, `migrations/` |
| `config` | `configs/`, `.env` |

### Exemplos

```
feat(domain): add Order entity with CalculateFinalPrice method
feat(usecase): implement ListOrdersUseCase with repository interface
feat(repository): add GetAll method to MySQL order repository
feat(rest): add GET /order endpoint for listing orders
feat(grpc): add ListOrders RPC to OrderService proto
feat(graphql): add listOrders query to schema and resolver
feat(di): wire ListOrdersUseCase across all three handlers
feat(infra): add Dockerfile multi-stage and docker-compose with healthcheck
docs: add api.http with create and list examples for all protocols
```

## Semantic Versioning

```
v<MAJOR>.<MINOR>.<PATCH>
```

| Evento | Bump | Exemplo |
|--------|------|---------|
| Breaking change (`feat!` / `BREAKING CHANGE`) | MAJOR | v1.0.0 -> v2.0.0 |
| Nova feature (`feat`) | MINOR | v0.1.0 -> v0.2.0 |
| Bug fix, docs, refactor, chore | PATCH | v0.1.0 -> v0.1.1 |

### Versoes planejadas

| Versao | Marco |
|--------|-------|
| v0.1.0 | Scaffolding + entidade Order |
| v0.2.0 | UseCase ListOrders + repository GetAll |
| v0.3.0 | REST GET /order |
| v0.4.0 | gRPC ListOrders |
| v0.5.0 | GraphQL listOrders |
| v0.6.0 | Wire DI integrado nos 3 protocolos |
| v0.7.0 | Docker automation completo |
| v0.8.0 | api.http + README |
| v1.0.0 | Release: todos os criterios de aceite passando |

## Regras criticas

- NUNCA fazer commit diretamente em `main` ou `develop`
- NUNCA usar `git push --force` em `main` ou `develop`
- NUNCA incluir `.env`, secrets, credenciais no commit
- NUNCA usar `--no-verify`
- Cada PR deve ter ao menos 1 criterio de aceite verificado antes do merge
- Tags sao criadas sempre em `develop` ou `main`, nunca em feature branches
