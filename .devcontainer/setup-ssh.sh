#!/bin/bash
# Copia as chaves SSH do host para dentro do container com permissões corretas.
# Necessário porque o SSH rejeita arquivos cujo dono não é o usuário atual.
# O host monta ~/.ssh em /root/host-ssh (read-only); aqui copiamos e ajustamos.
set -e

mkdir -p /root/.ssh

for f in github_resumos github_resumos.pub config known_hosts; do
    [ -f "/root/host-ssh/$f" ] && cp "/root/host-ssh/$f" /root/.ssh/
done

chmod 700 /root/.ssh
[ -f /root/.ssh/github_resumos ]     && chmod 600 /root/.ssh/github_resumos
[ -f /root/.ssh/config ]             && chmod 600 /root/.ssh/config
[ -f /root/.ssh/github_resumos.pub ] && chmod 644 /root/.ssh/github_resumos.pub

# Garante que github.com está no known_hosts para não travar no primeiro push
ssh-keyscan -H github.com >> /root/.ssh/known_hosts 2>/dev/null || true

echo "SSH configurado com sucesso."
