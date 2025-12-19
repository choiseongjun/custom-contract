#!/bin/bash

set -e

echo "=== Simple CRUD Contract Test ==="
echo ""

# Store contract
echo "1. Storing contract..."
STORE_TX=$(scontractd tx wasm store contracts/simple-crud/target/wasm32-unknown-unknown/release/simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json)

CODE_ID=$(echo $STORE_TX | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')

echo "   ✓ Code ID: $CODE_ID"
sleep 6

# Instantiate contract
echo "2. Instantiating contract..."
INIT_TX=$(scontractd tx wasm instantiate $CODE_ID '{}' \
  --from alice \
  --label "simple-crud" \
  --no-admin \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json)

CONTRACT=$(echo $INIT_TX | jq -r '.logs[0].events[] | select(.type=="instantiate") | .attributes[] | select(.key=="_contract_address") | .value')

echo "   ✓ Contract: $CONTRACT"
sleep 6

# Create data
echo "3. Creating data (name=Alice)..."
scontractd tx wasm execute $CONTRACT \
  '{"create":{"key":"name","value":"Alice"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json > /dev/null

sleep 6

# Query data
echo "4. Reading data..."
RESULT=$(scontractd query wasm contract-state smart $CONTRACT \
  '{"read":{"key":"name"}}' \
  --chain-id scontract \
  --output json)

echo "   ✓ Result: $(echo $RESULT | jq -r '.data.value')"

# Update data
echo "5. Updating data (name=Bob)..."
scontractd tx wasm execute $CONTRACT \
  '{"update":{"key":"name","value":"Bob"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json > /dev/null

sleep 6

# Query updated data
echo "6. Reading updated data..."
RESULT=$(scontractd query wasm contract-state smart $CONTRACT \
  '{"read":{"key":"name"}}' \
  --chain-id scontract \
  --output json)

echo "   ✓ Result: $(echo $RESULT | jq -r '.data.value')"

# Delete data
echo "7. Deleting data..."
scontractd tx wasm execute $CONTRACT \
  '{"delete":{"key":"name"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json > /dev/null

sleep 6

# Query after delete
echo "8. Reading after delete..."
RESULT=$(scontractd query wasm contract-state smart $CONTRACT \
  '{"read":{"key":"name"}}' \
  --chain-id scontract \
  --output json)

echo "   ✓ Result: $(echo $RESULT | jq -r '.data.value')"

echo ""
echo "=== Test Complete! ==="
echo "Contract Address: $CONTRACT"
