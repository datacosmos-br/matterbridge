# Upstream tracking (datacosmos fork)

State on 2026-11-07, after PR #3 integrated everything into `master`. This file is the
record of what is still external: branches proposed (or proposable) to
matterbridge-org/matterbridge, related upstream issues, and known follow-ups.

## Upstream conversation

- matterbridge-org/matterbridge#270: our first-contact issue (AI usage disclosed).
  A maintainer replied that it "smells automated agent authored". The answer must come
  from a human maintainer of this fork; no upstream PR is opened before upstream agrees.

## Branches kept for upstream PRs

Each is one commit on upstream master, already merged into our `master`. Open a PR with
`gh pr create -R matterbridge-org/matterbridge --base master --head datacosmos-br:<branch>`
only if upstream accepts it. When upstream merges or rejects one, delete the branch.

| Branch | Change | Upstream reference |
|---|---|---|
| `upstream/config-stringslice2d-absent-v1` | stop logging unset 2D string settings as errors | #256 |
| `upstream/whatsappmulti-nil-contextinfo-v3` | fix nil-pointer panics in whatsappmulti | #264 |
| `upstream/discord-webhook-delete-v2` | delete webhook-sent messages with the webhook token | #42 |
| `upstream/discord-forum-threads-v6` | Discord forum threads via `ForumID` (stacked on the delete branch) | 42wim#1914 |
| `upstream/whatsappmulti-stickers-v2` | relay WhatsApp stickers | — |
| `upstream/whatsappmulti-reactions-v2` | relay WhatsApp reactions | #189 (draft) |

Fork-only (not proposed): persistent message map and opt-out (`gateway/messagemap.go`,
`gateway/optout.go`), cross-platform delete propagation, Dify summaries and identity
(`gateway/dify`, `docs/dify.md`), governance and `make` tooling.

Adopted from upstream ahead of merge: #139 (matrix), #170 (discord attachments). Drop the
duplicates when upstream merges them.

## Known follow-ups

- Graceful shutdown: `matterbridge.go` has no shutdown path, so `Router.Close()` only runs
  when startup fails; add signal handling (stop bridges, then close stores).
- Discord: messages sent through a webhook with file uploads cannot be deleted through the
  webhook, because their IDs are discarded (`webhookSendFilesOnly`).
- Teams meeting webhook and Google Meet: no consumer; the old scaffolds remain only in
  `legacy/datacosmos-2026-11`.
- Deployment: no local project deploys matterbridge yet (no chart, gitops manifest or unit).
  Production needs the cluster lane (cosmos-charts/cosmos-gitops), Dify credentials and a
  WhatsApp session.
- Harvest source: the 42wim-era contributor branches that still had unique commits are
  listed in the legacy record; their commits remain fetchable by SHA from the GitHub fork
  network, and the old fork is frozen as `legacy/datacosmos-2026-11`
  (tag `legacy-datacosmos-2026-11-07`).
