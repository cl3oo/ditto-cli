TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzczMzk4MjIsImlkIjoiNjllZWI4ZDZkMDMwZTU0OTkxYmEyMjRmIn0.W3jEsd1SOaLV9pyAK7tu-1KvtEfPfYpM_6AjR_kjZRQ"

test_curl() {
    local method=$1
    local url=$2
    local data=$3
    echo "Testing: $method $url"
    res=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$url" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "$data")
    echo "Result: $res"
    if [[ "$res" == "200" || "$res" == "201" ]]; then
        echo "SUCCESS: $url with $data"
        exit 0
    fi
}

PROTOCOLS=("http" "https")
SLASHES=("" "/")

for proto in "${PROTOCOLS[@]}"; do
    for slash in "${SLASHES[@]}"; do
        # Variation 1: /v1/posts with community_name
        test_curl "POST" "${proto}://api.ditto.local/v1/posts${slash}" '{"title": "Test", "content": "Test", "community_name": "random"}'
        
        # Variation 2: /v1/posts with community_id
        test_curl "POST" "${proto}://api.ditto.local/v1/posts${slash}" '{"title": "Test", "content": "Test", "community_id": "69eeb8e0b58542dea1be8549"}'
        
        # Variation 3: /v1/communities/random/posts
        test_curl "POST" "${proto}://api.ditto.local/v1/communities/random/posts${slash}" '{"title": "Test", "content": "Test"}'
        
        # Variation 4: /v1/communities/69eeb8e0b58542dea1be8549/posts
        test_curl "POST" "${proto}://api.ditto.local/v1/communities/69eeb8e0b58542dea1be8549/posts${slash}" '{"title": "Test", "content": "Test"}'
    done
done
