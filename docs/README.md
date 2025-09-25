# Sistema de Auditoria Distribuída para Transações Financeiras

## 1. Introdução

### 1.1. Propósito

Este documento descreve os requisitos funcionais e não funcionais, regras de negócio, regras de acesso, atores, casos de uso e entidades do sistema **Sistema de Auditoria Distribuída para Transações Financeiras**.

### 1.2. Escopo

O sistema processa transações financeiras de forma assíncrona, registrando cada etapa em uma blockchain privada para auditoria completa, rastreabilidade e compliance financeiro.

### 1.3. Público-Alvo

- Desenvolvedores
- Auditores financeiros
- Administradores do sistema
- Gestores de projetos

---

## 2. Definição de Entidades

- **Usuário** – pessoa que utiliza o sistema (Cliente, Auditor ou Administrador).
- **Transação** – operação financeira com origem, destino, valor e status.
- **Bloco** – registro na blockchain contendo transações e hash anterior.
- **Notificação** – mensagens de status enviadas ao cliente sobre transações.

---

## 3. Atores

- **Cliente** – envia transações e recebe status.
- **Auditor** – consulta histórico completo e gera relatórios.
- **Administrador** – gerencia filas, workers, nós da blockchain e configurações do sistema.
- **Sistema Externo (Blockchain)** – armazena de forma imutável os registros de transações.
- **Fila de Processamento** – gerencia tarefas assíncronas.

---

## 4. Requisitos Funcionais

- **RF01**: Receber transações via gRPC.
- **RF02**: Validar dados e gerar ID único da transação.
- **RF03**: Processar transações de forma assíncrona.
- **RF04**: Validar regras financeiras (saldo, limites, integridade).
- **RF05**: Registrar transações na blockchain privada.
- **RF06**: Consultar histórico de transações.
- **RF07**: Notificar status das transações.

---

## 5. Requisitos Não Funcionais

- **RNF01**: Processamento paralelo com goroutines.
- **RNF02**: Escalabilidade para absorver picos de transações.
- **RNF03**: Comunicação gRPC segura (TLS).
- **RNF04**: Retry automático em falhas.
- **RNF05**: Auditoria com histórico imutável.
- **RNF06**: Alta disponibilidade.

---

## 6. Regras de Negócio

- **RN01**: Cada transação deve ter um ID único e rastreável.
- **RN02**: Saldo insuficiente impede processamento e registra falha.
- **RN03**: Apenas usuários autenticados podem enviar transações ou consultar histórico.
- **RN04**: Transações falhadas devem ser reprocessadas até limite definido.
- **RN05**: Blockchain deve registrar hash da transação anterior para integridade.
- **RN06**: Notificações são enviadas somente após validação ou finalização.

---

## 7. Regras de Acesso

- **Cliente**: enviar transações e consultar histórico próprio.
- **Auditor**: consultar todas as transações e gerar relatórios.
- **Administrador**: gerenciar filas, workers, nós da blockchain e configurações.

---

## 8. Casos de Uso

### UC01 – Enviar Transação

**Ator:** Cliente

### UC02 – Receber Status da Transação

**Ator:** Cliente

### UC03 – Auditar Histórico

**Ator:** Auditor

### UC04 – Gerenciar Sistema

**Ator:** Administrador

---

## 9. Diagramas UML

### 9.1. Diagrama de Casos de Uso

```mermaid
usecaseDiagram
  actor Cliente
  actor Auditor
  actor Administrador

  usecase UC01 as "Enviar Transação"
  usecase UC02 as "Receber Status"
  usecase UC03 as "Auditar Histórico"
  usecase UC04 as "Gerenciar Sistema"

  Cliente --> UC01
  Cliente --> UC02
  Auditor --> UC03
  Administrador --> UC04
```

## 9.2. Diagrama de Classes

    ````classDiagram

class Usuario {
+id: int
+nome: string
+tipo: string
+credenciais: string
+autenticar()
}

class Transacao {
+id: int
+origem: string
+destino: string
+valor: float
+status: string
+timestamp: datetime
+validar()
+processar()
}

class Bloco {
+hash: string
+transacoes: list<Transacao>
+hashAnterior: string
+timestamp: datetime
+assinar()
}

class Notificacao {
+id: int
+transacao_id: int
+mensagem: string
+status: string
+timestamp: datetime
+enviar()
}

Usuario "1" --> "N" Transacao
Transacao "1" --> "1" Bloco
Transacao "1" --> "N" Notificacao

``````



## 9.3. Diagrama de Sequência – Fluxo Assíncrono

    ````sequenceDiagram
    Cliente ->> TransactionService: Envia Transação (gRPC)
    TransactionService ->> TransactionService: Valida Transação
    TransactionService ->> TransactionService: Gera ID Único
    TransactionService ->> Queue: Envia Transação
    TransactionService -->> Cliente: Retorna ID

        Queue ->> WorkerPool: Processa Transação
        WorkerPool ->> WorkerPool: Valida Saldo/Regras
        WorkerPool ->> BlockchainWriter: Prepara dados
        BlockchainWriter ->> Blockchain: Grava Bloco
        BlockchainWriter -->> WorkerPool: Confirmação
        WorkerPool ->> NotificationService: Envia Status
        NotificationService -->> Cliente: Notificação
    `````
``````
