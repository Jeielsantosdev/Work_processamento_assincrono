# Sistema de Auditoria Distribuída para Transações Financeiras

## 📌 Visão Geral

O **Sistema de Auditoria Distribuída** processa transações financeiras de forma assíncrona, registrando cada etapa em uma **blockchain privada**.\
Cada transação recebe um identificador único, permitindo rastreabilidade completa, auditoria e compliance financeiro.\
O sistema garante **performance**, **confiabilidade** e **segurança**, sendo ideal para bancos, fintechs e empresas que precisam de monitoramento detalhado de transações.

---

## 🎯 Objetivos

- Processar transações financeiras de forma assíncrona.\
- Registrar todas as transações em blockchain privada para auditoria.\
- Notificar clientes sobre o status das transações (sucesso, falha, pendente).\
- Permitir consultas históricas e relatórios para auditores.\
- Garantir **confiabilidade**, **integridade** e **escabilidade** do sistema.\
- Fornecer interface simples para administradores monitorarem o sistema.

---

## 🛠️ Tecnologias Utilizadas

- **Backend:** Go (para APIs, processamento assíncrono e workers).\
- **Frontend:** React.js + TailwindCSS (painel administrativo e dashboards).\
- **Banco de Dados:** PostgreSQL (produção) e SQLite (desenvolvimento).\
- **Blockchain Privada:** Para armazenamento imutável das transações.\
- **Filas Assíncronas:** RabbitMQ ou Redis Streams para processamento de transações.\
- **Autenticação:** JWT (para clientes e auditores).\
- **Testes:** Go testing (backend) + Jest (frontend).

---

## 🏛️ Arquitetura do Sistema

O sistema segue uma **arquitetura distribuída em camadas**:

1.  **Frontend (React + TailwindCSS):** Painel administrativo e dashboards de auditoria, consulta de histórico e status das transações.\
2.  **Backend (Go):** API gRPC para recebimento de transações, workers para processamento assíncrono e integração com blockchain.\
3.  **Fila de Processamento:** Gerencia tarefas assíncronas e picos de demanda.\
4.  **Blockchain Privada:** Armazena os blocos com metadados e hash das transações para garantir imutabilidade.\
5.  **Camada de Autenticação (JWT):** Controle de acesso para clientes, auditores e administradores.\

---

## 📐 Metodologia de Desenvolvimento

- **Metodologia:** Ágil (Scrum).\
- **Controle de Versão:** Git + GitHub.\
- **Integração Contínua (CI/CD):** GitHub Action\
- **Code Review:** Pull Requests obrigatórios.\
- **Documentação:** Markdown + Swagger/OpenAPI para API.

---

## 📂 Estrutura do Projeto

    /sistema-auditoria
    │── backend/                  # Código em Go
    │   ├── cmd/                  # Pontos de entrada da aplicação
    │   ├── internal/             # Lógica de negócio, workers e serviços
    │   ├── pkg/                  # Pacotes utilitários
    │   └── tests/                # Testes unitários e de integração
    │
    │── frontend/                 # Código em React + Tailwind
    │   ├── src/
    │   │   ├── components/       # Componentes reutilizáveis
    │   │   ├── pages/            # Dashboards e páginas de auditoria
    │   │   ├── hooks/            # Hooks customizados
    │   │   ├── services/         # Comunicação com a API
    │   │   └── styles/           # Estilização com Tailwind
    │
    │── docs/                     # Documentação do projeto
    │── docker-compose.yml        # Configuração de containers
    │── README.md                 # Guia rápido do projeto

---

## 🚀 Fluxo Básico do Sistema

1.  Cliente envia transação via **gRPC**.\
2.  **Transaction Service** valida os dados e gera um ID único.\
3.  Transação é enviada para a **fila assíncrona**.\
4.  **Worker Pool** processa a transação, valida regras financeiras e prepara dados para blockchain.\
5.  **Blockchain Writer** cria bloco com hash anterior e grava na blockchain privada.\
6.  **Notification Service** envia status da transação ao cliente.\
7.  Auditor ou administrador pode consultar histórico ou gerar relatórios a qualquer momento.

---

## 🔒 Segurança

- Uso de **TLS/HTTPS** em toda a comunicação.\
- Tokens **JWT** para autenticação e autorização.\
- Filas e workers isolados para impedir execução não autorizada.\
- Sanitização e validação de todas as entradas de dados.\
- Registro imutável de transações na blockchain para auditoria confiável.

---
