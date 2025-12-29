# cloak
Encrypt on disk. Inject in RAM. Invisible everywhere else.

[Installation](https://google.com/) • [Quick Start](https://google.com/) • [Features](https://google.com/) • [Security](https://google.com/)

## Problem
You are building modern software, but you are managing secrets like it's 1999.
- **The ENV struggle** - You can't commit them. You have to DM them to new hires. They get lost or even accidentally leaked.
- **Overkill cloud vaults** - You don't need to pay $500/mo enterprise license for a company that rhymes with BashiBorp Gault just to run a simple Next.js app

## Solution
**cloak** is a developer first secrets manager that lives right in your terminal. It encrypts your environment variables into a binary file that you ***can*** commit to Git.

When you run your app, cloak basically acts as a "ghost shell." It decrypts your secrets right into the process memory. So, **no plaintext files ever touch your disk**.

## Installation

### via Go (recommended)
```bash
go install github.com/atomisadev/cloak@latest
```

### via Homebrew
```bash
brew tap atomisadev/tap
brew install cloak
```

## Quick Start
### 1. Initialize the vault
Create a new encrypted store in your project root.
```
$ cloak init
> Generating Master Key...
> MASTER KEY: xxxx-yyyy-zzzz (Save this!)
> Created cloak.encrypted
```

### 2. Add a secret
Don't exit text files. Add secrets securely.
```
$ cloak set STRIPE_KEY sk_text_51Mz...
> Encrypting... [ OK ]
```

### 3. Run "Ghost Mode"
Run your project. Cloak will inject the variables directly into the child process.
```
# before: bun run dev (it will crash because you are missing envs)
# after:
$ cloak run -- bun run dev

> [CLOAK] Injecting 12 secrets into process...
> ready - started server on 0.0.0.0:3000, url: http://localhost:3000
```

## Powerful Features
### Slick TUI (`cloak edit`)
Don't like CLI flags? You can launch the interactive "Deck" to manage secrets visually with a clean interface.
- Vim-style navigation (`j`/`k`)
- Masking toggle (`h` to hide/show values)
- Audit metadata (see who last modified a key)

### Dead Drop Sharing (`cloak share`)
Need to give the Master Key to a new team member? Don't paste it in your Slack. Instead, use Cloak to generate a Zero-Knowledge one-time URL. The server sees the encrypted blob, but the decrypted key is in the URL hash fragment (which is never sent to the server).
```
$ cloak share
> Encrypting Master Key...
> DEAD DROP: https://getcloak.xyz/drop/a1b2-c3d4#FRAGMENT_KEY
> This link will self-destruct after 1 view.
```

### Polyglot Intellisense (`cloak types`)
Don't guess variable names. Cloak reads your encrypted vault and generates type definitions for your IDE.
```
# for TypeScript/Next.js
$ cloak types --lang=ts
> Generated env.d.ts

# For Python/Pydantic
$ cloak eject --lang=python
> Generated config.py
```

### Conflict Resolution (`cloak fix-merge`)
Binary conflict in Git?
```
$ git pull
> CONFLICT (content): Merge conflict in cloak.encrypted

$ cloak fix-merge
> Decrypting LOCAL...
> Decrypting REMOTE...
> Smart Merging keys...
> Re-encrypting... [ FIXED ]
```

## Security Architecture
Cloak is built on the philosophy of **Trust No One**.
- **AES-256-GCM** - Industry standard for authenticated encryption. Used for the `cloak.encrypted` file.
- **Zero Knowledge Sharing** - The `cloak share` command uses client side encryption. The server hosting the "Dead Drop" can't read your keys at all.
- **Memory Only Injection** - Secrets are decrypted into RAM and pased directly to the `syscall.Exec` environment. They are never written to a temporary files (preventing attacks via `/tmp` scanning)

## Under Development
This project is still under development, and so is not yet installable.
