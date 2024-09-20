#!/bin/bash

# Set DOCKER to 'docker' if it is not set
: ${DOCKER:=docker}

# Set DEBUG to 'true' to enable debugging
: ${DEBUG:=false}

[[ "${DEBUG}" == 'true' ]] && echo Debugging is enabled

VAULT_ADDR='http://127.0.0.1:8200'
VAULT_TOKEN='root'
VAULT_CONTAINER_NAME='my-vault'
VAULT_IMAGE_NAME='hashicorp/vault:latest'

# Function to clean up the Vault containerreconcileRawDataLanding
cleanup() {
    echo "Cleaning up the Vault container..."
    ${DOCKER} rm -f ${VAULT_CONTAINER_NAME}
    exit 0
}

# Check for the "clean" argument
if [ "$1" == "clean" ]; then
    cleanup
fi

# Create the Vault container if it does not exist
if [ ! $(${DOCKER} ps -aq -f name=my-vault) ]; then
    ${DOCKER} run --name ${VAULT_CONTAINER_NAME} -d -p 8200:8200 \
     -e "VAULT_DEV_ROOT_TOKEN_ID=${VAULT_TOKEN}" \
     -e 'VAULT_LOCAL_CONFIG={
                              "backend": {
                                            "file": {"path": "/vault/file"}
                                          },
                              "default_lease_ttl": "168h",
                              "max_lease_ttl": "720h"}' \
                              ${VAULT_IMAGE_NAME}
else
    echo "Vault container already exists"
fi

# Vault Commands (using docker exec)
vault_cmd() {
    ${DOCKER} exec ${VAULT_CONTAINER_NAME} sh -c \
    "export VAULT_ADDR='http://0.0.0.0:8200' && \
    export VAULT_TOKEN='$VAULT_TOKEN' && \
    vault $*"
}

# Enable the KV secrets engine at the path "apps/"
if ! vault_cmd "secrets list | grep -q 'apps/'"; then
    echo "Enabling KV secrets engine at 'apps/' path..."
    vault_cmd "secrets enable -path=apps kv"
else
    echo "KV secrets engine at 'apps/' path is already enabled"
fi

# Define app-policy for the apps/ path
app_policy=$(cat <<EOF
path "apps/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}
EOF
)
# Write app-policy
vault_cmd "policy write app-policy - <<'EOF'
${app_policy}
EOF
"
git
# Read the app-policy
[[ "${DEBUG}" == 'true' ]] && vault_cmd "policy read app-policy"

# Enable AppRole Authentication (only if not already enabled)
if ! vault_cmd "auth list | grep -q approle"; then
    echo "Enabling approle auth method..."
    vault_cmd "auth enable approle"
else
    echo "AppRole authentication is already enabled"
fi

# Create AppRole Role
vault_cmd "write auth/approle/role/my-role policies=\"app-policy\" secret_id_ttl=0"
[[ "${DEBUG}" == 'true' ]]  && vault_cmd "read auth/approle/role/my-role "

# Get Role ID and Secret ID
export VAULT_ROLE_ID=$(vault_cmd "read -field=role_id auth/approle/role/my-role/role-id")
export VAULT_SECRET_ID=$(vault_cmd "write -f -field=secret_id auth/approle/role/my-role/secret-id")

sed -i -e "s|role-id: '.*'|role-id: '$VAULT_ROLE_ID'|" config.yaml
sed -i -e "s|secret-id: '.*'|secret-id: '$VAULT_SECRET_ID'|" config.yaml
sed -i -e "s|address: '.*'|address: '$VAULT_ADDR'|" config.yaml

[[ "${DEBUG}" == 'true' ]] && cat config.yaml

[[ "${DEBUG}" == 'true' ]]  && echo "role-id: $VAULT_ROLE_ID"
[[ "${DEBUG}" == 'true' ]]  && echo "secret-id: $VAULT_SECRET_ID"

[[ "${DEBUG}" == 'true' ]] && export VAULT_TOKEN=$(vault_cmd "write -field=token auth/approle/login role_id=$VAULT_ROLE_ID secret_id=$VAULT_SECRET_ID")
[[ "${DEBUG}" == 'true' ]]  && echo "token: $VAULT_TOKEN"

# write a secret to the apps/ path
[[ "${DEBUG}" == 'true' ]] && vault_cmd "kv put -mount=apps data/ddis-asteroid-infra/dice/foo foo=bar"
[[ "${DEBUG}" == 'true' ]] && vault_cmd "kv get -mount=apps data/ddis-asteroid-infra/dice/foo"



