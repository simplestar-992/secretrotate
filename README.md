# SecretRotate | Secrets Rotation Manager

<p align="center">
  <img src="https://img.shields.io/badge/Security-Secrets%20Rotation-9B59B6?style=for-the-badge" alt=""/>
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt=""/>
</p>

---

### Automate your secrets rotation

Rotate API keys, passwords, and tokens without downtime. SecretRotate handles the complexity of updating running services.

```bash
secretrotate rotate --service api-key --env production
```

---

## How It Works

1. Generate new secret
2. Update secrets store (Vault, AWS, etc.)
3. Roll out to services gradually
4. Verify service health
5. Revoke old secret

---

## Features

- 🔄 **Automatic rotation** - Set it and forget it
- 🔍 **Health checks** - Verify services after rotation
- 📊 **Audit logging** - Track all secret changes
- 🌐 **Multi-backend** - Vault, AWS Secrets, Azure Key Vault
- ⏰ **Scheduled rotation** - Cron-based automation

---

## Usage

```bash
# Rotate a secret
secretrotate rotate -s api-key -e prod

# List secrets
secretrotate list

# View rotation status
secretrotate status -s database-password

# Configure rotation schedule
secretrotate schedule -s api-key -every 30d
```

---

## Supported Backends

| Provider | Status |
|----------|--------|
| AWS Secrets Manager | ✅ |
| HashiCorp Vault | ✅ |
| Azure Key Vault | ✅ |
| GCP Secret Manager | 🚧 |
| Local encrypted store | ✅ |

---

MIT © 2024 [simplestar-992](https://github.com/simplestar-992)
