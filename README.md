# Sistema de Auditoria Distribuída

Plataforma em Go para processar transações financeiras de forma assíncrona, com auditoria, rastreabilidade e persistência em banco e blockchain privada.

## O que o projeto entrega

- API HTTP para criação e consulta de transações.
- Processamento assíncrono com worker separado da API.
- Registro de eventos e auditoria para conformidade.
- Autenticação com JWT e hashing seguro de senha.
- Camadas bem separadas para facilitar manutenção e testes.
- Suporte a fila em memória e base para RabbitMQ.
- Infra pronta para Docker Compose.

## Qual dor ele resolve

Sistemas financeiros e de auditoria normalmente precisam conciliar três exigências ao mesmo tempo: resposta rápida para o usuário, rastreabilidade completa das operações e confiabilidade dos registros. Quando tudo acontece no fluxo síncrono, a API fica lenta, difícil de escalar e mais sujeita a falhas em cadeia.

Este projeto resolve isso separando a entrada da operação do seu processamento pesado. A API recebe a solicitação, persiste o evento e delega o processamento ao worker. Assim, a aplicação mantém agilidade sem abrir mão de auditoria e controle.

## Benefícios

- Menor latência na API.
- Melhor escalabilidade para volumes altos de transações.
- Rastreabilidade das operações do início ao fim.
- Arquitetura mais simples de testar e manter.
- Facilidade para evoluir a fila, a blockchain ou os provedores de notificação sem refatoração grande.
- Base clara para ambientes locais, homologação e produção.

## Como o sistema funciona

1. O cliente envia uma transação para a API.
2. A aplicação valida, grava e enfileira a operação.
3. O worker consome a mensagem e executa o processamento.
4. O sistema atualiza o estado, registra auditoria e envia notificações.

## Estrutura principal

- `cmd/api` - entrypoint da API.
- `cmd/worker` - entrypoint do processador assíncrono.
- `internal/domain` - entidades e contratos do domínio.
- `internal/usecase` - regras de negócio.
- `internal/infra` - banco, fila, autenticação, blockchain e HTTP.
- `internal/container` - composição de dependências.
- `config` - carregamento e validação de configuração.
- `docs` - documentação técnica detalhada.

## Começar rápido

```bash
cp .env.example .env
docker-compose up -d
```

Ou rodando localmente:

```bash
go run ./cmd/api
go run ./cmd/worker
```

## Documentação

- [Arquitetura](docs/ARCHITECTURE.md)
- [Guia geral](docs/README.md)
- [Quickstart](docs/QUICKSTART.md)
- [Integração com API](docs/API_INTEGRATION_GUIDE.md)
- [Testes](docs/TESTING_GUIDE.md)
- [Deploy](docs/DEPLOYMENT_GUIDE.md)
- [Segurança](docs/SECURITY_GUIDE.md)
