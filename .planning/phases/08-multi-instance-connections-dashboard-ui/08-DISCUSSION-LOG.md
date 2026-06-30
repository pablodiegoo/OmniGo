# Phase 8: Multi-Instance Connections & Dashboard UI - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-06-29
**Phase:** 8-multi-instance-connections-dashboard-ui
**Areas discussed:** Transição de Banco de Dados, Fluxo de Re-pareamento de WhatsApp Web, Limite de Conexões por Workspace

---

## Transição de Banco de Dados (Tabelas Legadas vs Nova tabela 'connections')

| Option | Description | Selected |
|--------|-------------|----------|
| Opção A (Recomendada) | Criar uma migração Goose (SQL) que migra automaticamente todos os dados existentes das tabelas antigas para a tabela connections e as remove. | ✓ |
| Opção B | Criar a tabela connections limpa e desconsiderar os dados antigos. | |

**User's choice:** Opção A
**Notes:** Garante transição de dados suave e preservação de todas as credenciais em ambientes locais de desenvolvimento/produção.

---

## Fluxo de Re-pareamento de WhatsApp Web (whatsmeow desconectado)

| Option | Description | Selected |
|--------|-------------|----------|
| Opção A (Recomendada) | Exibir status 'Desconectado' no card do canal e botão 'Re-parear' abrindo modal para gerar novo QR Code sob o mesmo ID de conexão, mantendo logs históricos vinculados. | ✓ |
| Opção B | Excluir a conexão automaticamente e forçar o operador a criar uma nova do zero. | |

**User's choice:** Opção A
**Notes:** Mantém a rastreabilidade e integridade dos logs históricos e estatísticas vinculadas ao remetente mesmo após desconexões físicas do chip.

---

## Limite de Conexões por Workspace (Segurança e Controle de RAM)

| Option | Description | Selected |
|--------|-------------|----------|
| Opção A (Recomendada) | Limite máximo configurável por env var (PERGO_MAX_WHATSAPP_CONNECTIONS=5) para WhatsApp Web (whatsmeow), bloqueando com HTTP 422 se violado. Canais stateless continuam ilimitados. | ✓ |
| Opção B | Deixar totalmente ilimitado, permitindo que a CPU/RAM do servidor atinja o limite físico. | |

**User's choice:** Opção A
**Notes:** Protege servidores pequenos/VPS contra vazamentos de memória e erros de Out of Memory (OOM).

---

## the agent's Discretion
- Seleção e implementação de bibliotecas e plugins daisyUI adicionais no frontend.
- Detalhes de tratamento e checagem de erros SQL específicos durante a migração Goose.

## Deferred Ideas
- stripe/billing managed cloud backend (recorded in seed `cpaas-billing-saas-infra.md`).
