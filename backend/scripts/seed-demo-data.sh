#!/usr/bin/env bash
# Walks realistic articles through the real HTTP API (not direct DB writes) so
# every dashboard has believable content across every pipeline stage. Safe to
# re-run: it only ever creates new articles, never deletes anything.
#
# Usage: API_BASE=http://localhost:8080/api/v1 ./scripts/seed-demo-data.sh
set -euo pipefail

API="${API_BASE:-http://localhost:8080/api/v1}"
BANNER_DIR="${BANNER_DIR:-/tmp/banners}"

login() {
  curl -sf -X POST "$API/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"email\":\"$1\",\"password\":\"password123\"}" | jq -r '.accessToken'
}

ADMIN_TOKEN=$(login admin@team1.blog)
MODERATOR_TOKEN=$(login moderator@team1.blog)
DESIGNER_TOKEN=$(login designer@team1.blog)
PUBLISHER_TOKEN=$(login publisher@team1.blog)
CHIDI_TOKEN=$(login contributor@team1.blog)
AMARA_TOKEN=$(login amara@team1.blog)
TUNDE_TOKEN=$(login tunde@team1.blog)

create_article() {
  local token="$1" title="$2" content="$3" citation="$4"
  curl -sf -X POST "$API/articles/" \
    -H "Authorization: Bearer $token" -H 'Content-Type: application/json' \
    -d "$(jq -n --arg t "$title" --arg c "$content" --arg s "$citation" \
      '{title:$t, content:$c, sourceCitation:$s}')" | jq -r '.id'
}

submit() {
  curl -sf -X POST "$API/articles/$1/submit" -H "Authorization: Bearer $2" -o /dev/null
}

review() {
  local article="$1" decision="$2" summary="$3" token="$4"
  curl -sf -X POST "$API/reviews/" \
    -H "Authorization: Bearer $token" -H 'Content-Type: application/json' \
    -d "$(jq -n --arg a "$article" --arg d "$decision" --arg s "$summary" \
      '{articleId:$a, decision:$d, summary:$s, suggestions:[]}')" -o /dev/null
}

upload_banner() {
  curl -sf -X POST "$API/banners/$1/upload" -H "Authorization: Bearer $DESIGNER_TOKEN" \
    -F "file=@$2" -o /dev/null
}

mark_ready() {
  curl -sf -X POST "$API/banners/$1/mark-ready" -H "Authorization: Bearer $DESIGNER_TOKEN" -o /dev/null
}

publish() {
  curl -sf -X POST "$API/articles/$1/publish" \
    -H "Authorization: Bearer $PUBLISHER_TOKEN" -H 'Content-Type: application/json' \
    -d "$(jq -n --arg u "$2" '{substackUrl:$u}')" -o /dev/null
}

release_payment() {
  curl -sf -X POST "$API/payments/$1/release" -H "Authorization: Bearer $ADMIN_TOKEN" -o /dev/null
}

echo "== 1/9 draft: Chidi, never submitted =="
create_article "$CHIDI_TOKEN" \
  "Understanding Avalanche's Subnet Architecture" \
  "Avalanche's subnet model lets teams launch application-specific blockchains with their own validator sets, gas tokens, and execution rules, while still settling back to the Primary Network for interoperability. This piece walks through how subnets differ from sharding or sidechains, why that matters for throughput, and what it takes to launch one today using Avalanche-CLI and the HyperSDK toolchain. We cover validator bonding requirements, the role of the P-Chain in subnet validation, and a few production subnets already live, including DeFi Kingdoms and Dexalot." \
  "https://docs.avax.network/subnets" > /dev/null

echo "== 2/9 draft: Amara, never submitted =="
create_article "$AMARA_TOKEN" \
  "DeFi on Avalanche: A Beginner's Guide" \
  "From Trader Joe to Benqi, Avalanche's DeFi ecosystem has grown into one of the most active on any L1. This guide covers how to bridge assets onto the C-Chain via Core, the difference between liquid staking and traditional staking on Avalanche, and a walkthrough of providing liquidity on a DEX for the first time, with screenshots and a note on impermanent loss risk for newcomers." \
  "https://www.avax.network/ecosystem" > /dev/null

echo "== 3/9 submitted: Tunde, awaiting moderator review =="
A3=$(create_article "$TUNDE_TOKEN" \
  "How Avalanche Consensus Achieves Sub-Second Finality" \
  "Avalanche consensus relies on repeated sub-sampled voting rather than a single global leader, letting validators reach finality in under a second without sacrificing the safety guarantees of classical consensus protocols. This article breaks down Snowman++ in plain language: what a validator actually does when it receives a new transaction, how conflicting transactions get resolved, and why this design scales validator count without the latency blowup you'd see in PBFT-style systems." \
  "https://docs.avax.network/consensus")
submit "$A3" "$TUNDE_TOKEN"

echo "== 4/9 changes_requested: Chidi =="
A4=$(create_article "$CHIDI_TOKEN" \
  "Top 5 Avalanche Subnets to Watch in 2026" \
  "Subnets continue to be Avalanche's signature differentiator. This roundup covers five subnets worth tracking this year across gaming, real-world assets, and institutional finance, with a short technical note on what makes each one's validator design distinct from a generic C-Chain deployment." \
  "https://subnets.avax.network")
