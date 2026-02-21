# Money Savior - Telegram Bot

<div align="center">

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=FF9900)
![DynamoDB](https://img.shields.io/badge/DynamoDB-527FFF?style=for-the-badge&logo=amazon-dynamodb&logoColor=white)
![Telegram](https://img.shields.io/badge/Telegram-26A5E4?style=for-the-badge&logo=telegram&logoColor=white)
![Lambda](https://img.shields.io/badge/AWS%20Lambda-FF9900?style=for-the-badge&logo=aws-lambda&logoColor=white)

Um **assistente pessoal de controle de gastos** integrado ao Telegram com suporte a IDs sequenciais, navegação inteligente e gerenciamento completo de despesas.

[Repositório](#) • [Documentação](#documentação) • [Quick Start](#-quick-start) • [Comandos](#-comandos-disponíveis)

</div>

---

## Features

[+] **Registro de Gastos** - Registre despesas com valor, categoria e método de pagamento  
[+] **IDs Sequenciais** - Gastos salvos automaticamente com IDs em ordem (1, 2, 3...)  
[+] **Consulta Inteligente** - Liste todos os gastos ou veja um específico com navegação  
[+] **Navegação** - Botões para mover entre registros  
[+] **Delete com Confirmação** - Apague um gasto específico ou todos com confirmação inline  
[+] **Banco de Dados Cloud** - DynamoDB da AWS para armazenamento seguro  
[+] **Serverless** - Execução via AWS Lambda para escalabilidade  

---

## Arquitetura

```
┌─────────────────────────────────────────────────────────┐
│                  Telegram Bot (Go)                       │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  Bot Routes / Message Handlers                      │ │
│  │  - /gastei  - /consulta  - /deletar               │ │
│  │  - /deletartudo  - /help  - /start                │ │
│  └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│           AWS Lambda + API Gateway                       │
│  (Opcional para deploy serverless)                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              AWS DynamoDB                               │
│  Tabela: expenses                                        │
│  PK: user_id | SK: expense_id                          │
│  Atributos: amount, category, method, seq_id           │
└─────────────────────────────────────────────────────────┘
```

---

## Estrutura do Projeto

```
the-money-savior-telegram-bot/
├── cmd/
│   ├── bot/
│   │   └── main.go              # Entry point do bot Telegram
│   └── lambda/
│       ├── main.go              # Handler AWS Lambda
│       └── deploy.sh            # Script de deploy
├── internal/
│   ├── bot/
│   │   └── telegram.go          # Roteamento de mensagens
│   ├── database/
│   │   └── dynamodb.go          # Integração DynamoDB
│   ├── handlers/
│   │   ├── start.go             # /start
│   │   ├── help.go              # /help
│   │   ├── expense.go           # /gastei
│   │   ├── query.go             # /consulta
│   │   ├── delete.go            # /deletar e /deletartudo
│   │   └── invalid.go           # Comando inválido
│   └── models/
│       └── expense.go           # Struct Expense
├── go.mod
├── go.sum
└── README.md
```

---

## Quick Start

### Pré-requisitos

- **Go** >= 1.23.0
- **AWS Account** com credenciais configuradas
- **Telegram Bot Token** (criar via [@BotFather](https://t.me/botfather))
- **DynamoDB Table** chamada `expenses`

### Instalação Local

```bash
# Clone o repositório
git clone https://github.com/AmrmDev/the-money-savior-telegram-bot.git
cd the-money-savior-telegram-bot

# Instale as dependências
go mod download

# Configure variáveis de ambiente
export TELEGRAM_BOT_TOKEN="seu_token_aqui"
export AWS_REGION="us-east-1"
export TABLE_NAME="expenses"

# Execute o bot
go run cmd/bot/main.go
```

### Deploy AWS Lambda

```bash
cd cmd/lambda
./deploy.sh
```

---

## Comandos Disponíveis

### Registrar Gasto
```
/gastei <valor> <categoria> [método]
```
**Exemplo:** `/gastei 45.50 supermercado débito`

### Consultar Gastos
```
/consulta                  # Lista todos os gastos
/consulta <ID>             # Vê detalhes de um gasto específico
```
**Recurso:** Navegação entre registros, visualize gastos sequencialmente

### Deletar Gasto
```
/deletar <ID>              # Deleta um gasto específico (com confirmação)
/deletartudo               # Deleta todos os registros (com confirmação)
```

### Ajuda
```
/help                      # Exibe todos os comandos
/start                     # Mensagem de boas-vindas
```

---

## Variáveis de Ambiente

| Variável | Descrição | Obrigatória |
|----------|-----------|------------|
| `TELEGRAM_BOT_TOKEN` | Token do bot Telegram | Sim |
| `TABLE_NAME` | Nome da tabela DynamoDB | Sim |
| `AWS_REGION` | Região AWS (padrão: us-east-1) | Não |

---

## Modelo de Dados

### Expense
```go
type Expense struct {
    UserID    int64     // ID do usuário Telegram
    ChatID    int64     // ID do chat
    Username  string    // Username do Telegram
    Amount    float64   // Valor do gasto
    Category  string    // Categoria do gasto
    Method    string    // Método de pagamento
    CreatedAt time.Time // Data de criação
    ExpenseID string    // SK DynamoDB (user_id#timestamp)
    SeqID     int       // ID sequencial (1, 2, 3...)
}
```

### DynamoDB Table Schema
```
Table Name: expenses
Primary Key:
  - Partition Key: user_id (Number)
  - Sort Key: expense_id (String)

Attributes:
  - seq_id: Number (índice para ordenação)
  - amount: Number
  - category: String
  - method: String
  - created_at: String (RFC3339)
  - username: String
  - chat_id: Number
```

---

## Fluxo de Operações

### Registrar Gasto
```
1. Usuário: /gastei 50.00 uber pix
   |
2. Bot: Valida formato e valor
   |
3. Bot: Busca próximo SeqID disponível
   |
4. Bot: Salva no DynamoDB
   |
5. Bot: Confirma a operação
```

### Consultar com Navegação
```
1. Usuário: /consulta 3
   |
2. Bot: Busca gasto com SeqID=3
   |
3. Bot: Exibe card com detalhes
   |
4. Usuário: Clica em Anterior ou Próximo
   |
5. Bot: Atualiza card inline (sem nova mensagem)
```

### Deletar com Confirmação
```
1. Usuário: /deletar 3
   |
2. Bot: Exibe detalhes + botões [Confirmar] [Cancelar]
   |
3. Usuário: Clica confirmar
   |
4. Bot: Deleta do DynamoDB
   |
5. Bot: Confirma e atualiza IDs (resequencialização automática na próxima consulta)
```

---

## Tecnologias

- **Linguagem:** Go 1.23.0
- **API Telegram:** [go-telegram-bot-api v5](https://github.com/go-telegram-bot-api/telegram-bot-api)
- **Banco de Dados:** AWS DynamoDB
- **Serverless:** AWS Lambda
- **SDK AWS:** [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)

---

## Dependências

```go
require (
    github.com/aws/aws-lambda-go v1.52.0
    github.com/aws/aws-sdk-go-v2 v1.41.1
    github.com/aws/aws-sdk-go-v2/config v1.32.9
    github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.32
    github.com/aws/aws-sdk-go-v2/service/dynamodb v1.55.0
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
)
```

---

## Segurança

- [x] **Validação de entrada** em todos os comandos
- [x] **Proteção por ID de usuário** (cada usuário vê apenas seus gastos)
- [x] **Confirmação obrigatória** para delete
- [x] **Credenciais AWS** via ambiente (nunca hardcoded)
- [x] **DynamoDB** com controle de acesso IAM

---

## Roadmap Futuro

- [ ] Resumo mensal de gastos
- [ ] Gráficos de categoria
- [ ] Exportar dados (CSV/PDF)
- [ ] Limite de gasto diário
- [ ] Alertas de overspending
- [ ] Integração com mais bancos de dados
- [ ] Modo compartilhado (múltiplos usuários)

---

## Logs

O bot mantém logs detalhados em stdout com prefixos:
- `[INFO]` - Informações gerais
- `[WARN]` - Avisos
- `[ERROR]` - Erros
- `[DEBUG]` - Debug detalhado

---

## Contribuições

Contribuições são bem-vindas! Por favor:

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra uma Pull Request

---

## Licença

Distribuído sob a licença MIT. Veja `LICENSE` para mais detalhes.

---

## Autores

**Thyago Toledo**
- Email: thyago10a2007@gmail.com
- GitHub: [@ThyagoToledo](https://github.com/ThyagoToledo)

**Armando dos Santos**
- GitHub: [@AmrmDev](https://github.com/AmrmDev)

---

<div align="center">

**[⬆ Voltar ao topo](#-money-savior---telegram-bot)**

Made with effort using Go + AWS

</div>