submit "$A4" "$CHIDI_TOKEN"
review "$A4" "changes_requested" "Good topic, but each subnet needs at least 2 sentences on its actual validator/tokenomics setup, not just a one-line description. Also please cite a primary source per subnet rather than just the subnets.avax.network directory." "$MODERATOR_TOKEN"

echo "== 5/9 resubmitted: Amara, awaiting re-review =="
A5=$(create_article "$AMARA_TOKEN" \
  "NFTs on Avalanche: The Creator Economy" \
  "Avalanche's low fees and fast finality have made it a popular home for NFT marketplaces like Joepegs and Kalao. This piece looks at how creators are using Avalanche subnets for dedicated NFT drops, the role of royalties enforcement at the smart-contract level, and two case studies of artists who moved from Ethereum to Avalanche specifically to cut minting costs for their community." \
  "https://joepegs.com")
submit "$A5" "$AMARA_TOKEN"
review "$A5" "changes_requested" "Please add a section on royalty enforcement mechanics before resubmitting." "$MODERATOR_TOKEN"
# contributor addresses feedback and resubmits
curl -sf -X PUT "$API/articles/$A5" -H "Authorization: Bearer $AMARA_TOKEN" -H 'Content-Type: application/json' \
  -d "$(jq -n --arg t "NFTs on Avalanche: The Creator Economy" \
    --arg c "Avalanche's low fees and fast finality have made it a popular home for NFT marketplaces like Joepegs and Kalao. This piece looks at how creators are using Avalanche subnets for dedicated NFT drops, the role of royalties enforcement at the smart-contract level, and two case studies of artists who moved from Ethereum to Avalanche specifically to cut minting costs for their community. On royalty enforcement: Avalanche-native marketplaces increasingly use on-chain royalty registries enforced at the token-contract level via EIP-2981, closing the loophole where secondary marketplaces simply ignore creator royalties set off-chain." \
    --arg s "https://joepegs.com" '{title:$t, content:$c, sourceCitation:$s}')" -o /dev/null
submit "$A5" "$AMARA_TOKEN"

echo "== 6/9 editorial_approved: Tunde, awaiting designer banner =="
A6=$(create_article "$TUNDE_TOKEN" \
  "Avalanche vs Ethereum: A Technical Comparison" \
  "Both Avalanche and Ethereum run the EVM, but their consensus and finality models diverge sharply. This comparison covers block time, finality guarantees, validator hardware requirements, and gas fee dynamics side by side, aimed at developers deciding where to deploy a new contract rather than at general crypto audiences." \
  "https://ethereum.org, https://docs.avax.network")
submit "$A6" "$TUNDE_TOKEN"
review "$A6" "approved" "Solid technical comparison, accurate on finality timing. Approved." "$MODERATOR_TOKEN"

echo "== 7/9 banner_uploaded: Chidi, awaiting publisher =="
A7=$(create_article "$CHIDI_TOKEN" \
  "The Rise of GameFi on Avalanche" \
  "Avalanche's subnet architecture turned out to be a natural fit for game studios that need predictable, dedicated throughput without competing with unrelated DeFi traffic for block space. This article profiles three GameFi subnets, how their token economies are designed to avoid the inflationary death spiral that hit earlier play-to-earn titles, and what onboarding looks like for a player who has never held a crypto wallet before." \
  "https://www.avax.network/gaming")
submit "$A7" "$CHIDI_TOKEN"
review "$A7" "approved" "Great real-world examples, approved as-is." "$MODERATOR_TOKEN"
upload_banner "$A7" "$BANNER_DIR/banner1.jpg"
mark_ready "$A7"

echo "== 8/9 published, payment pending: Amara =="
A8=$(create_article "$AMARA_TOKEN" \
  "Avalanche's Path to Institutional Adoption" \
  "From JPMorgan's Onyx pilots to tokenized fund products, several institutions have chosen Avalanche subnets for regulated, permissioned deployments rather than the public C-Chain. This piece explains why a permissioned subnet appeals to institutions that need KYC'd validator sets, and reviews three live institutional pilots and what they reveal about Avalanche's enterprise strategy." \
  "https://www.avax.network/institutions")
submit "$A8" "$AMARA_TOKEN"
review "$A8" "approved" "Well-sourced, approved." "$MODERATOR_TOKEN"
upload_banner "$A8" "$BANNER_DIR/banner2.jpg"
mark_ready "$A8"
publish "$A8" "https://blogdeteam1.substack.com/p/institutional-adoption-demo"

echo "== 9/9 published, payment confirmed: Tunde =="
A9=$(create_article "$TUNDE_TOKEN" \
  "Building Your First Smart Contract on Avalanche" \
  "This walkthrough takes a developer from zero to a deployed ERC-20 token on Avalanche's Fuji testnet using Foundry, covering RPC setup, funding a test wallet from the faucet, writing a minimal contract, and verifying it on Snowtrace. It closes with notes on what changes when moving from Fuji to C-Chain mainnet." \
  "https://docs.avax.network/build")
submit "$A9" "$TUNDE_TOKEN"
review "$A9" "approved" "Clear, reproducible steps. Approved." "$MODERATOR_TOKEN"
upload_banner "$A9" "$BANNER_DIR/banner3.jpg"
mark_ready "$A9"
publish "$A9" "https://blogdeteam1.substack.com/p/first-smart-contract-demo"
release_payment "$A9"
sleep 3

echo "done"
